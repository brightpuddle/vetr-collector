package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"collector/pkg/aci"
	"collector/pkg/archive"
	"collector/pkg/req"

	"golang.org/x/sync/errgroup"

	"collector/pkg/log"

	"github.com/tidwall/gjson"
)

// Config is CLI conifg
type Config struct {
	Host              string
	Username          string
	Password          string
	RetryDelay        int
	RequestRetryCount int
	BatchSize         int
	PageSize          int
	Confirm           bool
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
	// Sanatize username against quotes
	cfg.Password = strings.ReplaceAll(cfg.Password, "\"", "\\\"")
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

func fetchWithRetry(
	client aci.Client,
	path string,
	cfg Config,
	mods []func(*aci.Req),
) (gjson.Result, error) {
	res, err := client.Get(path, mods...)
	if err != nil && err.Error() == "result dataset is too big" {
		return res, err
	}

	// Retry for requestRetryCount times
	for retries := 0; err != nil && retries < cfg.RequestRetryCount; retries++ {
		log.Warn().Err(err).Msgf("request failed for %s. Retrying after %d seconds.",
			path, cfg.RetryDelay)
		time.Sleep(time.Second * time.Duration(cfg.RetryDelay))
		res, err = client.Get(path, mods...)
	}
	if err != nil {
		return res, fmt.Errorf("request failed for %s: %v", path, err)
	}
	return res, nil
}

// Fetch fetches data via API and writes it to the provided archive.
func Fetch(client aci.Client, req req.Request, arc archive.Writer, cfg Config) error {
	path := "/api/class/" + req.Class
	startTime := time.Now()
	log.Debug().Time("start_time", startTime).Msgf("begin: %s", req.Class)

	log.Info().Msgf("fetching %s...", req.Class)

	mods := []func(*aci.Req){}
	for k, v := range req.Query {
		mods = append(mods, aci.Query(k, v))
	}

	// Handle tenants individually for scale purposes
	res, err := fetchWithRetry(client, path, cfg, mods)
	if err != nil && err.Error() == "result dataset is too big" {
		if err := paginate(client, req, arc, cfg, mods); err != nil {
			return err
		}
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

func paginate(
	client aci.Client,
	req req.Request,
	arc archive.Writer,
	cfg Config,
	mods []func(*aci.Req),
) error {
	path := "/api/class/" + req.Class
	log.Info().Msgf("fetching large dataset for %s...", req.Class)
	mods = append(mods, aci.Query("page-size", strconv.Itoa(cfg.PageSize)))

	log.Info().Msgf("fetching page 0 for %s...", req.Class)
	res, err := fetchWithRetry(client, path, cfg, mods)
	if err != nil {
		return err
	}

	cnt, _ := strconv.Atoi(res.Get("totalCount").Str)

	log.Info().Msgf("Total record count for %s: %d", req.Class, cnt)
	pages := cnt / cfg.PageSize

	batch := 1
	for i := 0; i < pages; i += cfg.BatchSize {
		var g errgroup.Group
		fmt.Println(strings.Repeat("*", 30))
		fmt.Println("Fetching paginated request batch", batch)
		fmt.Println(strings.Repeat("*", 30))
		for j := i; j < i+cfg.BatchSize && j < pages; j++ {
			page := j
			g.Go(func() error {
				log.Info().Msgf("fetching page %d of %d for %s...", page, pages, req.Class)
				mods := append(mods, aci.Query("page", strconv.Itoa(page)))
				res, err := fetchWithRetry(client, path, cfg, mods)
				if err != nil {
					return fmt.Errorf("failed to fetch large dataset for %s", req.Class)
				}
				log.Info().Msgf("%d of %d for %s complete", page, pages, req.Class)
				err = arc.Add(fmt.Sprintf("%s-%d.json", req.Class, page), []byte(res.Raw))
				if err != nil {
					return fmt.Errorf("failed to write large dataset for %s", req.Class)
				}
				return nil
			})
		}
		err = g.Wait()
		if err != nil {
			log.Error().Err(err).Msg("Error fetching data.")
		}
		batch++
	}
	return nil
}
