package datasource

import (
	"fmt"
	"sort"

	"github.com/tokuhirom/dashyard/internal/config"
	"github.com/tokuhirom/dashyard/internal/prometheus"
)

// Registry manages named datasource clients created from datasource configs.
type Registry struct {
	clients     map[string]Datasource
	defaultName string
}

// NewRegistry creates a Registry from the given datasource configurations.
func NewRegistry(datasources []config.DatasourceConfig) (*Registry, error) {
	clients := make(map[string]Datasource, len(datasources))
	var defaultName string

	for _, ds := range datasources {
		switch ds.Type {
		case "prometheus":
			var opts []prometheus.ClientOption
			if len(ds.Headers) > 0 {
				opts = append(opts, prometheus.WithHeaders(ds.Headers))
			}
			clients[ds.Name] = prometheus.NewClient(ds.URL, ds.Timeout, opts...)
		default:
			return nil, fmt.Errorf("unsupported datasource type %q for %q", ds.Type, ds.Name)
		}
		if ds.Default {
			defaultName = ds.Name
		}
	}

	return &Registry{
		clients:     clients,
		defaultName: defaultName,
	}, nil
}

// Get returns the datasource for the given name.
// If name is empty, the default datasource is returned.
func (r *Registry) Get(name string) (Datasource, error) {
	if name == "" {
		return r.Default(), nil
	}
	client, ok := r.clients[name]
	if !ok {
		return nil, fmt.Errorf("unknown datasource %q", name)
	}
	return client, nil
}

// Default returns the default datasource.
func (r *Registry) Default() Datasource {
	return r.clients[r.defaultName]
}

// DefaultName returns the name of the default datasource.
func (r *Registry) DefaultName() string {
	return r.defaultName
}

// Names returns all registered datasource names in sorted order.
func (r *Registry) Names() []string {
	names := make([]string, 0, len(r.clients))
	for name := range r.clients {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
