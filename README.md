# HEPIC App Server v2

Advanced REST API Server based on Echo v4 with PostgreSQL connection and JWT authentication.

## 🚀 Features

- **REST API** based on Echo v4
- **PostgreSQL** database with automatic initialization
- **JWT authentication** with role-based system
- **Swagger documentation** with auto-generation
- **Middleware** for CORS, logging, security
- **Docker** support
- **Graceful shutdown**
- **Pagination** and data filtering
- **Input validation**
- **Clean Architecture** with services and handlers separation

## 📋 Requirements

- Go 1.21+
- PostgreSQL 12+
- Docker (optional)

## 🛠 Quick Start

### Local Installation

1. **Clone and navigate to directory:**
```bash
cd /home/shurik/Projects/hepic-app-server
```

2. **Install dependencies:**
```bash
go mod tidy
```

3. **Setup database:**
```bash
# Create PostgreSQL database
createdb hepic_db
```

4. **Configure environment:**
```bash
cp env.example .env
# Edit .env with your settings
```

5. **Build and run:**
```bash
make build
make run
```

### Docker Installation

```bash
make docker
make docker-run
```

## 📚 Documentation

- **[Main Documentation](docs/README.md)** - Complete setup and usage guide
- **[Configuration Guide](docs/CONFIG_README.md)** - Configuration with Viper framework
- **[API Documentation](http://localhost:8080/api/v1/docs/)** - Swagger UI (after server start)

## 🔧 Configuration

The server supports multiple configuration sources:

1. **Environment variables** (highest priority)
2. **Configuration files** (JSON, YAML)
3. **Default values** (lowest priority)

See [Configuration Guide](docs/CONFIG_README.md) for detailed information.

## 🔐 Authentication

### Get Token

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'
```

### Use Token

```bash
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 📊 Main API Endpoints

### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/logout` - User logout
- `GET /api/v1/auth/me` - Current user info

### Users
- `GET /api/v1/users` - List users (with pagination)
- `GET /api/v1/users/{id}` - Get user by ID
- `POST /api/v1/users` - Create user (admin only)
- `PUT /api/v1/users/{id}` - Update user
- `DELETE /api/v1/users/{id}` - Delete user (admin only)

### HEP Records
- `GET /api/v1/hep` - List HEP records (with filtering and pagination)
- `GET /api/v1/hep/{id}` - Get HEP record by ID
- `POST /api/v1/hep` - Create HEP record
- `GET /api/v1/hep/stats` - HEP records statistics

## 🐳 Docker

### Build Image
```bash
make docker
```

### Run Container
```bash
make docker-run
```

## 🧪 Development

```bash
# Install dependencies
make deps

# Build
make build

# Run
make run

# Run in development mode
make dev

# Run tests
make test

# Clean
make clean
```

## 🏗 Project Structure

```
├── config/          # Configuration
├── database/        # Database connection
├── docs/           # Documentation
├── handlers/       # HTTP handlers (controllers)
├── middleware/     # Middleware
├── models/         # Data models
├── routes/         # API routes
├── services/       # Business logic (services)
├── main.go         # Entry point
├── Dockerfile      # Docker image
├── Makefile        # Build commands
└── config.json     # Configuration
```

## 🎯 Architecture

The project follows **Clean Architecture** principles with layer separation:

- **Handlers** - HTTP controllers, handle requests and responses
- **Services** - Business logic, data operations
- **Models** - Data structures and DTOs
- **Database** - Data access layer
- **Routes** - API route configuration
- **Middleware** - Cross-cutting concerns (CORS, auth, logging)

## 📄 License

MIT License

## 👥 Authors

HEPIC Development Team

## 📞 Support

- Email: support@hepic.local
- Documentation: http://localhost:8080/api/v1/docs/
- Issues: GitHub Issues

---

**Version:** v2.0.0  
**Go Version:** 1.21+  
**Echo Version:** v4.13.4