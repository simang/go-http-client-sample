package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	HeaderAPIKey = "api-key"
)

// Client internal client interface
type Client interface {
	NewRequest(method, path string, body interface{}) (*http.Request, error)
	Do(*http.Request, interface{}) (*http.Response, error)
}

// APIClient is implementation of Client
type APIClient struct {
	BaseURL    *url.URL
	APIKey     string
	httpClient *http.Client
}

// NewRequest returns http.Request
func (c *APIClient) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	//fmt.Println(buf)

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set(HeaderAPIKey, c.APIKey)
	return req, nil
}

// Do returns response
func (c *APIClient) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}

// SetTimeout sets http client timeout (sec)
func (c *APIClient) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = time.Duration(timeout * time.Second)
}

// NewAPIClient returns new client instance
func NewAPIClient(httpClient *http.Client, baseURL string, apiKey string) (*APIClient, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	httpClient.Timeout = time.Duration(30 * time.Second)

	bu, err := url.Parse(strings.TrimRight(baseURL, "/"))
	if err != nil {
		return nil, err
	}
	c := &APIClient{
		httpClient: httpClient,
		BaseURL:    bu,
		APIKey:     apiKey,
	}
	return c, nil
}
