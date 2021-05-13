package main

import (
	"fmt"
	"os"
	"strings"

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
	log = newLogger()
	args = newArgs()

	// Initialize ACI HTTP client
	client, err := getClient(args.APIC, args.Username, args.Password)
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing ACI client.")
	}

	// Create results archive
	os.Remove(args.Output)
	arc, err := newArchiveWriter(args.Output)
	if err != nil {
		log.Fatal().Err(err).Msgf("Error creating archive file: %s.", args.Output)
	}
	defer arc.close()

	// Initiate requests
	fmt.Println(strings.Repeat("=", 30))
	var g errgroup.Group

	for _, req := range getRequests() {
		req := req
		g.Go(func() error {
			return fetchResource(client, req, arc)
		})
	}
	err = g.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("Error fetching data.")
	}
	fmt.Println(strings.Repeat("=", 30))

	if err != nil {
		log.Debug().Err(err).Msg("some data could not be fetched")
	} else {
		log.Info().Msg("Collection complete.")
		log.Info().Msgf("Please provide %s to Cisco Services for further analysis.", args.Output)
	}
	pause("Press enter to exit.")
}
