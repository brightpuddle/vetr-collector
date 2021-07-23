package main

import (
	"fmt"
	"time"

	"collector/aci"
)

func getClient(host, usr, pwd string) (aci.Client, error) {
	client, err := aci.NewClient(
		host, usr, pwd,
		aci.RequestTimeout(600),
	)
	if err != nil {
		return aci.Client{}, fmt.Errorf("failed to create ACI client: %v", err)
	}

	// Authenticate
	log.Info().Str("host", host).Msg("APIC host")
	log.Info().Str("user", usr).Msg("APIC username")
	log.Info().Msg("Authenticating to the APIC...")
	if err := client.Login(); err != nil {
		return aci.Client{}, fmt.Errorf("cannot authenticate to the APIC at %s: %v", host, err)
	}
	return client, nil
}

// Fetch data via API.
func fetchResource(client aci.Client, req Request, arc archiveWriter) error {
	startTime := time.Now()
	log.Debug().Time("start_time", startTime).Msgf("begin: %s", req.Prefix)

	log.Info().Msgf("fetching %s...", req.Prefix)
	log.Debug().Str("url", req.path).Msg("requesting resource")

	var mods []func(*aci.Req)
	for k, v := range req.Query {
		mods = append(mods, aci.Query(k, v))
	}

	res, err := client.Get(req.path, mods...)
	// Retry for requestRetryCount times
	for retries := 0; err != nil && retries < args.RequestRetryCount; retries++ {
		log.Warn().Err(err).Msgf("request failed for %s. Retrying after %d seconds.",
			req.path, args.RetryDelay)
		time.Sleep(time.Second * time.Duration(args.RetryDelay))
		res, err = client.Get(req.path, mods...)
	}
	if err != nil {
		return fmt.Errorf("request failed for %s: %v", req.path, err)
	}
	log.Info().Msgf("%s complete", req.Prefix)
	// err = arc.add(req.Prefix+".json", []byte(res.Get(req.Filter).Raw))
	err = arc.add(req.Prefix+".json", []byte(res.Raw))
	if err != nil {
		return err
	}
	log.Debug().
		TimeDiff("elapsed_time", time.Now(), startTime).
		Msgf("done: %s", req.Prefix)
	return nil
}

func pause(msg string) {
	fmt.Println("Press enter to exit.")
	var throwaway string
	fmt.Scanln(&throwaway)
}
