package cli

import (
	"collector/pkg/aci"
	"collector/pkg/archive"
	"collector/pkg/req"
	"fmt"
	"time"

	"github.com/aci-vetr/bats/logger"
)

// Config is CLI conifg
type Config struct {
	Host              string
	Username          string
	Password          string
	RetryDelay        int
	RequestRetryCount int
	BatchSize         int
}

// NewConfig populates default values
func NewConfig() Config {
	return Config{
		RequestRetryCount: 3,
		RetryDelay:        10,
		BatchSize:         10,
	}
}

// GetClient creates an ACI host client
func GetClient(cfg Config) (aci.Client, error) {
	log := logger.Get()
	client, err := aci.NewClient(
		cfg.Host, cfg.Username, cfg.Password,
		aci.RequestTimeout(600),
	)
	if err != nil {
		return aci.Client{}, fmt.Errorf("failed to create ACI client: %v", err)
	}

	// Authenticate
	log.Info().Str("host", cfg.Host).Msg("APIC host")
	log.Info().Str("user", cfg.Username).Msg("APIC username")
	log.Info().Msg("Authenticating to the APIC...")
	if err := client.Login(); err != nil {
		return aci.Client{}, fmt.Errorf("cannot authenticate to the APIC at %s: %v", cfg.Host, err)
	}
	return client, nil
}

// FetchResource fetches data via API and writes it to the provided archive.
func FetchResource(client aci.Client, req req.Request, arc archive.Writer, cfg Config) error {
	path := "/api/class/" + req.Class
	log := logger.Get()
	startTime := time.Now()
	log.Debug().Time("start_time", startTime).Msgf("begin: %s", req.Class)

	log.Info().Msgf("fetching %s...", req.Class)
	log.Debug().Str("url", path).Msg("requesting resource")

	var mods []func(*aci.Req)
	for k, v := range req.Query {
		mods = append(mods, aci.Query(k, v))
	}

	// Handle tenants individually for scale purposes
	res, err := client.Get(path, mods...)
	// Retry for requestRetryCount times
	for retries := 0; err != nil && retries < cfg.RequestRetryCount; retries++ {
		log.Warn().Err(err).Msgf("request failed for %s. Retrying after %d seconds.",
			path, cfg.RetryDelay)
		time.Sleep(time.Second * time.Duration(cfg.RetryDelay))
		res, err = client.Get(path, mods...)
	}
	if err != nil {
		return fmt.Errorf("request failed for %s: %v", path, err)
	}
	log.Info().Msgf("%s complete", req.Class)
	err = arc.Add(req.Class+".json", []byte(res.Raw))
	if err != nil {
		return err
	}
	log.Debug().
		TimeDiff("elapsed_time", time.Now(), startTime).
		Msgf("done: %s", req.Class)
	return nil
}
