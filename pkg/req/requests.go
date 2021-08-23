package req

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
	Path   string
}

// Normalize chooses correct class and path.
// This should be run on all manually created Requests structs.
func (req *Request) Normalize() *Request {
	if req.Prefix == "" {
		req.Prefix = req.Class
	}
	req.Path = "/api/class/" + req.Class
	return req
}

// GetRequests returns normalized requests
func GetRequests() (reqs []Request, err error) {
	err = yaml.Unmarshal(reqsData, &reqs)
	if err != nil {
		return
	}
	for i := 0; i < len(reqs); i++ {
		reqs[i].Normalize()
	}
	return
}
