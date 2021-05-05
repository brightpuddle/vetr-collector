package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/brightpuddle/goaci"
	"github.com/mholt/archiver"
	"github.com/rs/zerolog"
	"github.com/tidwall/buntdb"
	"github.com/tidwall/gjson"
	"golang.org/x/sync/errgroup"
)

// Version comes from CI
var (
	version string
	mux     sync.Mutex
)

const (
	resultZip  = "aci-vetr-data.zip"
	scriptName = "vetr-collector.sh"
	logFile    = "collector.log"
	dbName     = "data.db"
)

// Write requests to script to be run on the APIC.
// Note, this is a more complicated collection methodology and should rarely
// be used.
func writeScript(log zerolog.Logger) error {
	var (
		final     = "aci-vetr-raw.zip"
		tmpFolder = "/tmp/aci-vetr-collections"
	)
	os.Remove(scriptName)
	script := []string{
		"#!/bin/bash",
		"",
		"mkdir " + tmpFolder,
		"",
		"# Fetch data from API",
	}

	client := goaci.Client{}

	for _, request := range getRequests() {
		req := client.NewReq("GET", request.path, nil, request.mods...)
		cmd := fmt.Sprintf("icurl -kG https://localhost/%s", req.HttpReq.URL.Path)

		for key, value := range req.HttpReq.URL.Query() {
			if len(value) >= 1 {
				cmd = fmt.Sprintf("%s -d '%s=%s'", cmd, key, value[0])
			}
		}
		cmd = fmt.Sprintf("%s > %s/%s", cmd, tmpFolder, request.prefix+".json")
		script = append(script, cmd)
	}

	script = append(script, []string{
		"",
		"# Zip result",
		fmt.Sprintf("zip -mj ~/%s %s/*.json", final, tmpFolder),
		"",
		"# Cleanup",
		"rm -rf " + tmpFolder,
		"",
		"echo Collection complete.",
		fmt.Sprintf("echo Provide Cisco Services the %s file.", final),
	}...)

	err := ioutil.WriteFile(scriptName, []byte(strings.Join(script, "\n")), 0755)
	if err != nil {
		return err
	}
	log.Info().Msgf("Script complete. Run %s on the APIC.", scriptName)
	return nil
}

// Translate raw (script) data to aci-vetr-data.zip file for backend consumption.
func readRaw(in, out string, log zerolog.Logger) error {
	results := make(map[string]goaci.Res)
	// Read data from zip
	err := archiver.Walk(in, func(f archiver.File) error {
		zfh, ok := f.Header.(zip.FileHeader)
		if ok && strings.HasSuffix(zfh.Name, ".json") {
			prefix := strings.TrimSuffix(zfh.Name, ".json")
			b, err := ioutil.ReadAll(f)
			if err != nil {
				return err
			}
			json := gjson.ParseBytes(b)
			results[prefix] = json
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error reading from archive: %v", err)
	}

	// Apply filters
	for _, request := range getRequests() {
		if res, ok := results[request.prefix]; ok {
			results[request.prefix] = res.Get("imdata." + request.filter)
		}
	}

	// Write to DB
	if err := writeToDB(results); err != nil {
		return fmt.Errorf("error writing to DB: %v", err)
	}
	defer os.Remove(dbName)

	// Create archive
	log.Info().Msg("Creating archive")
	os.Remove(out) // Remove any old archives and ignore errors
	if err := archiver.Archive([]string{dbName}, out); err != nil {
		return fmt.Errorf("cannot create archive: %v", err)
	}

	// Cleanup
	fmt.Println(strings.Repeat("=", 30))
	log.Info().Msgf("Please provide %s to Cisco Services for further analysis.", out)
	return nil
}

// Write results to db file.
func writeToDB(responses map[string]goaci.Res) error {
	db, err := buntdb.Open(dbName)
	if err != nil {
		return fmt.Errorf("cannot open output file: %v", err)
	}
	defer db.Close()

	for prefix, res := range responses {
		if err := db.Update(func(tx *buntdb.Tx) error {
			for _, record := range res.Array() {
				key := fmt.Sprintf("%s:%s", prefix, record.Get("dn").Str)
				if _, _, err := tx.Set(key, record.Raw, nil); err != nil {
					return fmt.Errorf("cannot set key: %v", err)
				}
			}
			return nil
		}); err != nil {
			return fmt.Errorf("cannot write to DB file: %v", err)
		}
	}

	// Add metadata
	metadata := goaci.Body{}.
		Set("collectorVersion", version).
		Set("timestamp", time.Now().String()).
		Str
	return db.Update(func(tx *buntdb.Tx) error {
		if _, _, err := tx.Set("meta", string(metadata), nil); err != nil {
			return fmt.Errorf("cannot write metadata to db: %v", err)
		}
		return nil
	})
}

func fetch(client goaci.Client, reqs []*Request, log Logger) (map[string]goaci.Res, error) {
	responses := make(map[string]goaci.Res)
	var g errgroup.Group

	for _, req := range reqs {
		req := req

		g.Go(func() error {
			startTime := time.Now()
			log.Debug().Time("start_time", startTime).Msgf("begin: %s", req.prefix)

			log.Info().Str("resource", req.prefix).Msg("fetching resource...")
			log.Debug().Str("url", req.path).Msg("requesting resource")

			res, err := client.Get(req.path, req.mods...)
			if err != nil {
				log.Warn().Err(err).Msgf("skipping failed request for %s", req.path)
				return fmt.Errorf("failed to make request for %s: %v", req.path, err)
			}
			mux.Lock()
			responses[req.prefix] = res.Get("imdata." + req.filter)
			mux.Unlock()
			log.Debug().
				TimeDiff("elapsed_time", time.Now(), startTime).
				Msgf("done: %s", req.prefix)
			return nil
		})
	}

	err := g.Wait()
	return responses, err
}

// Fetch data via API.
func fetchHTTP(args Args, log zerolog.Logger) error {
	client, err := goaci.NewClient(
		args.APIC,
		args.Username,
		args.Password,
		goaci.RequestTimeout(600),
	)
	if err != nil {
		return fmt.Errorf("failed to create ACI client: %v", err)
	}

	// Authenticate
	log.Info().Str("host", args.APIC).Msg("APIC host")
	log.Info().Str("user", args.Username).Msg("APIC username")
	log.Info().Msg("Authenticating to the APIC...")
	if err := client.Login(); err != nil {
		return fmt.Errorf("cannot authenticate to the APIC at %s: %v", args.APIC, err)
	}

	// Fetch data from API
	fmt.Println(strings.Repeat("=", 30))

	responses, err := fetch(client, getRequests(), log)
	if err != nil {
		// TODO checking the error string is a hack. Push code into goaci upstream.
		if !strings.Contains(err.Error(), "400") {
			return err
		}
	}

	if err := writeToDB(responses); err != nil {
		return fmt.Errorf("error writing to DB: %v", err)
	}

	fmt.Println(strings.Repeat("=", 30))

	// Create archive
	log.Info().Msg("Creating archive")
	os.Remove(args.Output) // Remove any old archives and ignore errors
	if err := archiver.Archive([]string{dbName, logFile}, args.Output); err != nil {
		return fmt.Errorf("cannot create archive: %v", err)
	}

	// Cleanup
	fmt.Println(strings.Repeat("=", 30))
	log.Info().Msg("Collection complete.")
	log.Info().Msgf("Please provide %s to Cisco Services for further analysis.", args.Output)
	return nil
}

func main() {
	log := newLogger()
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				log.Error().Err(err).Msg("unexpected error")
			}
			log.Error().Msg("Collection failed.")
		} else {
			// TODO move cleanup into the archive lib, e.g. zip -m
			os.Remove(logFile)
		}
		os.Remove(dbName)
		fmt.Println("Press enter to exit.")
		var throwaway string
		fmt.Scanln(&throwaway)
	}()
	args, err := newArgs()
	if err != nil {
		panic(err)
	}
	switch {
	case args.WriteScript:
		err := writeScript(log)
		if err != nil {
			log.Error().Err(err).Msg("cannot create script")
		}
	case args.ReadRaw != "":
		err := readRaw(args.ReadRaw, args.Output, log)
		if err != nil {
			log.Error().Err(err).Msg("cannot read script output")
		}
	default:
		err := fetchHTTP(args, log)
		if err != nil {
			log.Debug().Err(err).Msg("some data could not be fetched")
		}
	}
}
