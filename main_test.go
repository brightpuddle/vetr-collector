package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/brightpuddle/goaci"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
)

func TestWriteScript(t *testing.T) {
	a := assert.New(t)
	log := zerolog.New(&bytes.Buffer{})

	err := writeScript(log)
	a.NoError(err)
	defer os.Remove(logFile)
	fs, err := os.Stat("vetr-collector.sh")
	if a.NoError(err) {
		a.True(fs.Size() > 300)
	}
}

func TestReadRaw(t *testing.T) {
	a := assert.New(t)
	log := zerolog.New(&bytes.Buffer{})

	inPath := filepath.Join("testdata", "aci-vetr-raw.zip")
	outPath := filepath.Join("testdata", "script-data.zip")
	err := readRaw(inPath, outPath, log)
	a.NoError(err)
	defer os.Remove(outPath)
	fs, err := os.Stat(outPath)
	if a.NoError(err) {
		a.True(fs.Size() > 300)
	}
}

func TestFetch(t *testing.T) {
	a := assert.New(t)
	defer gock.Off()

	gock.New("https://apic").
		Get("/api/class/fvTenant.json").
		Reply(200).
		BodyString(goaci.Body{}.
			Set("imdata.0.fvTenant.attributes.dn", "uni/tn-zero").
			Set("imdata.1.fvTenant.attributes.dn", "uni/tn-one").
			Str)
	client, _ := goaci.NewClient("apic", "usr", "pwd")
	client.LastRefresh = time.Now()
	gock.InterceptClient(client.HttpClient)

	log := zerolog.New(&bytes.Buffer{})
	reqs := []*Request{{
		class:  "fvTenant",
		path:   "/api/class/fvTenant",
		filter: "#.fvTenant.attribute",
	}}
	results, err := fetch(client, reqs, log)
	a.NoError(err)
	if tenants, ok := results["fvTenant"]; ok {
		a.Equal("uni/tn-zero", tenants.Get("0.dn").Str)
		a.Equal("uni/tn-one", tenants.Get("1.dn").Str)
	}
}
