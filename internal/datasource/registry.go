package datasource

import (
	"fmt"
	"sort"

	"github.com/tokuhirom/dashyard/internal/config"
	"github.com/tokuhirom/dashyard/internal/prometheus"
)

// Registry manages named Prometheus clients created from datasource configs.
type Registry struct {
	clients     map[string]*prometheus.Client
	defaultName string
}

// NewRegistry creates a Registry from the given datasource configurations.
func NewRegistry(datasources []config.DatasourceConfig) *Registry {
	clients := make(map[string]*prometheus.Client, len(datasources))
	var defaultName string

	for _, ds := range datasources {
		clients[ds.Name] = prometheus.NewClient(ds.URL, ds.Timeout)
		if ds.Default {
			defaultName = ds.Name
		}
	}

	return &Registry{
		clients:     clients,
		defaultName: defaultName,
	}
}

// Get returns the Prometheus client for the given datasource name.
// If name is empty, the default datasource is returned.
func (r *Registry) Get(name string) (*prometheus.Client, error) {
	if name == "" {
		return r.Default(), nil
	}
	client, ok := r.clients[name]
	if !ok {
		return nil, fmt.Errorf("unknown datasource %q", name)
	}
	return client, nil
}

// Default returns the default Prometheus client.
func (r *Registry) Default() *prometheus.Client {
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
