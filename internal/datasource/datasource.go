package datasource

import (
	"context"
	"io"
)

// Datasource defines the interface for querying a metrics backend.
// Handlers use this interface so they are not coupled to a specific implementation.
type Datasource interface {
	QueryRange(ctx context.Context, query, start, end, step string) (io.ReadCloser, int, error)
	Ping(ctx context.Context) error
	LabelValues(ctx context.Context, label, match string) (io.ReadCloser, int, error)
}
