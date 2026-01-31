package prometheus

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is an HTTP client for the Prometheus query_range API.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Prometheus client.
func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
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

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("executing request: %w", err)
	}

	return resp.Body, resp.StatusCode, nil
}
