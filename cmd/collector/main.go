package main

import (
	"fmt"
	"os"
	"path/filepath"

	"collector/pkg/aci"
	"collector/pkg/archive"
	"collector/pkg/cli"
	"collector/pkg/config"
	"collector/pkg/log"
	"collector/pkg/req"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

var args Args

func pause(msg string) {
	fmt.Println(msg)
	var throwaway string
	fmt.Scanln(&throwaway)
}

func main() {
	args = newArgs()

	// Set log level based on verbose flag
	if args.Verbose {
		log.SetLevel(zerolog.DebugLevel)
	} else {
		log.SetLevel(zerolog.InfoLevel)
	}

	// If config file is provided, use multi-fabric mode
	if args.ConfigFile != "" {
		runMultiFabric()
		return
	}

	// Single fabric mode (backward compatibility)
	runSingleFabric()
}

func runSingleFabric() {
	cfg := cli.Config{
		Host:              args.URL,
		Username:          args.Username,
		Password:          args.Password,
		RetryDelay:        args.RetryDelay,
		RequestRetryCount: args.RequestRetryCount,
		BatchSize:         args.BatchSize,
		PageSize:          args.PageSize,
		Confirm:           args.Confirm,
		FabricName:        "",
	}

	// Initialize ACI HTTP client
	client, err := cli.GetClient(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing ACI client.")
	}

	// Create results archive
	arc, err := archive.NewWriter(args.Output)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error creating archive file: %s.", args.Output)
	}

	// Initiate requests
	reqs, err := req.GetRequests()
	if err != nil {
		log.Fatal().Err(err).Msgf("Error reading requests.")
	}

	// Allow overriding in-built queries with a single class query
	if args.Class != "" && args.Class != "all" {
		reqs = []req.Request{{
			Class: args.Class,
			Query: args.Query,
		}}
	}

	// Batch and fetch queries in parallel
	collectErr := collectFabric(client, arc, reqs, cfg)

	arc.Close()
	log.Info().Msg("==============================")
	log.Info().Msg("Complete")
	log.Info().Msg("==============================")

	path, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot read current working directory")
	}
	outPath := filepath.Join(path, args.Output)

	if collectErr != nil {
		log.Warn().Err(collectErr).Msg("some data could not be fetched")
		log.Info().Msgf("Available data written to %s.", outPath)
	} else {
		log.Info().Msg("Collection complete.")
		log.Info().Msgf("Please provide %s to Cisco Services for further analysis.", outPath)
	}
	if !cfg.Confirm {
		pause("Press enter to exit.")
	}
}

func runMultiFabric() {
	// Load config file
	cfg, err := config.LoadConfig(args.ConfigFile)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error loading config file: %s", args.ConfigFile)
	}

	// Set log level from config if verbose is set globally
	if cfg.Global.Verbose {
		log.SetLevel(zerolog.DebugLevel)
	}

	log.Info().Msgf("Loaded config file with %d fabric(s)", len(cfg.Fabrics))

	// Collect each fabric in parallel
	var g errgroup.Group
	for _, fabric := range cfg.Fabrics {
		fabric := fabric.MergeWithGlobal(cfg.Global)
		g.Go(func() error {
			return collectSingleFabric(fabric)
		})
	}

	if err := g.Wait(); err != nil {
		log.Error().Err(err).Msg("Error collecting one or more fabrics")
	}

	log.Info().Msg("Multi-fabric collection complete.")
}

func collectSingleFabric(fabric config.FabricConfig) error {
	fabricName := fabric.GetFabricName()
	outputFile := fabric.GetOutputFileName()

	logger := log.WithFabric(fabricName)
	logger.Info().Msgf("Starting collection for fabric: %s", fabricName)

	// Build CLI config from fabric config
	cliCfg := cli.Config{
		Host:              fabric.URL,
		Username:          fabric.Username,
		Password:          fabric.Password,
		RetryDelay:        fabric.GetRetryDelay(),
		RequestRetryCount: fabric.GetRequestRetryCount(),
		BatchSize:         fabric.GetBatchSize(),
		PageSize:          fabric.GetPageSize(),
		Confirm:           fabric.GetConfirm(),
		FabricName:        fabricName,
	}

	// Initialize ACI HTTP client
	client, err := cli.GetClient(cliCfg)
	if err != nil {
		return fmt.Errorf("error initializing ACI client for %s: %w", fabricName, err)
	}

	// Create results archive
	arc, err := archive.NewWriter(outputFile)
	if err != nil {
		return fmt.Errorf("error creating archive file %s: %w", outputFile, err)
	}
	defer arc.Close()

	// Initiate requests
	reqs, err := req.GetRequests()
	if err != nil {
		return fmt.Errorf("error reading requests for %s: %w", fabricName, err)
	}

	// Allow overriding in-built queries with a single class query
	if fabric.GetClass() != "all" {
		reqs = []req.Request{{
			Class: fabric.GetClass(),
			Query: fabric.Query,
		}}
	}

	// Batch and fetch queries in parallel
	collectErr := collectFabric(client, arc, reqs, cliCfg)

	path, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot read current working directory: %w", err)
	}
	outPath := filepath.Join(path, outputFile)

	if collectErr != nil {
		logger.Warn().Err(collectErr).Msgf("Some data could not be fetched for %s", fabricName)
	}

	logger.Info().Msgf("Collection complete for %s. Output: %s", fabricName, outPath)
	return collectErr
}

func collectFabric(client aci.Client, arc archive.Writer, reqs []req.Request, cfg cli.Config) error {
	var logger log.Logger
	if cfg.FabricName != "" {
		logger = log.WithFabric(cfg.FabricName)
	} else {
		logger = log.New()
	}

	batch := 1
	var firstErr error
	for i := 0; i < len(reqs); i += cfg.BatchSize {
		var g errgroup.Group
		logger.Info().Msg("==============================")
		logger.Info().Msgf("Fetching request batch %d", batch)
		logger.Info().Msg("==============================")
		for j := i; j < i+cfg.BatchSize && j < len(reqs); j++ {
			req := reqs[j]
			g.Go(func() error {
				return cli.Fetch(client, req, arc, cfg)
			})
		}
		err := g.Wait()
		if err != nil {
			logger.Error().Err(err).Msg("Error fetching data.")
			if firstErr == nil {
				firstErr = err
			}
		}
		batch++
	}
	return firstErr
}
