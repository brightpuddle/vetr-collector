package aci

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// Client is an HTTP ACI API client.
// Use aci.NewClient to initiate a client.
// This will ensure proper cookie handling and processing of modifiers.
type Client struct {
	// HTTPClient is the *http.Client used for API requests.
	HTTPClient *http.Client
	// host is the APIC IP or hostname, e.g. 10.0.0.1:80 (port is optional).
	host string
	// Usr is the APIC username.
	Usr string
	// Pwd is the APIC password.
	Pwd string
	// LastRefresh is the timestamp of the last token refresh interval.
	LastRefresh time.Time
	// Token is the current authentication token
	Token string
}

// NewClient creates a new ACI HTTP client.
// Pass modifiers in to modify the behavior of the client, e.g.
//  client, _ := NewClient("apic", "user", "password", RequestTimeout(120))
func NewClient(url, usr, pwd string, mods ...func(*Client)) (Client, error) {

	// Normalize the URL
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	cookieJar, _ := cookiejar.New(nil)
	httpClient := http.Client{
		Timeout:   300 * time.Second,
		Transport: tr,
		Jar:       cookieJar,
	}

	client := Client{
		HTTPClient: &httpClient,
		host:       url,
		Usr:        usr,
		Pwd:        pwd,
	}
	for _, mod := range mods {
		mod(&client)
	}
	return client, nil
}

// NewReq creates a new Req request for this client.
func (client Client) NewReq(method, uri string, body io.Reader, mods ...func(*Req)) Req {
	httpReq, err := http.NewRequest(method, client.host+":443"+uri+".json", body)
	if err != nil {
		panic(err)
	}
	req := Req{
		HTTPReq: httpReq,
		Refresh: true,
	}
	for _, mod := range mods {
		mod(&req)
	}
	return req
}

// RequestTimeout modifies the HTTP request timeout from the default of 60 seconds.
func RequestTimeout(x time.Duration) func(*Client) {
	return func(client *Client) {
		client.HTTPClient.Timeout = x * time.Second
	}
}

// Do makes a request.
// Requests for Do are built ouside of the client, e.g.
//
//  req := client.NewReq("GET", "/api/class/fvBD", nil)
//  res := client.Do(req)
func (client *Client) Do(req Req) (Res, error) {
	if req.Refresh && time.Now().Sub(client.LastRefresh) > 480*time.Second {
		if err := client.Refresh(); err != nil {
			return Res{}, err
		}
	}

	httpRes, err := client.HTTPClient.Do(req.HTTPReq)
	if err != nil {
		return Res{}, err
	}
	defer httpRes.Body.Close()

	body, err := ioutil.ReadAll(httpRes.Body)
	if err != nil {
		return Res{}, errors.New("cannot decode response body")
	}

	res := Res(gjson.ParseBytes(body))

	if httpRes.StatusCode == 400 {
		errStr := res.Get("imdata.0.error.attributes.text").Str
		if strings.Contains(errStr, "Unable to process the query, result dataset is too big") {
			return Res{}, errors.New("result dataset is too big")
		}
	}

	if httpRes.StatusCode != http.StatusOK {
		return Res{}, fmt.Errorf("received HTTP status %d", httpRes.StatusCode)
	}

	return res, nil
}

// Get makes a GET request and returns a GJSON result.
// Results will be the raw data structure as returned by the APIC, wrapped in imdata, e.g.
//
//  {
//    "imdata": [
//      {
//        "fvTenant": {
//          "attributes": {
//            "dn": "uni/tn-mytenant",
//            "name": "mytenant",
//          }
//        }
//      }
//    ],
//    "totalCount": "1"
//  }
func (client *Client) Get(path string, mods ...func(*Req)) (Res, error) {
	req := client.NewReq("GET", path, nil, mods...)
	res, err := client.Do(req)
	// if err != nil && err.Error() == "result dataset is too big" && len(mods) == 0 {
	// 	res, err = client.GetWithPagination(path, mods...)
	// }
	return res, err
}

// GetWithPagination performs a get request with pagination
func (client *Client) GetWithPagination(path string, mods ...func(*Req)) (Res, error) {
	pageSize := 10
	pageNumber := 0
	mods = append(mods, Query("page", strconv.Itoa(pageNumber)))
	mods = append(mods, Query("page-size", strconv.Itoa(pageSize)))
	req := client.NewReq("GET", path, nil, mods...)
	res, err := client.Do(req)

	if err != nil {
		return res, err
	}
	if !res.Get("imdata").IsArray() {
		return res, errors.New("imdata is not an array")
	}

	totalCount := res.Get("totalCount").Str
	count, err := strconv.Atoi(res.Get("totalCount").Str)
	if err != nil {
		return res, err
	}

	var tmp string
	for i, value := range res.Get("imdata").Array() {
		if i == 0 {
			tmp = value.Raw
		} else {
			tmp = tmp + "," + value.Raw
		}
	}

	count = count - pageSize
	for count > 0 {
		pageNumber = pageNumber + 1
		mods = append(mods, Query("page", strconv.Itoa(pageNumber)))
		mods = append(mods, Query("page-size", strconv.Itoa(pageSize)))
		req := client.NewReq("GET", path, nil, mods...)
		res, err := client.Do(req)
		if err != nil {
			return res, err
		}
		if !res.Get("imdata").IsArray() {
			return res, errors.New("imdata is not an array")
		}
		for _, value := range res.Get("imdata").Array() {
			tmp = tmp + "," + value.Raw
		}
		count = count - pageSize
	}

	json := fmt.Sprintf(`{"totalCount":%s,"imdata":[%s]}`, totalCount, tmp)
	return gjson.Parse(json), err
}

// GetClass makes a GET request by class and unwraps the results.
// Result is removed from imdata, but still wrapped in Class.attributes, e.g.
//  [
//    {
//      "fvTenant": {
//        "attributes": {
//          "dn": "uni/tn-mytenant",
//          "name": "mytenant",
//        }
//      }
//    }
//  ]
func (client *Client) GetClass(class string, mods ...func(*Req)) (Res, error) {
	res, err := client.Get(fmt.Sprintf("/api/class/%s", class), mods...)
	if err != nil {
		return res, err
	}
	return res.Get("imdata"), nil
}

// GetDn makes a GET request by DN.
// Result is removed from imdata and first result is removed from the list, e.g.
//  {
//    "fvTenant": {
//      "attributes": {
//        "dn": "uni/tn-mytenant",
//        "name": "mytenant",
//      }
//    }
//  }
func (client *Client) GetDn(dn string, mods ...func(*Req)) (Res, error) {
	res, err := client.Get(fmt.Sprintf("/api/mo/%s", dn), mods...)
	if err != nil {
		return res, err
	}
	return res.Get("imdata.0"), nil
}

// Post makes a POST request and returns a GJSON result.
// Hint: Use the Body struct to easily create POST body data.
func (client *Client) Post(path, data string, mods ...func(*Req)) (Res, error) {
	req := client.NewReq("POST", path, strings.NewReader(data), mods...)
	return client.Do(req)
}

// Login authenticates to the APIC.
func (client *Client) Login() error {
	data := fmt.Sprintf(`{"aaaUser":{"attributes":{"name":"%s","pwd":"%s"}}}`,
		client.Usr,
		client.Pwd,
	)
	res, err := client.Post("/api/aaaLogin", data, NoRefresh)
	if err != nil {
		return err
	}
	errText := res.Get("imdata.0.error.attributes.text").Str
	if errText != "" {
		return errors.New("authentication error")
	}
	client.Token = res.Get("imdata.0.aaaLogin.attributes.token").Str
	client.LastRefresh = time.Now()
	return nil
}

// Refresh refreshes the authentication token.
// Note that this will be handled automatically be default.
// Refresh will be checked every request and the token will be refreshed after 8 minutes.
// Pass aci.NoRefresh to prevent automatic refresh handling and handle it directly instead.
func (client *Client) Refresh() error {
	res, err := client.Get("/api/aaaRefresh", NoRefresh)
	if err != nil {
		return err
	}
	client.Token = res.Get("imdata.0.aaaRefresh.attributes.token").Str
	client.LastRefresh = time.Now()
	return nil
}
