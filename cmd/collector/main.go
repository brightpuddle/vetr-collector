package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"collector/pkg/archive"
	"collector/pkg/cli"
	"collector/pkg/req"

	"collector/pkg/logger"

	"golang.org/x/sync/errgroup"
)

var args Args

func pause(msg string) {
	fmt.Println(msg)
	var throwaway string
	fmt.Scanln(&throwaway)
}

func main() {
	log, err := logger.New(logger.Config{
		Filename:     "collector.log",
		ConsoleLevel: logger.InfoLevel,
	})
	if err != nil {
		panic(err)
	}
	args = newArgs()
	cfg := cli.Config{
		Host:              args.APIC,
		Username:          args.Username,
		Password:          args.Password,
		RetryDelay:        args.RetryDelay,
		RequestRetryCount: args.RequestRetryCount,
		BatchSize:         args.BatchSize,
		PageSize:          args.PageSize,
		Confirm:           args.Confirm,
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
	batch := 1
	for i := 0; i < len(reqs); i += args.BatchSize {
		var g errgroup.Group
		fmt.Println(strings.Repeat("=", 30))
		fmt.Println("Fetching request batch", batch)
		fmt.Println(strings.Repeat("=", 30))
		for j := i; j < i+args.BatchSize && j < len(reqs); j++ {
			req := reqs[j]
			g.Go(func() error {
				return cli.Fetch(client, req, arc, cfg)
			})
		}
		err = g.Wait()
		if err != nil {
			log.Error().Err(err).Msg("Error fetching data.")
		}
		batch++
	}
	arc.Close()
	fmt.Println(strings.Repeat("=", 30))
	fmt.Println("Complete")
	fmt.Println(strings.Repeat("=", 30))

	path, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot read current working directory")
	}
	outPath := filepath.Join(path, args.Output)

	if err != nil {
		log.Warn().Err(err).Msg("some data could not be fetched")
		log.Info().Err(err).Msgf("Available data written to %s.", outPath)
	} else {
		log.Info().Msg("Collection complete.")
		log.Info().Msgf("Please provide %s to Cisco Services for further analysis.", outPath)
	}
	if !cfg.Confirm {
		pause("Press enter to exit.")
	}
}
