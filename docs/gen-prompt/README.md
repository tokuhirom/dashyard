# gen-prompt: Real Monitoring Stack

A docker-compose stack with real services generating real Prometheus metrics. Use it to produce realistic `gen-prompt` output and to preview LLM-generated dashboards in Dashyard.

## Architecture

```
traffic-gen ──► Traefik ──► whoami
                  │
                  ▼ (scrape :8080/metrics)
              Prometheus ◄── OTel Collector ──┬── hostmetrics
                  │                           └── Redis
                  ▼
              Dashyard (:8080)
```

## Services

| Service | Port | Role |
|---------|------|------|
| Prometheus | 9090 | Metrics storage, remote_write receiver |
| OTel Collector | — | Collects host + Redis metrics → Prometheus |
| Traefik | 8888 | HTTP proxy, exposes Prometheus metrics |
| Redis | 6379 | Cache (metrics collected by OTel) |
| whoami | — | Simple backend for Traefik |
| traffic-gen | — | wget loop generating HTTP traffic |
| Dashyard | 8080 | Dashboard viewer |

## Files

| File | Description |
|------|-------------|
| `docker-compose.yaml` | Full monitoring stack |
| `prometheus.yaml` | Prometheus config (scrapes Traefik, receives remote_write) |
| `otelcol-config.yaml` | OTel Collector config (hostmetrics + Redis → Prometheus) |
| `traefik.yaml` | Traefik v3 static config |
| `traefik-dynamic.yaml` | Traefik routing rules |
| `config.yaml` | Dashyard config (points to Prometheus) |
| `output/` | Generated files and dashboards |
| `output/prompt.md` | Static prompt template (editable) |
| `output/prompt-metrics.md` | Metrics + label values (regenerated) |
| `output/dashboards/` | LLM-generated dashboard YAML files |

## Quick Start: Generate Prompt

Run from the repository root:

```bash
make gen-prompt
```

This will:
1. Start Prometheus, OTel Collector, Traefik, Redis, whoami, and traffic-gen
2. Wait 60 seconds for metrics to accumulate
3. Run `gen-prompt` against Prometheus
4. Write `output/prompt.md` and `output/prompt-metrics.md`
5. Shut down all services

## View Dashboards

After generating dashboards with an LLM, place the YAML files in `output/dashboards/` and start the full stack:

```bash
make gen-prompt-up
```

Open http://localhost:8080 (login: admin / admin).

## Workflow

1. `make gen-prompt` — start stack → accumulate metrics → run gen-prompt → stop
2. Feed `output/prompt.md` + `output/prompt-metrics.md` to an LLM to generate dashboard YAML files
3. Place generated YAML files in `docs/gen-prompt/output/dashboards/`
4. `make gen-prompt-up` — start Dashyard + monitoring stack
5. Open http://localhost:8080 to verify dashboards render with real metrics
6. If the dashboards need improvement, refine the prompt and repeat from step 1
