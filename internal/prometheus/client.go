package prometheus

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// ClientOption configures optional Client settings.
type ClientOption func(*Client)

// WithBearerToken sets a bearer token for authentication.
func WithBearerToken(token string) ClientOption {
	return func(c *Client) {
		c.bearerToken = token
	}
}

// Client is an HTTP client for the Prometheus query_range API.
type Client struct {
	baseURL     string
	httpClient  *http.Client
	bearerToken string
}

// NewClient creates a new Prometheus client.
func NewClient(baseURL string, timeout time.Duration, opts ...ClientOption) *Client {
	c := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) applyAuth(req *http.Request) {
	if c.bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.bearerToken)
	}
}

// QueryRange performs a Prometheus query_range request and returns the raw response body.
// The caller is responsible for closing the returned ReadCloser.
func (c *Client) QueryRange(ctx context.Context, query, start, end, step string) (io.ReadCloser, int, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, 0, fmt.Errorf("parsing base URL: %w", err)
	}
	u.Path = "/api/v1/query_range"

	params := url.Values{}
	params.Set("query", query)
	params.Set("start", start)
	params.Set("end", end)
	params.Set("step", step)
	u.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, 0, fmt.Errorf("creating request: %w", err)
	}
	c.applyAuth(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("executing request: %w", err)
	}

	return resp.Body, resp.StatusCode, nil
}

// LabelValues queries the Prometheus label values API and returns the raw response body.
// The caller is responsible for closing the returned ReadCloser.
func (c *Client) LabelValues(ctx context.Context, label, match string) (io.ReadCloser, int, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, 0, fmt.Errorf("parsing base URL: %w", err)
	}
	u.Path = fmt.Sprintf("/api/v1/label/%s/values", label)

	if match != "" {
		params := url.Values{}
		params.Set("match[]", match)
		u.RawQuery = params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, 0, fmt.Errorf("creating request: %w", err)
	}
	c.applyAuth(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("executing request: %w", err)
	}

	return resp.Body, resp.StatusCode, nil
}
