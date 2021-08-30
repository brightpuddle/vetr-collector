package main

import (
	"fmt"
	"strings"

	"collector/pkg/archive"
	"collector/pkg/logger"
	"collector/pkg/req"

	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

// Version comes from CI
var (
	version string
	log     zerolog.Logger
	args    Args
)

func main() {
	log = logger.New()
	args = newArgs()

	// Initialize ACI HTTP client
	client, err := getClient(args.APIC, args.Username, args.Password)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing ACI client.")
	}

	outFilename := nextFilename(args.Output)

	// Create results archive
	arc, err := archive.NewWriter(outFilename)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error creating archive file: %s.", outFilename)
	}
	defer arc.Close()

	// Initiate requests
	reqs, err := req.GetRequests()
	if err != nil {
		log.Fatal().Err(err).Msgf("Error reading requests.")
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
				return fetchResource(client, req, arc)
			})
		}
		err = g.Wait()
		if err != nil {
			log.Error().Err(err).Msg("Error fetching data.")
		}
		batch++
	}

	fmt.Println(strings.Repeat("=", 30))
	fmt.Println("Complete")
	fmt.Println(strings.Repeat("=", 30))

	if err != nil {
		log.Warn().Err(err).Msg("some data could not be fetched")
		log.Info().Err(err).Msgf("Available data written to %s.", outFilename)
	} else {
		log.Info().Msg("Collection complete.")
		log.Info().Msgf("Please provide %s to Cisco Services for further analysis.", outFilename)
	}
	pause("Press enter to exit.")
}
