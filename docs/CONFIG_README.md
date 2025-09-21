# HEPIC App Server v2 Configuration

## 🎯 Recommended Framework: **Viper**

For mapping JSON config to Go structures we use **Viper** - the most popular and powerful configuration framework for Go.

## 🚀 Viper Advantages:

- ✅ **Multiple formats**: JSON, YAML, TOML, ENV, etc.
- ✅ **Automatic reading** from files and environment variables
- ✅ **Validation** of configuration
- ✅ **Hot reloading** (optional)
- ✅ **Excellent documentation** and community
- ✅ **Used in large projects** (Docker, Kubernetes, etc.)

## 📁 Configuration Structure

```
├── config/
│   ├── config.go           # Основная конфигурация
│   └── config_viper.go     # Viper реализация
├── config.json             # JSON конфигурация
├── config.yaml             # YAML конфигурация
├── env.example             # Пример переменных окружения
└── main_viper_example.go   # Пример использования
```

## 🔧 Usage

### 1. JSON Configuration (config.json)
```json
{
  "database": {
    "host": "localhost",
    "port": 5432,
    "user": "hepic_user",
    "password": "hepic_password",
    "name": "hepic_db",
    "sslmode": "disable"
  },
  "server": {
    "port": "8080",
    "host": "0.0.0.0"
  },
  "jwt": {
    "secret": "your-super-secret-jwt-key-here",
    "expire_hours": 24
  },
  "logging": {
    "level": "info"
  }
}
```

### 2. YAML Configuration (config.yaml)
```yaml
database:
  host: localhost
  port: 5432
  user: hepic_user
  password: hepic_password
  name: hepic_db
  sslmode: disable

server:
  port: "8080"
  host: "0.0.0.0"

jwt:
  secret: "your-super-secret-jwt-key-here"
  expire_hours: 24

logging:
  level: info
```

### 3. Environment Variables
```bash
# Database
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=hepic_user
export DB_PASSWORD=hepic_password
export DB_NAME=hepic_db
export DB_SSLMODE=disable

# Server
export SERVER_PORT=8080
export SERVER_HOST=0.0.0.0

# JWT
export JWT_SECRET=your-super-secret-jwt-key-here
export JWT_EXPIRE_HOURS=24

# Logging
export LOG_LEVEL=info
```

## 🎯 Configuration Priority

Viper uses the following priority (from highest to lowest):

1. **Environment variables** (highest priority)
2. **Configuration file** (config.json, config.yaml, etc.)
3. **Default values** (lowest priority)

## 🔄 Hot Reloading

For automatic configuration reload:

```go
viper.WatchConfig()
viper.OnConfigChange(func(e fsnotify.Event) {
    log.Println("Config file changed:", e.Name)
    // Reload configuration
})
```

## 🛡️ Validation

Viper automatically validates:
- Required fields
- Data types
- Value ranges
- Formats (email, URL, etc.)

## 📊 Alternative Solutions

### 1. **Cleanenv** (Simple)
```bash
go get github.com/ilyakaznacheev/cleanenv
```

**Pros:**
- Very simple
- Minimal dependencies
- Good performance

**Cons:**
- Fewer features
- Limited flexibility

### 2. **Koanf** (Modern)
```bash
go get github.com/knadh/koanf/v2
```

**Pros:**
- Very flexible
- Excellent performance
- Modern API

**Cons:**
- Harder to learn
- Less documentation

### 3. **Envconfig** (ENV only)
```bash
go get github.com/kelseyhightower/envconfig
```

**Pros:**
- Only for ENV variables
- Very simple
- Fast

**Cons:**
- Only ENV variables
- Limited functionality

## 🚀 Recommendation

For **HEPIC App Server v2** we recommend **Viper** because:

1. **Time-tested** - used in large projects
2. **Multi-functional** - support for all formats
3. **Good ecosystem** - many examples and documentation
4. **Flexibility** - easily configurable for any needs
5. **Professional** - suitable for production

## 📝 Integration Example

```go
package main

import (
    "hepic-app-server/v2/config"
)

func main() {
    // Load configuration
    cfg := config.Load()
    
    // Usage
    log.Printf("Server starting on %s:%s", cfg.Server.Host, cfg.Server.Port)
    log.Printf("Database: %s@%s:%d/%s", cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
}
```

## 🔧 Installation

```bash
# Install Viper
go get github.com/spf13/viper

# Or update go.mod
go mod tidy
```

## 📚 Documentation

- [Viper GitHub](https://github.com/spf13/viper)
- [Viper Documentation](https://pkg.go.dev/github.com/spf13/viper)
- [Viper Examples](https://github.com/spf13/viper/tree/master/examples)

---

**Conclusion:** Viper is the best choice for configuration in Go projects! 🎯
