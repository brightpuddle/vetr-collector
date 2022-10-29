// Package req contains the collector requests
package req

import (
	"collector/pkg/aci"

	_ "embed"

	"gopkg.in/yaml.v2"
)

//go:embed reqs.json
var reqsData []byte

// Mod modifies an aci Request
type Mod = func(*aci.Req)

// Request is an HTTP request.
type Request struct {
	Class string            // MO class
	Query map[string]string // Query parameters
}

// GetRequests returns normalized requests
func GetRequests() (reqs []Request, err error) {
	err = yaml.Unmarshal(reqsData, &reqs)
	return
}
