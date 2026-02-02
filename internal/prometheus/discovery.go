package prometheus

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

// promAPIResponse is the JSON envelope for Prometheus HTTP API responses.
type promAPIResponse struct {
	Status string          `json:"status"`
	Data   json.RawMessage `json:"data"`
	Error  string          `json:"error,omitempty"`
}

// MetricMetadataEntry represents a single metadata entry from /api/v1/metadata.
type MetricMetadataEntry struct {
	Type string `json:"type"`
	Help string `json:"help"`
	Unit string `json:"unit"`
}

// MetricInfo aggregates all discovered information about a metric.
type MetricInfo struct {
	Name        string
	Type        string
	Help        string
	Unit        string
	Labels      []string
	LabelValues map[string][]string // label name -> values
}

// doGet performs a GET request to the given Prometheus API path, applies auth,
// and decodes the standard JSON envelope. It returns the raw data field.
func (c *Client) doGet(ctx context.Context, path string, params url.Values) (json.RawMessage, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("parsing base URL: %w", err)
	}
	u = u.JoinPath(path)
	if params != nil {
		u.RawQuery = params.Encode()
	}

	requestURL := u.String()
	slog.Debug("prometheus API request", "url", requestURL)

	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	c.applyAuth(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request (%s): %w", requestURL, err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body (%s): %w", requestURL, err)
	}

	var apiResp promAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("decoding response (status %d, url %s): %w", resp.StatusCode, requestURL, err)
	}

	if apiResp.Status != "success" {
		return nil, fmt.Errorf("prometheus API error (%s): %s", requestURL, apiResp.Error)
	}

	return apiResp.Data, nil
}

// MetricNames returns all metric names from /api/v1/label/__name__/values.
func (c *Client) MetricNames(ctx context.Context) ([]string, error) {
	data, err := c.doGet(ctx, "/api/v1/label/__name__/values", nil)
	if err != nil {
		return nil, fmt.Errorf("fetching metric names: %w", err)
	}

	var names []string
	if err := json.Unmarshal(data, &names); err != nil {
		return nil, fmt.Errorf("decoding metric names: %w", err)
	}
	return names, nil
}

// MetricMetadata returns metric metadata from /api/v1/metadata.
func (c *Client) MetricMetadata(ctx context.Context) (map[string][]MetricMetadataEntry, error) {
	data, err := c.doGet(ctx, "/api/v1/metadata", nil)
	if err != nil {
		return nil, fmt.Errorf("fetching metadata: %w", err)
	}

	var metadata map[string][]MetricMetadataEntry
	if err := json.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("decoding metadata: %w", err)
	}
	return metadata, nil
}

// MetricLabels returns label names for a given metric from /api/v1/labels,
// filtering out the __name__ label.
func (c *Client) MetricLabels(ctx context.Context, metricName string) ([]string, error) {
	params := url.Values{}
	params.Set("match[]", fmt.Sprintf(`{__name__="%s"}`, metricName))

	data, err := c.doGet(ctx, "/api/v1/labels", params)
	if err != nil {
		return nil, fmt.Errorf("fetching labels for %s: %w", metricName, err)
	}

	var labels []string
	if err := json.Unmarshal(data, &labels); err != nil {
		return nil, fmt.Errorf("decoding labels for %s: %w", metricName, err)
	}

	// Filter out __name__
	filtered := make([]string, 0, len(labels))
	for _, l := range labels {
		if l != "__name__" {
			filtered = append(filtered, l)
		}
	}
	return filtered, nil
}

// MetricLabelValues returns the values of a specific label for a given metric
// using /api/v1/label/{label}/values?match[]={__name__="metric"}.
func (c *Client) MetricLabelValues(ctx context.Context, metricName, labelName string) ([]string, error) {
	params := url.Values{}
	params.Set("match[]", metricName)

	path := fmt.Sprintf("/api/v1/label/%s/values", labelName)
	data, err := c.doGet(ctx, path, params)
	if err != nil {
		return nil, fmt.Errorf("fetching label values for %s/%s: %w", metricName, labelName, err)
	}

	var values []string
	if err := json.Unmarshal(data, &values); err != nil {
		return nil, fmt.Errorf("decoding label values for %s/%s: %w", metricName, labelName, err)
	}
	return values, nil
}
