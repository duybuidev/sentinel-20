-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Bảng containers: lưu trạng thái các container
CREATE TABLE IF NOT EXISTS containers (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(255) NOT NULL UNIQUE,
    image       VARCHAR(255) NOT NULL,
    status      VARCHAR(50)  NOT NULL DEFAULT 'unknown',
    -- running | exited | dead | unknown
    cpu_percent NUMERIC(5,2) DEFAULT 0,
    mem_usage   BIGINT       DEFAULT 0, -- bytes
    mem_limit   BIGINT       DEFAULT 0, -- bytes
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- Bảng logs: lưu toàn bộ log từ containers
CREATE TABLE IF NOT EXISTS logs (
    id           BIGSERIAL PRIMARY KEY,
    container_id UUID        REFERENCES containers(id) ON DELETE CASCADE,
    level        VARCHAR(20) NOT NULL DEFAULT 'INFO',
    -- INFO | WARNING | ERROR | CRITICAL
    message      TEXT        NOT NULL,
    service      VARCHAR(255),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Bảng events: lưu các sự kiện quan trọng (start/stop/crash)
CREATE TABLE IF NOT EXISTS events (
    id           BIGSERIAL PRIMARY KEY,
    container_id UUID        REFERENCES containers(id) ON DELETE CASCADE,
    type         VARCHAR(50) NOT NULL,
    -- started | stopped | crashed | restarted
    exit_code    INT,
    description  TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes để query nhanh
CREATE INDEX IF NOT EXISTS idx_logs_container_id  ON logs(container_id);
CREATE INDEX IF NOT EXISTS idx_logs_level         ON logs(level);
CREATE INDEX IF NOT EXISTS idx_logs_created_at    ON logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_events_container_id ON events(container_id);
CREATE INDEX IF NOT EXISTS idx_events_created_at  ON events(created_at DESC);

-- Auto-update updated_at
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN NEW.updated_at = NOW(); RETURN NEW; END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_containers_updated_at
    BEFORE UPDATE ON containers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
