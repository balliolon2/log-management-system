-- Database Schema & RLS Strategy
-- 1. สร้าง Extension ที่จำเป็น
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 2. ตาราง Users (สำหรับ AuthN/AuthZ)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'viewer', -- 'admin' หรือ 'viewer'
    tenant_id TEXT NOT NULL, -- Key สำคัญสำหรับแยกข้อมูลลูกค้า
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Seed Data (ข้อมูลตัวอย่างสำหรับ Login)
-- Password hashes (bcrypt):
-- admin123 -> $2a$10$vHDZ/YqHWtf8jxBNATtAy.xYKbO7769jy2Yutc5UuGjL6HY4HgOBu
-- viewer123 -> $2a$10$y.zfB8a0paF62rFCbqCT6erpitlf/av1WEzPb5wkgcYBjE8I0VJHO
INSERT INTO users (username, password_hash, role, tenant_id)
VALUES
    ('admin', '$2a$10$vHDZ/YqHWtf8jxBNATtAy.xYKbO7769jy2Yutc5UuGjL6HY4HgOBu', 'admin', 'all'),
    ('viewer1', '$2a$10$y.zfB8a0paF62rFCbqCT6erpitlf/av1WEzPb5wkgcYBjE8I0VJHO', 'viewer', 'demo'),
    ('viewer2', '$2a$10$y.zfB8a0paF62rFCbqCT6erpitlf/av1WEzPb5wkgcYBjE8I0VJHO', 'viewer', 'demoA')
ON CONFLICT (username) DO NOTHING;

-- 3. ตาราง Logs (Schema กลาง)
CREATE TABLE IF NOT EXISTS logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id TEXT NOT NULL, -- Must-have สำหรับ RLS
    timestamp TIMESTAMPTZ NOT NULL, -- @timestamp
    source TEXT NOT NULL, -- firewall, api, etc.
    event_type TEXT,
    severity INT,

    -- Field ที่เหลือเก็บเป็น JSONB เพื่อความยืดหยุ่น (Schema-less)
    body JSONB,

    -- เก็บ Raw message เผื่อไว้ debug
    raw_message TEXT
);

-- 4. Indexes เพื่อความเร็วในการค้นหา
-- Index สำหรับค้นหาใน JSONB (สำคัญมากสำหรับ ClickHouse/Postgres search)
CREATE INDEX IF NOT EXISTS idx_logs_body ON logs USING GIN (body);
-- Index สำหรับ Filter ตามเวลาและ Tenant
CREATE INDEX IF NOT EXISTS idx_logs_tenant_time ON logs (tenant_id, timestamp DESC);

-- 5. Row Level Security (RLS) - หัวใจของ Multi-tenant
ALTER TABLE logs ENABLE ROW LEVEL SECURITY;

-- สร้าง Policy: "User จะเห็นเฉพาะแถวที่ tenant_id ตรงกับค่า session variable"
CREATE POLICY tenant_isolation_policy ON logs
    FOR SELECT
    USING (tenant_id = current_setting('app.current_tenant', true)::text);

-- อนุญาตให้ Insert ได้ (Ingestion service จะเป็นคน insert)
CREATE POLICY ingest_policy ON logs
    FOR INSERT
    WITH CHECK (true);

-- 6. ตาราง Alert Rules (กฎการแจ้งเตือน)
CREATE TABLE IF NOT EXISTS alert_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    rule_name TEXT NOT NULL UNIQUE,
    description TEXT,
    event_type TEXT, -- เช่น "login_failed"
    condition_field TEXT, -- field ที่ต้องตรวจสอบ เช่น "src_ip"
    condition_operator TEXT, -- เช่น "count_gt" (มากกว่า)
    threshold INT, -- เกณฑ์ เช่น 5
    time_window_seconds INT, -- ช่วงเวลา เช่น 300 (5 นาที)
    severity INT DEFAULT 5, -- ระดับความรุนแรง 1-10
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 7. ตาราง Alert History (ประวัติ alerts ที่เกิดขึ้น)
CREATE TABLE IF NOT EXISTS alert_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id TEXT NOT NULL,
    rule_id UUID REFERENCES alert_rules(id),
    rule_name TEXT NOT NULL,
    triggered_at TIMESTAMPTZ DEFAULT NOW(),
    details JSONB, -- เก็บรายละเอียดเพิ่มเติม เช่น IP, user ที่เกี่ยวข้อง
    severity INT,
    acknowledged BOOLEAN DEFAULT false
);

-- Index สำหรับ Alert History
CREATE INDEX IF NOT EXISTS idx_alert_history_tenant_time
    ON alert_history (tenant_id, triggered_at DESC);

CREATE INDEX IF NOT EXISTS idx_alert_history_rule
    ON alert_history (rule_id);

-- 8. RLS สำหรับ Alert History
ALTER TABLE alert_history ENABLE ROW LEVEL SECURITY;

CREATE POLICY alert_tenant_isolation_policy ON alert_history
    FOR SELECT
    USING (tenant_id = current_setting('app.current_tenant', true)::text);

CREATE POLICY alert_insert_policy ON alert_history
    FOR INSERT
    WITH CHECK (true);

-- 9. Seed Alert Rules (Hardcode กฎเริ่มต้น)
INSERT INTO alert_rules (rule_name, description, event_type, condition_field, condition_operator, threshold, time_window_seconds, severity)
VALUES
    ('brute_force_login', 'Detect brute force login attempts', 'login_failed', 'src_ip', 'count_gt', 5, 300, 8),
    ('repeated_malware', 'Detect repeated malware detection', 'malware_detected', 'host', 'count_gt', 3, 600, 9)
ON CONFLICT (rule_name) DO NOTHING;
