-- Create database for HEPIC analytics
CREATE DATABASE IF NOT EXISTS hepic_analytics;

-- Use the database
USE hepic_analytics;

-- Create HEP analytics table
CREATE TABLE IF NOT EXISTS hep_analytics (
    id UInt64,
    call_id String,
    source_ip IPv4,
    destination_ip IPv4,
    protocol String,
    method String,
    status_code UInt16,
    timestamp DateTime64(3),
    raw_data String,
    created_at DateTime64(3) DEFAULT now64(3)
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (timestamp, call_id)
SETTINGS index_granularity = 8192;

-- Create materialized view for real-time statistics
CREATE MATERIALIZED VIEW IF NOT EXISTS hep_stats_mv
ENGINE = SummingMergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (timestamp, protocol, method, status_code)
AS SELECT
    toStartOfMinute(timestamp) as timestamp,
    protocol,
    method,
    status_code,
    count() as count
FROM hep_analytics
GROUP BY timestamp, protocol, method, status_code;

-- Create table for user analytics
CREATE TABLE IF NOT EXISTS user_analytics (
    user_id UInt64,
    action String,
    timestamp DateTime64(3),
    ip_address IPv4,
    user_agent String,
    created_at DateTime64(3) DEFAULT now64(3)
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (timestamp, user_id)
SETTINGS index_granularity = 8192;

-- Create table for system metrics
CREATE TABLE IF NOT EXISTS system_metrics (
    metric_name String,
    metric_value Float64,
    timestamp DateTime64(3),
    tags Map(String, String),
    created_at DateTime64(3) DEFAULT now64(3)
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (timestamp, metric_name)
SETTINGS index_granularity = 8192;
