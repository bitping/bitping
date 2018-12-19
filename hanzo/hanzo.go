package hanzo

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
)

type Client struct {
	ApiKey   string
	Endpoint string

	client *http.Client
	ctx    context.Context
}

func New(apiKey string, client *http.Client) *Client {
	return &Client{
		ApiKey:   apiKey,
		Endpoint: "https://api.hanzo.io",
		client:   client,
	}
}

func (c *Client) Request(method, url string, body interface{}, dst interface{}) (*http.Response, error) {
	var data *bytes.Buffer

	// Encode body
	if body != nil {
		if json, err := json.Marshal(body); err != nil {
			log.Printf("Failed to marshal request body: %v", err)
			return nil, err
		} else {
			data = bytes.NewBuffer(json)
		}
	} else {
		data = bytes.NewBufferString("")
	}

	// Create request
	req, err := http.NewRequest(method, c.Endpoint+url, data)
	if err != nil {
		log.Fatalf("Failed to create hanzo request: %v", err)
		return nil, err
	}

	// Set headers
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth("", c.ApiKey)

	// Do request
	r, err := c.client.Do(req)

	dump, _ := httputil.DumpRequest(req, true)
	log.Printf("hanzo request:\n%s", dump)

	// Request failed
	if err != nil {
		log.Printf("hanzo request failed: %v", err)
		return r, err
	}

	dump, _ = httputil.DumpResponse(r, true)
	log.Printf("hanzo response:\n%s", dump)

	defer r.Body.Close()

	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to convert io.ReaderCloser: %v", err)
		return nil, err
	}

	// Decode response wrapper
	if err = json.Unmarshal(body, dst); err != nil {
		log.Printf("Failed to decode response: %v", err)
		return nil, err
	}

	return r, err
}
