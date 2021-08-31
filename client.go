package main

import (
	"fmt"
	"sync"
	"time"

	"collector/pkg/aci"
	"collector/pkg/archive"
	"collector/pkg/req"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"golang.org/x/sync/errgroup"
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

// fetchTenants splits the fvTenant query for large configs
func fetchTenants(client aci.Client, req req.Request, mods []func(*aci.Req)) (gjson.Result, error) {
	all := ""
	log.Info().Msg("fetching tenant list...")
	// Fetch only tenants
	res, err := client.Get(req.Path)
	if err != nil {
		return gjson.Result{}, err
	}
	tns := res.Get("imdata.#.fvTenant.attributes").Array()

	// Batch tenant requests
	batch := 1
	for i := 0; i < len(tns); i += args.BatchSize {
		var (
			g  errgroup.Group
			mu sync.Mutex
		)
		log.Debug().Msgf("Fetching tenant request batch %d", batch)
		for j := i; j < i+args.BatchSize && j < len(tns); j++ {
			tn := tns[j]
			g.Go(func() error {
				dn := tn.Get("dn").Str
				log.Info().Msgf("fetching Tenant %s", tn.Get("name").Str)

				res, err := client.Get("/api/mo/"+dn, mods...)
				if err != nil {
					return fmt.Errorf("cannot fetch tenant %s: %w", tn.Get("name").Str, err)
				}
				log.Info().Msgf("Tenant %s complete", tn.Get("name").Str)
				mu.Lock()
				all, _ = sjson.SetRaw(all, "imdata.-1", res.Get("imdata.0").Raw)
				mu.Unlock()
				return nil
			})
		}
		if err := g.Wait(); err != nil {
			return gjson.Result{}, err
		}
		batch++
	}
	return gjson.Parse(all), nil
}

// Fetch data via API.
func fetchResource(client aci.Client, req req.Request, arc archive.Writer) error {
	startTime := time.Now()
	log.Debug().Time("start_time", startTime).Msgf("begin: %s", req.Prefix)

	log.Info().Msgf("fetching %s...", req.Prefix)
	log.Debug().Str("url", req.Path).Msg("requesting resource")

	var mods []func(*aci.Req)
	for k, v := range req.Query {
		mods = append(mods, aci.Query(k, v))
	}
	var (
		res gjson.Result
		err error
	)
	// Handle tenants individually for scale purposes
	if req.Prefix == "fvTenant" {
		res, err = fetchTenants(client, req, mods)
	} else {
		res, err = client.Get(req.Path, mods...)
		// Retry for requestRetryCount times
		for retries := 0; err != nil && retries < args.RequestRetryCount; retries++ {
			log.Warn().Err(err).Msgf("request failed for %s. Retrying after %d seconds.",
				req.Path, args.RetryDelay)
			time.Sleep(time.Second * time.Duration(args.RetryDelay))
			res, err = client.Get(req.Path, mods...)
		}
	}
	if err != nil {
		return fmt.Errorf("request failed for %s: %v", req.Path, err)
	}
	log.Info().Msgf("%s complete", req.Prefix)
	// err = arc.add(req.Prefix+".json", []byte(res.Get(req.Filter).Raw))
	err = arc.Add(req.Prefix+".json", []byte(res.Raw))
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
