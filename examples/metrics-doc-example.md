# Prometheus Metrics Reference for Dashyard

This document lists all available Prometheus metrics. Use it as context for generating Dashyard dashboard YAML files.

## Dashboard YAML Format

```yaml
title: "Dashboard Title"
variables:                          # optional
  - name: device                    # variable name used as $device in queries
    label: "Network Device"         # display label
    query: "label_values(metric_name, label_name)"
rows:
  - title: "Row Title"
    repeat: device                  # optional: repeat row for each variable value
    panels:
      - title: "Panel Title"
        type: graph                 # "graph" or "markdown"
        query: 'promql_expression'  # required for graph
        unit: bytes                 # bytes, percent, count, seconds
        chart_type: line            # line, bar, area, scatter, pie, doughnut
        legend: "{label_name}"      # legend template
        y_min: 0                    # optional y-axis bounds
        y_max: 100
        stacked: false              # stack series
        thresholds:                 # optional reference lines
          - value: 80
            color: orange
            label: "Warning"
      - title: "Notes"
        type: markdown
        content: |
          Markdown content here.
```

## PromQL Guidelines

- **counter** metrics: always wrap with `rate(...[5m])` or `increase(...[5m])`
- **histogram** `_bucket` metrics: use `histogram_quantile(0.99, rate(...[5m]))`
- **gauge** metrics: use directly, or apply aggregation like `avg()`, `sum()`
- **summary** `_sum`/`_count` metrics: `rate(sum[5m]) / rate(count[5m])` for average
- Use `{label="value"}` for filtering, `by (label)` for grouping
- Use `$variable` syntax to reference dashboard variables in queries

## Unit Selection

- `bytes` — metrics measuring bytes (memory, disk, network I/O)
- `percent` — ratios and utilization (0-100 scale)
- `seconds` — durations and latencies
- `count` — counts, rates, and dimensionless values

## Available Metrics

### system_cpu

**`system_cpu_load_average_1m_ratio`** (gauge)
  1-minute CPU load average.

**`system_cpu_utilization_ratio`** (gauge)
  CPU utilization as a ratio between 0 and 1.
  Labels: `cpu`

### system_disk

**`system_disk_io_bytes_total`** (counter)
  Total disk I/O bytes by device and direction.
  Labels: `device`, `direction`

### system_memory

**`system_memory_usage_bytes`** (gauge)
  Memory usage in bytes by state.
  Labels: `state`

### system_network

**`system_network_io_bytes_total`** (counter)
  Total network I/O bytes by device and direction.
  Labels: `device`, `direction`

