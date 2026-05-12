# Sentinel-20

Nền tảng **quan sát & vận hành** (observability) cho kiến trúc nhiều container: API Sentinel (Go), PostgreSQL, Redis, pipeline log/metrics qua Loki/Prometheus/Grafana, gateway Nginx và CI/CD GitHub Actions.

---

## Tiến độ triển khai

| # | Nội dung | Trạng thái |
|---|----------|:----------:|
| 1 | PostgreSQL + Redis | ✅ |
| 2 | Sentinel API foundation | ✅ |
| 3 | Repository layer | ✅ |
| 4 | Service + Handler layer | ✅ |
| 5 | SysWatch Agent | ✅ |
| 6 | Observability stack (Loki, Prometheus, Grafana, Promtail, cAdvisor, …) | ✅ |
| 7 | Mock ecosystem (8 services) | ✅ |
| 8 | Nginx API gateway (`/api`, `/ui`, `/grafana`, …) | ✅ |
| 9 | CI/CD GitHub Actions | ✅ |
| 10 | Alerting (Slack / Grafana contact points) | ⬜ |
| 11 | Frontend React | ⬜ |
| 12 | Swagger UI | ⬜ |

**Ghi chú**

- **Mục 11:** Trong repo đã có skeleton UI React ([`web/`](./web/) — Vite + TypeScript + Tailwind, Docker + Nginx **`/ui/`**). Đổi ô thành ✅ trên GitHub khi bạn coi milestone frontend đã xong.
- **Mục 10:** Grafana và biến `DISCORD_WEBHOOK_URL` có thể dùng cho thông báo; tích hợp **Slack** theo roadmap vẫn có thể đánh dấu hoàn thành riêng khi bạn cấu hình xong.
- **Mục 12:** API documentation tương tác (Swagger/OpenAPI) — chưa có trong repo.

---

## Tech stack

| Layer | Công nghệ |
|-------|-----------|
| **API** | Go + Gin |
| **Database** | PostgreSQL, Redis |
| **Logs** | Promtail → Loki → Grafana |
| **Metrics** | cAdvisor → Prometheus → Grafana |
| **Dashboard vận hành** | Grafana |
| **Gateway** | Nginx |
| **Frontend** | React (Vite) + Tailwind CSS |
| **CI/CD** | GitHub Actions |

---

## Kiến trúc

```
Internet → Nginx (port 80)
               ├── /api      → Sentinel API (Go)     → PostgreSQL
               │                                         → Redis
               ├── /ui/      → React frontend (static)
               ├── /grafana/ → Grafana
               └── /prometheus/ → Prometheus (UI)

Pipeline log:     Promtail → Loki → Grafana
Pipeline metrics: cAdvisor → Prometheus → Grafana
```

---

## Quick start (Docker)

```bash
git clone git@github.com:<your-username>/sentinel-20.git
cd sentinel-20
cp .env.example .env
# Chỉnh POSTGRES_*, REDIS_*, GRAFANA_*, DISCORD_WEBHOOK_URL (nếu dùng) trong .env
docker compose up -d --build
```

- **API (qua gateway):** `http://<host>/api/v1/...`
- **Health:** `http://<host>/health`
- **Frontend:** `http://<host>/ui/`
- **Grafana:** `http://<host>/grafana/` (compose map host **3001** → Grafana nếu truy cập trực tiếp container)

---

## Phát triển frontend (local)

Cần API đang chạy (ví dụ `docker compose up sentinel-api` hoặc `go run` trong `api/` trên port **8080**).

```bash
cd web
npm install
npm run dev
```

Mở **`http://localhost:3000/ui/`** (có dấu `/` cuối). Vite proxy `/api` và `/health` sang `http://127.0.0.1:8080`.

Build production:

```bash
cd web && npm run build
```

---

## Dịch vụ & cổng (tham khảo)

| Dịch vụ | Cổng (host) | Ghi chú |
|---------|-------------|---------|
| Nginx gateway | **80** | Điểm vào chính cho `/api`, `/ui`, `/grafana` |
| Sentinel API | **8080** | Truy cập trực tiếp container; production nên qua `/api` |
| Frontend (dev) | **3000** | Chỉ khi `npm run dev`; production qua `/ui/` |
| Grafana | **3001** | Map từ compose; có thể vào qua gateway `/grafana/` |
| Prometheus | **9090** | Có route `/prometheus/` qua gateway |
| PostgreSQL | **5432** | |
| Redis | **6379** | |
| Loki | **3100** | |
| Mock services | nội bộ | 8 service trong `mock-ecosystem` |

---

## Cấu trúc thư mục (rút gọn)

```
sentinel-20/
├── api/                 # Sentinel API (Go)
├── web/                 # React frontend (Vite)
├── syswatch-agent/      # Agent đồng bộ container → API
├── nginx/               # Cấu hình gateway
├── observability/       # Prometheus, Loki, Grafana provisioning, …
├── mock-ecosystem/      # Mock microservices
├── scripts/             # SQL init, …
└── .github/workflows/   # CI/CD
```

---

## License

Thêm license của bạn tại đây (ví dụ MIT) nếu repo public.
