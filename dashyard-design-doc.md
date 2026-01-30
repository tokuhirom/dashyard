# Dashyard Design Document

## Overview

Dashyard は Prometheus メトリクスを可視化するための軽量なダッシュボードツール。YAML で定義したダッシュボードをディレクトリに配置するだけで、シンプルな Web UI からメトリクスを閲覧できる。

### Goals

- **シンプル**: Grafana のような多機能は不要。メトリクスが見れればいい
- **Dashboard as Code**: YAML ファイルでダッシュボードを定義、Git 管理可能
- **軽量**: Docker イメージ一つで動作、外部依存なし
- **最小限の認証**: Basic 認証 (将来的に OIDC)

### Non-Goals

- アラーティング機能
- ダッシュボードの GUI エディタ
- マルチテナント / 細かい認可制御
- Prometheus 以外のデータソース対応

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                      Docker Container                    │
│  ┌───────────────────────────────────────────────────┐  │
│  │                   Dashyard Server                  │  │
│  │  ┌─────────────┐  ┌─────────────┐  ┌───────────┐  │  │
│  │  │   Web UI    │  │  REST API   │  │   Auth    │  │  │
│  │  │  (embed.FS) │  │             │  │ Middleware│  │  │
│  │  └─────────────┘  └─────────────┘  └───────────┘  │  │
│  │         │               │                │        │  │
│  │         └───────────────┼────────────────┘        │  │
│  │                         │                         │  │
│  │  ┌─────────────────────────────────────────────┐  │  │
│  │  │              Dashboard Loader               │  │  │
│  │  │         (YAML files → in-memory)            │  │  │
│  │  └─────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────┘  │
│                            │                             │
│              ┌─────────────┴─────────────┐              │
│              ▼                           ▼              │
│    /etc/dashyard/config.yaml      /dashboards/*.yaml    │
│         (mount: ro)                  (mount: ro)        │
└─────────────────────────────────────────────────────────┘
                            │
                            ▼
                    ┌───────────────┐
                    │  Prometheus   │
                    │    Server     │
                    └───────────────┘
```

## Configuration

### config.yaml

```yaml
# /etc/dashyard/config.yaml
server:
  port: 8080
  # host: "0.0.0.0"  # default

prometheus:
  url: "http://prometheus:9090"
  # timeout: 30s  # default

dashboards_dir: "/dashboards"

auth:
  users:
    - id: admin
      password: "$6$rounds=5000$salt$hashedpassword..."
    - id: viewer
      password: "$6$..."

  # 将来的に OIDC 対応
  # oidc:
  #   issuer: "https://accounts.google.com"
  #   client_id: "xxx"
  #   client_secret: "xxx"
  #   redirect_uri: "http://localhost:8080/callback"
```

### Dashboard YAML

```yaml
# /dashboards/infra/network.yaml
title: "Network Overview"

rows:
  - title: "Traffic"
    panels:
      - type: graph
        title: "Inbound"
        query: "rate(node_network_receive_bytes_total[5m])"
        unit: bytes

      - type: graph
        title: "Outbound"
        query: "rate(node_network_transmit_bytes_total[5m])"
        unit: bytes

      - type: graph
        title: "Errors"
        query: "rate(node_network_receive_errs_total[5m])"
        unit: count

  - title: null  # タイトルなしの row
    panels:
      - type: markdown
        content: |
          ## 運用メモ
          - しきい値: 100MB/s 超えたら要注意
          - 担当: @infra-team

  - title: "Connections"
    panels:
      - type: graph
        title: "TCP Connections"
        query: "node_netstat_Tcp_CurrEstab"
        unit: count
```

### Dashboard Directory Structure

```
/dashboards/
├── infra/
│   ├── network.yaml
│   └── storage.yaml
├── apps/
│   ├── api-server.yaml
│   └── batch.yaml
└── overview.yaml
```

UI のサイドバーはディレクトリ構造をそのまま反映:

```
📁 infra
   📊 network
   📊 storage
📁 apps
   📊 api-server
   📊 batch
📊 overview
```

## Data Model

### Config

```go
type Config struct {
    Server      ServerConfig     `yaml:"server"`
    Prometheus  PrometheusConfig `yaml:"prometheus"`
    DashboardsDir string         `yaml:"dashboards_dir"`
    Auth        AuthConfig       `yaml:"auth"`
}

type ServerConfig struct {
    Port int    `yaml:"port"`
    Host string `yaml:"host"`
}

type PrometheusConfig struct {
    URL     string        `yaml:"url"`
    Timeout time.Duration `yaml:"timeout"`
}

type AuthConfig struct {
    Users []User `yaml:"users"`
}

type User struct {
    ID       string `yaml:"id"`
    Password string `yaml:"password"` // SHA-512 crypt format
}
```

### Dashboard

```go
type Dashboard struct {
    Title string `yaml:"title"`
    Rows  []Row  `yaml:"rows"`
}

type Row struct {
    Title  *string `yaml:"title"` // nil = タイトルなし
    Panels []Panel `yaml:"panels"`
}

type Panel struct {
    Type    PanelType `yaml:"type"`    // graph, markdown
    Title   string    `yaml:"title"`   // graph のみ
    Query   string    `yaml:"query"`   // graph のみ
    Unit    Unit      `yaml:"unit"`    // graph のみ: bytes, percent, count
    Content string    `yaml:"content"` // markdown のみ
    Width   int       `yaml:"width"`   // optional: 相対幅 (default: 1)
}

type PanelType string

const (
    PanelTypeGraph    PanelType = "graph"
    PanelTypeMarkdown PanelType = "markdown"
)

type Unit string

const (
    UnitBytes   Unit = "bytes"   // 自動で KB/MB/GB に変換
    UnitPercent Unit = "percent" // 0-100%
    UnitCount   Unit = "count"   // そのまま表示
)
```

## API Design

### Authentication

```
POST /api/login
Content-Type: application/json

{
  "id": "admin",
  "password": "secret"
}

Response:
Set-Cookie: session=<signed-token>; HttpOnly; Secure; SameSite=Strict
{
  "ok": true
}
```

### Dashboard List

```
GET /api/dashboards

Response:
{
  "dashboards": [
    {
      "path": "infra/network",
      "title": "Network Overview"
    },
    {
      "path": "infra/storage",
      "title": "Storage Metrics"
    },
    {
      "path": "apps/api-server",
      "title": "API Server"
    },
    {
      "path": "overview",
      "title": "Overview"
    }
  ],
  "tree": {
    "infra": {
      "network": { "title": "Network Overview" },
      "storage": { "title": "Storage Metrics" }
    },
    "apps": {
      "api-server": { "title": "API Server" }
    },
    "overview": { "title": "Overview" }
  }
}
```

### Dashboard Detail

```
GET /api/dashboards/{path}
# path: "infra/network" など

Response:
{
  "title": "Network Overview",
  "rows": [
    {
      "title": "Traffic",
      "panels": [
        {
          "type": "graph",
          "title": "Inbound",
          "query": "rate(node_network_receive_bytes_total[5m])",
          "unit": "bytes"
        }
      ]
    }
  ]
}
```

### Query Prometheus

```
GET /api/query?query={promql}&start={unix}&end={unix}&step={seconds}

Response:
{
  "status": "success",
  "data": {
    "resultType": "matrix",
    "result": [
      {
        "metric": { "__name__": "up", "job": "prometheus" },
        "values": [
          [1234567890, "1"],
          [1234567900, "1"]
        ]
      }
    ]
  }
}
```

## UI Design

### Layout

```
┌─────────────────────────────────────────────────────────────────┐
│  🏠 Dashyard                              [1h] [6h] [24h] [7d]  │
├───────────────┬─────────────────────────────────────────────────┤
│               │                                                 │
│  📁 infra     │   Network Overview                              │
│     📊 network│   ─────────────────────────────────────────────│
│     📊 storage│                                                 │
│               │   ┌─ Traffic ─────────────────────────────────┐ │
│  📁 apps      │   │ [Inbound] [Outbound] [Errors]             │ │
│     📊 api    │   └───────────────────────────────────────────┘ │
│     📊 batch  │                                                 │
│               │   ┌───────────────────────────────────────────┐ │
│  📊 overview  │   │ ## 運用メモ                                │ │
│               │   │ - しきい値: 100MB/s ...                   │ │
│               │   └───────────────────────────────────────────┘ │
│               │                                                 │
│               │   ┌─ Connections ─────────────────────────────┐ │
│               │   │ [TCP Connections]                         │ │
│               │   └───────────────────────────────────────────┘ │
│               │                                                 │
└───────────────┴─────────────────────────────────────────────────┘
```

### Time Range Selector

- プリセット: 1h, 6h, 24h, 7d
- URL パラメータで保持: `?from=now-1h&to=now`

### Panel Rendering

**Graph Panel**:
- 折れ線グラフ
- Y軸は unit に応じて自動フォーマット
- ホバーで値表示
- 複数系列対応 (query が複数 metric を返す場合)

**Markdown Panel**:
- GitHub Flavored Markdown
- コードハイライトなし (シンプルに)

### Row Layout

- Row 内のパネルは flexbox で横並び
- 幅が足りなければ折り返し
- Panel の `width` で相対幅を指定可能 (default: 1)

## Tech Stack

### Backend

- **Language**: Go 1.23+
- **Web Framework**: net/http (標準ライブラリ)
- **Config**: gopkg.in/yaml.v3
- **Password Verification**: golang.org/x/crypto (SHA-512 crypt)
- **Session**: gorilla/securecookie or 自前実装

### Frontend

- **Framework**: React or Svelte (検討中)
- **Charting**: Chart.js or Recharts
- **Markdown**: marked or remark
- **Build**: Vite
- **Embed**: Go の embed.FS で single binary 化

### Infrastructure

- **Container**: Docker (Alpine base)
- **Build**: Multi-stage Dockerfile

## Deployment

### Docker Compose

```yaml
version: "3.8"
services:
  dashyard:
    image: ghcr.io/dashyard/dashyard:latest
    ports:
      - "8080:8080"
    volumes:
      - ./config.yaml:/etc/dashyard/config.yaml:ro
      - ./dashboards:/dashboards:ro
    environment:
      - DASHYARD_CONFIG=/etc/dashyard/config.yaml
```

### Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dashyard
spec:
  replicas: 1
  template:
    spec:
      containers:
        - name: dashyard
          image: ghcr.io/dashyard/dashyard:latest
          ports:
            - containerPort: 8080
          volumeMounts:
            - name: config
              mountPath: /etc/dashyard
              readOnly: true
            - name: dashboards
              mountPath: /dashboards
              readOnly: true
      volumes:
        - name: config
          configMap:
            name: dashyard-config
        - name: dashboards
          configMap:
            name: dashyard-dashboards
```

## Security Considerations

### Authentication

- パスワードは SHA-512 crypt (mkpasswd -m sha-512 互換)
- セッションは signed cookie (HMAC-SHA256)
- Cookie 属性: HttpOnly, Secure (HTTPS時), SameSite=Strict

### Network

- Prometheus への接続は内部ネットワークを想定
- HTTPS 終端はリバースプロキシ (nginx, Traefik) で行う

### Input Validation

- PromQL クエリはそのまま Prometheus に渡す (Prometheus 側で検証)
- Dashboard path は英数字、ハイフン、スラッシュのみ許可
- YAML パースエラーは起動時に検出、UI に表示

## Development Roadmap

### Phase 1: MVP

- [ ] Config loader
- [ ] Dashboard loader (YAML → in-memory)
- [ ] Basic auth (SHA-512 crypt)
- [ ] REST API
- [ ] Frontend (React + Chart.js)
- [ ] Dockerfile
- [ ] Basic documentation

### Phase 2: Polish

- [ ] Hot reload (ファイル監視)
- [ ] YAML validation with better error messages
- [ ] Panel width control
- [ ] Multiple queries per panel
- [ ] Legend customization

### Phase 3: Enterprise Features

- [ ] OIDC authentication
- [ ] Dashboard variables / templating
- [ ] Dashboard embedding (iframe)
- [ ] Prometheus basic auth support

## Open Questions

1. **Frontend framework**: React vs Svelte
   - React: エコシステムが大きい、Chart.js / Recharts が使いやすい
   - Svelte: バンドルサイズが小さい、シンプル

2. **Session storage**: Cookie only vs Cookie + server-side
   - Cookie only: stateless でスケールしやすい
   - Server-side: revocation が簡単

3. **Dashboard reload**: Hot reload vs manual
   - Hot reload: 開発体験がいい
   - Manual: シンプル、予期しない変更を防げる

## References

- [Prometheus HTTP API](https://prometheus.io/docs/prometheus/latest/querying/api/)
- [Grafana Dashboard JSON Model](https://grafana.com/docs/grafana/latest/dashboards/build-dashboards/view-dashboard-json-model/)
- [Perses Dashboard Spec](https://perses.dev/)
