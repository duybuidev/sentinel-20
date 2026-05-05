# Sentinel-20: Enterprise Observability Ecosystem

> Production-grade monitoring, logging & alerting platform for 20-container microservices architecture deployed on Azure.

## Tech Stack

| Layer | Technology |
|---|---|
| **API** | Go + Gin |
| **Database** | PostgreSQL + Redis |
| **Logs** | Promtail → Loki |
| **Metrics** | cAdvisor → Prometheus |
| **Dashboard** | Grafana |
| **Gateway** | Nginx |
| **CI/CD** | GitHub Actions |
| **Frontend** | ReactJS + TailwindCSS |

## Architecture

\`\`\`
Internet → Nginx (Gateway)
               ├── /api    → Sentinel API (Go)  → PostgreSQL
               │                                → Redis
               ├── /ui     → React Frontend
               └── /grafana→ Grafana Dashboard

Log Pipeline:  Promtail → Loki → Grafana
Metrics:       cAdvisor → Prometheus → Grafana
Alerting:      Grafana → Slack / PagerDuty
\`\`\`

## Quick Start

\`\`\`bash
git clone git@github.com:<your-username>/sentinel-20.git
cd sentinel-20
cp .env.example .env
# Edit .env với credentials của bạn
docker compose up -d
\`\`\`

## Services

| Service | Port | Language |
|---|---|---|
| Sentinel API | 8080 | Go |
| Frontend | 3000 | React |
| Grafana | 3001 | - |
| Prometheus | 9090 | - |
| PostgreSQL | 5432 | - |
| Redis | 6379 | - |
| Order Service | 8081 | Go |
| Auth Service | 8083 | Java |
| Shipping Service | 8091-8093 | Python |

## Roadmap
- [x] Mock microservices ecosystem
- [ ] Sentinel API (Go + Gin)
- [ ] PostgreSQL + Redis integration
- [ ] Loki + Prometheus + Grafana stack
- [ ] React Dashboard
- [ ] GitHub Actions CI/CD
- [ ] Slack Alerting
