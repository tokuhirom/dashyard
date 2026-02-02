package prometheus

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/tokuhirom/dashyard/internal/metrics"
)

// ClientOption configures optional Client settings.
type ClientOption func(*Client)

// Header represents a single HTTP header as a name/value pair.
type Header struct {
	Name  string
	Value string
}

// WithHeaders sets custom HTTP headers to include in every request.
func WithHeaders(headers []Header) ClientOption {
	return func(c *Client) {
		c.headers = append(c.headers, headers...)
	}
}

// Client is an HTTP client for the Prometheus query_range API.
type Client struct {
	baseURL    string
	httpClient *http.Client
	headers    []Header
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
	for _, h := range c.headers {
		req.Header.Add(h.Name, h.Value)
	}
}

// QueryRange performs a Prometheus query_range request and returns the raw response body.
// The caller is responsible for closing the returned ReadCloser.
func (c *Client) QueryRange(ctx context.Context, query, start, end, step string) (io.ReadCloser, int, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, 0, fmt.Errorf("parsing base URL: %w", err)
	}
	u = u.JoinPath("api/v1/query_range")

	params := url.Values{}
	params.Set("query", query)
	params.Set("start", start)
	params.Set("end", end)
	params.Set("step", step)
	u.RawQuery = params.Encode()

	slog.Debug("prometheus query_range request", "url", u.String())

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, 0, fmt.Errorf("creating request: %w", err)
	}
	c.applyAuth(req)

	reqStart := time.Now()
	resp, err := c.httpClient.Do(req)
	duration := time.Since(reqStart).Seconds()
	metrics.DatasourceQueryDuration.Observe(duration)
	if err != nil {
		metrics.DatasourceQueryTotal.WithLabelValues("error").Inc()
		return nil, 0, fmt.Errorf("executing request: %w", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		metrics.DatasourceQueryTotal.WithLabelValues("success").Inc()
	} else {
		metrics.DatasourceQueryTotal.WithLabelValues("error").Inc()
	}

	return resp.Body, resp.StatusCode, nil
}

// Ping checks whether the Prometheus server is reachable by hitting the /-/ready endpoint.
func (c *Client) Ping(ctx context.Context) error {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return fmt.Errorf("parsing base URL: %w", err)
	}
	u = u.JoinPath("-/ready")

	slog.Debug("prometheus ping request", "url", u.String())

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	c.applyAuth(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	return nil
}

// LabelValues queries the Prometheus label values API and returns the raw response body.
// The caller is responsible for closing the returned ReadCloser.
func (c *Client) LabelValues(ctx context.Context, label, match string) (io.ReadCloser, int, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, 0, fmt.Errorf("parsing base URL: %w", err)
	}
	u = u.JoinPath("api/v1/label", label, "values")

	if match != "" {
		params := url.Values{}
		params.Set("match[]", match)
		u.RawQuery = params.Encode()
	}

	slog.Debug("prometheus label values request", "url", u.String())

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
