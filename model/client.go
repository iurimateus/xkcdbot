package model

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client exports 'client' a wrapper struct of a http.Client pointer with a base (non changing) url
// field that'll be shared across different requests.
var Client client

// client is a wrapper struct of a http.Client pointer with a base (non changing) url field that'll
// be shared across different requests.
type client struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
}

func init() {
	Client.BaseURL = &url.URL{
		Scheme: "http",
		Host:   "xkcd.com",
	}
	Client.HTTPClient = &http.Client{}
	Client.HTTPClient.Timeout = 5 * time.Second
}

func (c *client) newRequest(method, path string) (*http.Request, error) {
	relPath := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(relPath)

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	// prevents connection from being reused; should stop EOF errors.
	// see: https://stackoverflow.com/a/34474535
	req.Close = true

	return req, nil
}

// Do sends a request to the API, decoding its response body into a generic parameter 'v'.
// Argument should be a pointer
func (c *client) do(req *http.Request, v interface{}, backoff bool) (*http.Response, error) {
	var (
		res *http.Response
		err error
	)

	if backoff {
		res, err = c.doWithBackoff(req)
	} else {
		res, err = c.HTTPClient.Do(req)
	}

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(v)
	return res, err
}

// Wraps 'HTTPClient.Do' around a exponential backoff implementation
func (c *client) doWithBackoff(req *http.Request) (*http.Response, error) {
	var (
		res     *http.Response
		err     error
		retries uint
	)

	knownErrors := []string{"Timeout exceeded", "connection reset by peer"}
	known := func(err error) bool {
		for _, s := range knownErrors {
			if strings.Contains(err.Error(), s) {
				return true
			}
		}
		return false
	}

	const maxRetries = 5
	for retries <= maxRetries {
		res, err = c.HTTPClient.Do(req)
		if err == nil {
			break
		}

		if !known(err) {
			// Unknown error --> stop calling the API
			log.Printf("unknown error:\n%v\n", err)
			break
		}

		// (2 ^ retries) * 5 in minutes
		// last  wait time: 160 minutes
		// total wait time: 310 minutes / 5h10min (worst case)
		exp := 2 << retries
		sleepTime := time.Duration(exp) * (5 * time.Minute)

		log.Printf("Retry: %d, sleeping for %s\n", retries+1, sleepTime.String())
		time.Sleep(sleepTime)

		retries++
	}

	return res, err
}
