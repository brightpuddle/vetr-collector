package main

import (
	"collector/pkg/aci"

	_ "embed"

	"gopkg.in/yaml.v2"
)

//go:embed reqs.yaml
var reqsData []byte

// Mod modifies an aci Request
type Mod = func(*aci.Req)

// Request is an HTTP request.
type Request struct {
	Class  string            // MO class
	Prefix string            // Name for filename and class in DB
	Query  map[string]string // Query parameters
	path   string
}

func (req *Request) normalize() {
	if req.Prefix == "" {
		req.Prefix = req.Class
	}
	req.path = "/api/class/" + req.Class
}

func getRequests() (reqs []Request, err error) {
	err = yaml.Unmarshal(reqsData, &reqs)
	if err != nil {
		return
	}
	for i := 0; i < len(reqs); i++ {
		reqs[i].normalize()
	}
	return
}
