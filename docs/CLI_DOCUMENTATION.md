# HEPIC App Server v2 CLI Documentation

## Overview

HEPIC App Server v2 provides a comprehensive command-line interface built with Cobra CLI framework. The CLI offers advanced features for server management, configuration, health monitoring, and deployment.

## Installation

### Build from Source

```bash
git clone <repository-url>
cd hepic-app-server
go build -o hepic-app-server-v2 .
```

### Docker

```bash
docker build -t hepic-app-server-v2 .
```

## Global Commands

### Basic Usage

```bash
hepic-app-server-v2 [command] [flags]
```

### Global Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--config` | Config file path | `config.json` |
| `--log-level` | Log level (debug, info, warn, error) | `info` |
| `--log-format` | Log format (json, text) | `json` |
| `-v, --verbose` | Verbose output | `false` |
| `--version` | Show version information | - |

## Commands

### 1. Serve Command

Start the HEPIC App Server with ClickHouse integration.

```bash
hepic-app-server-v2 serve [flags]
```

#### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-p, --port` | Port to listen on | `8080` |
| `-H, --host` | Host to bind to | `0.0.0.0` |
| `--db-host` | ClickHouse host | `localhost` |
| `--db-port` | ClickHouse port | `9000` |
| `--db-user` | ClickHouse user | `default` |
| `--db-password` | ClickHouse password | - |
| `--db-database` | ClickHouse database | `hepic_analytics` |
| `--db-compress` | Enable ClickHouse compression | `true` |
| `--jwt-secret` | JWT secret key | - |
| `--jwt-expire-hours` | JWT token expiration in hours | `24` |

#### Examples

```bash
# Start server with default settings
hepic-app-server-v2 serve

# Start server on custom port and host
hepic-app-server-v2 serve --port 9090 --host 127.0.0.1

# Start with custom ClickHouse settings
hepic-app-server-v2 serve \
  --db-host clickhouse.example.com \
  --db-port 9440 \
  --db-user analytics \
  --db-password secret123 \
  --db-database hepic_prod

# Start with custom JWT settings
hepic-app-server-v2 serve \
  --jwt-secret "your-super-secret-key" \
  --jwt-expire-hours 48

# Start with debug logging
hepic-app-server-v2 serve --log-level debug --log-format text
```

### 2. Config Command

Configuration management commands.

```bash
hepic-app-server-v2 config [command]
```

#### Subcommands

##### Validate Configuration

```bash
hepic-app-server-v2 config validate [flags]
```

**Flags:**
- `--check-db` - Check ClickHouse connectivity

**Examples:**
```bash
# Validate configuration file
hepic-app-server-v2 config validate

# Validate with database connectivity check
hepic-app-server-v2 config validate --check-db

# Validate custom config file
hepic-app-server-v2 config validate --config /path/to/config.yaml
```

##### Show Configuration

```bash
hepic-app-server-v2 config show [flags]
```

**Flags:**
- `--show-secrets` - Show sensitive data (passwords, secrets)

**Examples:**
```bash
# Show current configuration
hepic-app-server-v2 config show

# Show configuration with secrets
hepic-app-server-v2 config show --show-secrets
```

##### Generate Configuration

```bash
hepic-app-server-v2 config generate [flags]
```

**Flags:**
- `--format` - Output format (json, yaml, env, docker) | `json`
- `--output` - Output directory | `.`

**Examples:**
```bash
# Generate JSON configuration
hepic-app-server-v2 config generate

# Generate YAML configuration
hepic-app-server-v2 config generate --format yaml

# Generate environment variables
hepic-app-server-v2 config generate --format env

# Generate Docker Compose
hepic-app-server-v2 config generate --format docker

# Generate to specific directory
hepic-app-server-v2 config generate --output /etc/hepic-app-server
```

### 3. Health Command

Health check commands for monitoring and diagnostics.

```bash
hepic-app-server-v2 health [command]
```

#### Subcommands

##### Health Check

```bash
hepic-app-server-v2 health check [flags]
```

**Flags:**
- `--timeout` - Health check timeout | `10s`
- `--verbose` - Verbose output

**Examples:**
```bash
# Basic health check
hepic-app-server-v2 health check

# Health check with timeout
hepic-app-server-v2 health check --timeout 30s

# Verbose health check
hepic-app-server-v2 health check --verbose
```

##### Health Server

```bash
hepic-app-server-v2 health server [flags]
```

**Flags:**
- `--port` - Health server port | `8081`
- `--host` - Health server host | `0.0.0.0`

**Examples:**
```bash
# Start health check server
hepic-app-server-v2 health server

# Start on custom port
hepic-app-server-v2 health server --port 9091

# Start on specific host
hepic-app-server-v2 health server --host 127.0.0.1
```

**Health Endpoints:**
- `GET /health` - Basic health check
- `GET /health/ready` - Readiness check
- `GET /health/live` - Liveness check
- `GET /health/detailed` - Detailed health information

### 4. Version Command

Show version information.

```bash
hepic-app-server-v2 version [flags]
```

**Flags:**
- `--verbose` - Show verbose version information
- `--json` - Output in JSON format

**Examples:**
```bash
# Show basic version
hepic-app-server-v2 version

# Show verbose version with features
hepic-app-server-v2 version --verbose

# Show version in JSON format
hepic-app-server-v2 version --json
```

## Configuration Files

### JSON Configuration

```json
{
  "database": {
    "host": "localhost",
    "port": 9000,
    "user": "default",
    "password": "",
    "database": "hepic_analytics",
    "sslmode": "disable",
    "compress": true
  },
  "server": {
    "port": "8080",
    "host": "0.0.0.0"
  },
  "jwt": {
    "secret": "your-super-secret-jwt-key-here-change-in-production",
    "expire_hours": 24
  },
  "logging": {
    "level": "info"
  }
}
```

### YAML Configuration

```yaml
database:
  host: localhost
  port: 9000
  user: default
  password: ""
  database: hepic_analytics
  sslmode: disable
  compress: true

server:
  port: "8080"
  host: "0.0.0.0"

jwt:
  secret: "your-super-secret-jwt-key-here-change-in-production"
  expire_hours: 24

logging:
  level: info
```

### Environment Variables

```bash
# Database Configuration (ClickHouse)
HEPIC_DATABASE_HOST=localhost
HEPIC_DATABASE_PORT=9000
HEPIC_DATABASE_USER=default
HEPIC_DATABASE_PASSWORD=
HEPIC_DATABASE_DATABASE=hepic_analytics
HEPIC_DATABASE_SSLMODE=disable
HEPIC_DATABASE_COMPRESS=true

# Server Configuration
HEPIC_SERVER_PORT=8080
HEPIC_SERVER_HOST=0.0.0.0

# JWT Configuration
HEPIC_JWT_SECRET=your-super-secret-jwt-key-here-change-in-production
HEPIC_JWT_EXPIRE_HOURS=24

# Logging
HEPIC_LOGGING_LEVEL=info
```

## Docker Integration

### Docker Compose

```yaml
version: '3.8'

services:
  clickhouse:
    image: clickhouse/clickhouse-server:latest
    container_name: hepic-clickhouse
    environment:
      CLICKHOUSE_DB: hepic_analytics
      CLICKHOUSE_USER: default
      CLICKHOUSE_PASSWORD: ""
    ports:
      - "9000:9000"
      - "8123:8123"
    volumes:
      - clickhouse_data:/var/lib/clickhouse
    networks:
      - hepic-network

  hepic-app-server:
    build: .
    container_name: hepic-app-server-v2
    ports:
      - "8080:8080"
    environment:
      - HEPIC_DATABASE_HOST=clickhouse
      - HEPIC_DATABASE_PORT=9000
      - HEPIC_DATABASE_USER=default
      - HEPIC_DATABASE_PASSWORD=
      - HEPIC_DATABASE_DATABASE=hepic_analytics
      - HEPIC_SERVER_PORT=8080
      - HEPIC_SERVER_HOST=0.0.0.0
      - HEPIC_JWT_SECRET=your-super-secret-jwt-key-here-change-in-production
      - HEPIC_JWT_EXPIRE_HOURS=24
      - HEPIC_LOGGING_LEVEL=info
    depends_on:
      - clickhouse
    networks:
      - hepic-network

volumes:
  clickhouse_data:

networks:
  hepic-network:
    driver: bridge
```

### Docker Run

```bash
# Run with environment variables
docker run -d \
  --name hepic-app-server-v2 \
  -p 8080:8080 \
  -e HEPIC_DATABASE_HOST=clickhouse \
  -e HEPIC_DATABASE_PORT=9000 \
  -e HEPIC_DATABASE_DATABASE=hepic_analytics \
  -e HEPIC_JWT_SECRET=your-super-secret-jwt-key \
  hepic-app-server-v2:latest serve

# Run with config file
docker run -d \
  --name hepic-app-server-v2 \
  -p 8080:8080 \
  -v /path/to/config.json:/app/config.json \
  hepic-app-server-v2:latest serve --config /app/config.json
```

## Production Deployment

### Systemd Service

```ini
[Unit]
Description=HEPIC App Server v2
After=network.target

[Service]
Type=simple
User=hepic
Group=hepic
WorkingDirectory=/opt/hepic-app-server
ExecStart=/opt/hepic-app-server/hepic-app-server-v2 serve
Restart=always
RestartSec=5
Environment=HEPIC_LOGGING_LEVEL=info
Environment=HEPIC_LOGGING_FORMAT=json

[Install]
WantedBy=multi-user.target
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hepic-app-server-v2
spec:
  replicas: 3
  selector:
    matchLabels:
      app: hepic-app-server-v2
  template:
    metadata:
      labels:
        app: hepic-app-server-v2
    spec:
      containers:
      - name: hepic-app-server-v2
        image: hepic-app-server-v2:latest
        ports:
        - containerPort: 8080
        env:
        - name: HEPIC_DATABASE_HOST
          value: "clickhouse-service"
        - name: HEPIC_DATABASE_PORT
          value: "9000"
        - name: HEPIC_DATABASE_DATABASE
          value: "hepic_analytics"
        - name: HEPIC_JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: hepic-secrets
              key: jwt-secret
        command: ["hepic-app-server-v2", "serve"]
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

## Monitoring and Logging

### Log Levels

- `debug` - Detailed debug information
- `info` - General information (default)
- `warn` - Warning messages
- `error` - Error messages only

### Log Formats

- `json` - Structured JSON logging (default)
- `text` - Human-readable text logging

### Health Monitoring

```bash
# Check server health
curl http://localhost:8080/health

# Check readiness
curl http://localhost:8080/health/ready

# Check liveness
curl http://localhost:8080/health/live

# Get detailed health information
curl http://localhost:8080/health/detailed
```

## Troubleshooting

### Common Issues

1. **Configuration not found**
   ```bash
   # Check config file location
   hepic-app-server-v2 config show
   
   # Generate example config
   hepic-app-server-v2 config generate
   ```

2. **ClickHouse connection failed**
   ```bash
   # Test ClickHouse connectivity
   hepic-app-server-v2 health check --check-db
   
   # Check ClickHouse status
   docker ps | grep clickhouse
   ```

3. **Port already in use**
   ```bash
   # Use different port
   hepic-app-server-v2 serve --port 9090
   
   # Check port usage
   netstat -tlnp | grep :8080
   ```

4. **JWT secret not configured**
   ```bash
   # Set JWT secret
   export HEPIC_JWT_SECRET="your-super-secret-key"
   
   # Or use command line flag
   hepic-app-server-v2 serve --jwt-secret "your-super-secret-key"
   ```

### Debug Mode

```bash
# Enable debug logging
hepic-app-server-v2 serve --log-level debug --log-format text

# Verbose health check
hepic-app-server-v2 health check --verbose

# Show detailed version information
hepic-app-server-v2 version --verbose
```

## Advanced Usage

### Custom Configuration

```bash
# Use custom config file
hepic-app-server-v2 serve --config /etc/hepic-app-server/production.json

# Override specific settings
hepic-app-server-v2 serve \
  --config /etc/hepic-app-server/config.json \
  --port 9090 \
  --db-host clickhouse-prod.example.com
```

### Health Check Server

```bash
# Start dedicated health check server
hepic-app-server-v2 health server --port 8081

# Use with load balancer
curl http://localhost:8081/health/ready
```

### Configuration Validation

```bash
# Validate configuration before starting
hepic-app-server-v2 config validate --check-db

# Show current configuration
hepic-app-server-v2 config show
```

This CLI provides comprehensive management capabilities for HEPIC App Server v2, making it easy to deploy, configure, and monitor in any environment.
