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

// Client internal client interface
type Client interface {
	NewRequest(options *Options) (*http.Request, error)
	Do(*http.Request, interface{}) (*http.Response, error)
}

// APIClient is implementation of Client
type APIClient struct {
	BaseURL    *url.URL
	httpClient *http.Client
}

type Header struct {
	Key   string
	Value string
}

type Options struct {
	Method  string
	Path    string
	Headers []*Header
	Body    interface{}
}

// NewRequest returns http.Request
func (c *APIClient) NewRequest(options *Options) (*http.Request, error) {
	rel := &url.URL{Path: options.Path}
	u := c.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if options.Body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(options.Body)
		if err != nil {
			return nil, err
		}
	}

	//fmt.Println(buf)

	req, err := http.NewRequest(options.Method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if options.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	for _, header := range options.Headers {
		req.Header.Set(header.Key, header.Value)
	}
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

// NewIamClient returns new client instance
func NewAPIClient(httpClient *http.Client, baseURL string) (*APIClient, error) {
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
	}
	return c, nil
}
