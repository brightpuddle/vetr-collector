package main

import (
	"bytes"
	"testing"
	"time"

	"collector/aci"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"gopkg.in/h2non/gock.v1"
)

type mockArchiveWriter struct {
	files map[string][]byte
}

func (a mockArchiveWriter) close() error {
	return nil
}

func (a mockArchiveWriter) add(name string, content []byte) error {
	a.files[name] = content
	return nil
}

func TestFetch(t *testing.T) {
	a := assert.New(t)
	defer gock.Off()

	// Overwrite logger with bin bucket
	log = zerolog.New(&bytes.Buffer{})

	// Mock API
	gock.New("https://apic").
		Get("/api/class/myClass.json").
		Reply(200).
		BodyString(aci.Body{}.
			Set("imdata.0.myClass.attributes.dn", "uni/my-zero").
			Set("imdata.1.myClass.attributes.dn", "uni/my-one").
			Str)

	// Test client
	client, _ := aci.NewClient("apic", "usr", "pwd")
	client.LastRefresh = time.Now()
	gock.InterceptClient(client.HTTPClient)

	// Test request
	req := Request{
		Class: "myClass",
	}
	req.normalize()

	// Mock archive
	arc := mockArchiveWriter{
		files: make(map[string][]byte),
	}

	// Write zip
	err := fetchResource(client, req, arc)
	a.NoError(err)

	// Verify content written to mock archive
	content, ok := arc.files["myClass.json"]
	a.True(ok)
	classes := gjson.ParseBytes(content).Get("imdata.#.myClass.attributes")
	a.Equal("uni/my-zero", classes.Get("0.dn").Str)
	a.Equal("uni/my-one", classes.Get("1.dn").Str)
}
