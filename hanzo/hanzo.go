package hanzo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"
)

type Client struct {
	ApiKey   string
	Endpoint string

	client *http.Client
	ctx    context.Context
}

func New(c context.Context, apiKey string) *Client {
	c, _ = context.WithTimeout(c, time.Second*30)
	client := urlfetch.Client(c)
	client.Transport = &http.Transport{
		Context: ctx,
	}

	return &Client{
		ApiKey:   apiKey,
		Endpoint: "https://api.hanzo.io",
		client:   client,
		ctx:      c,
	}
}

func (c *Client) Request(method, url string, body interface{}, dst interface{}) (*http.Response, error) {
	var data *bytes.Buffer

	// Encode body
	if body != nil {
		data = bytes.NewBuffer(json.EncodeBytes(body))
	} else {
		data = bytes.NewBufferString("")
	}

	// Create request
	req, err := http.NewRequest(method, c.Endpoint+url, data)
	if err != nil {
		fmt.Errorf("Failed to create hanzo request: %v", err)
		return nil, err
	}

	// Set headers
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth("", c.ApiKey)

	// Do request
	r, err := c.client.Do(req)

	dump, _ := httputil.DumpRequest(req, true)
	log.Warn("hanzo request:\n%s", dump, c.ctx)

	// Request failed
	if err != nil {
		log.Error("hanzo request failed: %v", err, c.ctx)
		return r, err
	}

	dump, _ = httputil.DumpResponse(r, true)
	log.Warn("hanzo response:\n%s", dump, c.ctx)

	defer r.Body.Close()

	// Decode response wrapper
	if err := json.Decode(r.Body, dst); err != nil {
		log.Warn("Failed to decode response:%v", err, c.ctx)
		return nil, err
	}

	return r, err
}
