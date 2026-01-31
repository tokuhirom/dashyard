# Dashboard YAML Format

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
        chart_type: line            # line, bar, area, scatter
        legend: "{label_name}"      # legend template
        y_min: 0                    # optional y-axis bounds
        y_max: 100
        stacked: false              # stack series
        y_scale: linear             # linear (default) or log
        thresholds:                 # optional reference lines
          - value: 80
            color: orange
            label: "Warning"
      - title: "Notes"
        type: markdown
        content: |
          Markdown content here.
        full_width: true            # span entire row width (markdown only)
```

# Rules

## PromQL by Metric Type

- **counter** (`_total` suffix): always wrap with `rate(...[5m])` or `increase(...[5m])`. Never display a raw counter.
  - `rate()` converts a counter into a per-second rate. The resulting unit depends on what the counter measures:
    - `rate(network_io_bytes_total[5m])` produces bytes/sec — use `unit: bytes`
    - `rate(http_requests_total[5m])` produces requests/sec — use `unit: count`
- **histogram** (`_bucket` suffix): use `histogram_quantile(0.99, rate(...[5m]))`
- **gauge**: use directly, or apply `avg()`, `sum()`, `min()`, `max()`
- **summary** (`_sum`/`_count`): `rate(sum[5m]) / rate(count[5m])` for average

## Ratio and Percent Metrics

- Metrics named `_ratio` are typically 0–1 scale. Multiply by 100 in the PromQL expression when using `unit: percent` (e.g. `system_cpu_utilization_ratio * 100`), or use as-is with a suitable `y_max`.
- Check the actual metric range if unsure — some exporters use 0–100 even for `_ratio` names.

## Filtering and Grouping

- Filter with `{label="value"}`, group with `by (label)`
- Use `$variable` to reference dashboard variables in queries

## Unit Selection

| Unit | When to use | Y-axis behavior |
|------|------------|-----------------|
| `bytes` | memory, disk, network I/O | Human-readable (KB, MB, GB) |
| `percent` | utilization, ratios (0–100 scale) | Fixed 0–100 unless `y_min`/`y_max` override |
| `seconds` | durations, latencies | Human-readable time |
| `count` | counts, rates, dimensionless (default when omitted) | SI suffixes (k, M) |

## Division Queries

When dividing two `rate()` expressions (e.g. hit rate = hits / (hits + misses)), guard against division by zero with `> 0` on the denominator. Without this, periods with no traffic produce `NaN`.

```
rate(hits_total[5m]) / (rate(hits_total[5m]) + rate(misses_total[5m]) > 0) * 100
```

The `> 0` filter drops zero-denominator samples so the panel shows no data instead of `NaN`.

## Stacked Charts

Use `stacked: true` when series represent parts of a whole that sum to a meaningful total (e.g. memory by state: used + cached + free + buffers = total). Do not stack independent series that overlap (e.g. CPU utilization per core).
