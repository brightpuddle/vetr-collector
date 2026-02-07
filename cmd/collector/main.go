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

func pause(msg string) {
	fmt.Println(msg)
	var throwaway string
	fmt.Scanln(&throwaway)
}

func main() {
	cfg, err := readArgs()
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading configuration.")
	}

	// Set log level based on verbose flag
	if anyVerbose(cfg) {
		log.SetLevel(zerolog.DebugLevel)
	} else {
		log.SetLevel(zerolog.InfoLevel)
	}

	if len(cfg.Fabrics) > 1 {
		runMultiFabric(cfg)
		return
	}

	runSingleFabric(cfg)
}

func runSingleFabric(cfg *config.Config) {
	fabric := cfg.Fabrics[0].MergeWithGlobal(cfg.Global)

	// Initialize ACI HTTP client
	client, err := cli.GetClient(fabric)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing ACI client.")
	}

	// Create results archive
	outputFile := fabric.GetOutputFileName()
	arc, err := archive.NewWriter(outputFile)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error creating archive file: %s.", outputFile)
	}

	// Initiate requests
	reqs, err := req.GetRequests()
	if err != nil {
		log.Fatal().Err(err).Msgf("Error reading requests.")
	}

	// Allow overriding in-built queries with a single class query
	if fabric.GetClass() != "all" {
		reqs = []req.Request{{
			Class: fabric.GetClass(),
			Query: fabric.Query,
		}}
	}

	// Batch and fetch queries in parallel
	collectErr := collectFabric(client, arc, reqs, fabric)

	arc.Close()
	log.Info().Msg("==============================")
	log.Info().Msg("Complete")
	log.Info().Msg("==============================")

	path, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot read current working directory")
	}
	outPath := filepath.Join(path, outputFile)

	if collectErr != nil {
		log.Warn().Err(collectErr).Msg("some data could not be fetched")
		log.Info().Msgf("Available data written to %s.", outPath)
	} else {
		log.Info().Msg("Collection complete.")
		log.Info().Msgf("Please provide %s to Cisco Services for further analysis.", outPath)
	}
	if !fabric.GetConfirm() {
		pause("Press enter to exit.")
	}
}

func runMultiFabric(cfg *config.Config) {
	log.Info().Msgf("Loaded config with %d fabric(s)", len(cfg.Fabrics))

	// Collect each fabric in parallel
	var g errgroup.Group
	outputFiles := make([]string, 0, len(cfg.Fabrics))
	for _, fabric := range cfg.Fabrics {
		fabric := fabric.MergeWithGlobal(cfg.Global)
		outputFiles = append(outputFiles, fabric.GetOutputFileName())
		g.Go(func() error {
			return collectSingleFabric(fabric)
		})
	}

	if err := g.Wait(); err != nil {
		log.Error().Err(err).Msg("Error collecting one or more fabrics")
	}

	if err := createAggregateArchive(outputFiles); err != nil {
		log.Error().Err(err).Msg("Failed to create aggregate archive")
	}

	log.Info().Msg("Multi-fabric collection complete.")
}

func collectSingleFabric(fabric config.FabricConfig) error {
	fabricName := fabric.GetFabricName()
	outputFile := fabric.GetOutputFileName()

	log := log.WithFabric(fabricName)
	log.Info().Msgf("Starting collection for fabric: %s", fabricName)

	// Initialize ACI HTTP client
	client, err := cli.GetClient(fabric)
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
	collectErr := collectFabric(client, arc, reqs, fabric)

	path, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot read current working directory: %w", err)
	}
	outPath := filepath.Join(path, outputFile)

	if collectErr != nil {
		log.Warn().Err(collectErr).Msgf("Some data could not be fetched for %s", fabricName)
	}

	log.Info().Str("path", outPath).Msg("Collection complete.")
	return collectErr
}

func collectFabric(
	client aci.Client,
	arc archive.Writer,
	reqs []req.Request,
	cfg config.FabricConfig,
) error {
	var logger log.Logger
	if cfg.GetFabricName() != "" {
		logger = log.WithFabric(cfg.GetFabricName())
	} else {
		logger = log.New()
	}

	batch := 1
	var firstErr error
	for i := 0; i < len(reqs); i += cfg.GetBatchSize() {
		var g errgroup.Group
		logger.Info().Msgf("Fetching request batch %d", batch)
		for j := i; j < i+cfg.GetBatchSize() && j < len(reqs); j++ {
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

func anyVerbose(cfg *config.Config) bool {
	if cfg.Global.Verbose {
		return true
	}
	for _, fabric := range cfg.Fabrics {
		if fabric.Verbose != nil && *fabric.Verbose {
			return true
		}
	}
	return false
}

func createAggregateArchive(files []string) error {
	const aggregateZip = "aci-collection.zip"
	arc, err := archive.NewWriter(aggregateZip)
	if err != nil {
		return err
	}
	defer arc.Close()

	seen := make(map[string]bool)
	for _, file := range files {
		if file == "" {
			continue
		}
		name := filepath.Base(file)
		if seen[name] {
			continue
		}
		seen[name] = true
		if _, err := os.Stat(file); err != nil {
			log.Warn().Err(err).Msgf("Skipping missing archive: %s", file)
			continue
		}
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read archive %s: %w", file, err)
		}
		if err := arc.Add(name, content); err != nil {
			return fmt.Errorf("failed to add %s to aggregate archive: %w", file, err)
		}
	}
	return nil
}
