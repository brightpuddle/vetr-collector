package collector

import (
	_ "embed"
	"strings"
)

// Version is the vetR collector version
//go:embed VERSION
var Version string

func init() {
	Version = strings.Trim(Version, "\n")
}
