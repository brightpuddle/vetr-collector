package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestLogger(t *testing.T) {
	a := assert.New(t)

	fileBuf := &bytes.Buffer{}
	consoleBuf := &bytes.Buffer{}
	writer := MultiLevelWriter{
		file:    fileBuf,
		console: consoleBuf,
	}
	log := zerolog.New(writer).With().Timestamp().Logger()

	log.Debug().Msg("debug_test")
	log.Info().Msg("info_test")
	json := gjson.ParseBytes(fileBuf.Bytes())
	a.Equal("debug_test", json.Get(`..#(level="debug").message`).Str)
	a.Equal("info_test", json.Get(`..#(level="info").message`).Str)
	console := consoleBuf.String()
	a.True(strings.Contains(console, "info_test"))
	a.False(strings.Contains(console, "debug_test"))
}
