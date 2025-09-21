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

## 📋 Requirements

- Go 1.21+
- PostgreSQL 12+
- ClickHouse 22.0+
- Docker (опционально)

## 🛠 Installation

### Local Installation

1. **Клонирование и переход в директорию:**
```bash
cd /home/shurik/Projects/hepic-app-server
```

2. **Установка зависимостей:**
```bash
go mod tidy
```

3. **Настройка базы данных:**
```bash
# Создание базы данных PostgreSQL
createdb hepic_db
```

4. **Настройка конфигурации:**
```bash
cp config.env.example config.env
# Отредактируйте config.env под ваши настройки
```

5. **Сборка и запуск:**
```bash
make build
make run
```

### Docker Installation

```bash
# Запуск с PostgreSQL и ClickHouse
docker-compose -f docker-compose.clickhouse.yml up -d

# Проверка статуса
docker-compose -f docker-compose.clickhouse.yml ps

# Логи
docker-compose -f docker-compose.clickhouse.yml logs -f hepic-app-server
```

## 🔧 Configuration

Создайте файл `config.env` на основе `config.env.example`:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=hepic_user
DB_PASSWORD=hepic_password
DB_NAME=hepic_db
DB_SSLMODE=disable

# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-here
JWT_EXPIRE_HOURS=24

# Logging
LOG_LEVEL=info
```

## 📚 API Documentation

После запуска сервера документация доступна по адресу:
- **Swagger UI:** http://localhost:8080/api/v1/docs/
- **JSON API:** http://localhost:8080/api/v1/docs/doc.json

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

Authentication
- `POST /api/v1/auth/login` - Вход в систему
- `POST /api/v1/auth/logout` - Выход из системы
- `GET /api/v1/auth/me` - Информация о текущем пользователе

Users
- `GET /api/v1/users` - Список пользователей (с пагинацией)
- `GET /api/v1/users/{id}` - Получить пользователя по ID
- `POST /api/v1/users` - Создать пользователя (только админ)
- `PUT /api/v1/users/{id}` - Обновить пользователя
- `DELETE /api/v1/users/{id}` - Удалить пользователя (только админ)

HEP Records
- `GET /api/v1/hep` - Список HEP записей (с фильтрацией и пагинацией)
- `GET /api/v1/hep/{id}` - Получить HEP запись по ID
- `POST /api/v1/hep` - Создать HEP запись
- `GET /api/v1/hep/stats` - Статистика HEP записей

### Analytics (ClickHouse)
- `GET /api/v1/analytics/stats` - Общая аналитика
- `GET /api/v1/analytics/protocols` - Топ протоколов
- `GET /api/v1/analytics/methods` - Топ методов
- `GET /api/v1/analytics/traffic` - Трафик по часам
- `GET /api/v1/analytics/errors` - Статистика ошибок
- `GET /api/v1/analytics/performance` - Метрики производительности

## 🐳 Docker

### Сборка образа
```bash
make docker
```

### Запуск контейнера
```bash
make docker-run
```

### Docker Compose (рекомендуется)
```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_USER=hepic_user
      - DB_PASSWORD=hepic_password
      - DB_NAME=hepic_db
    depends_on:
      - postgres

  postgres:
    image: postgres:15
    environment:
      - POSTGRES_USER=hepic_user
      - POSTGRES_PASSWORD=hepic_password
      - POSTGRES_DB=hepic_db
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

## 🧪 Testing

```bash
# Запуск тестов
make test

# Запуск в режиме разработки
make dev
```

## 📝 Logging

Сервер использует структурированное JSON логирование с информацией о:
- Времени запроса
- IP адресе клиента
- HTTP методе и URI
- Статусе ответа
- Времени выполнения
- Размере данных

### Slog Middleware Features:
- ✅ **Structured JSON Logging** - Все логи в JSON формате
- ✅ **Request/Response Logging** - Автоматическое логирование HTTP запросов
- ✅ **Error Tracking** - Детальное логирование ошибок с контекстом
- ✅ **Panic Recovery** - Автоматическое восстановление от паник
- ✅ **Performance Metrics** - Время выполнения и размер запросов
- ✅ **Request ID Support** - Автоматическое включение Request ID

Подробная документация: [SLOG_MIDDLEWARE.md](SLOG_MIDDLEWARE.md)

## 🔒 Security

- JWT токены с настраиваемым временем жизни
- Хеширование паролей с bcrypt
- CORS настройки
- Заголовки безопасности
- Валидация входных данных
- Ролевая система доступа

## 📈 Monitoring

- Health check endpoint: `/api/v1/auth/me`
- Структурированные логи
- Метрики производительности
- Graceful shutdown

## 🚀 Performance

- Пул соединений с БД
- Сжатие ответов (Gzip)
- Кэширование заголовков
- Оптимизированные SQL запросы
- Индексы для быстрого поиска

## 📋 TODO

- [ ] Метрики Prometheus
- [ ] Rate limiting
- [ ] Кэширование Redis
- [ ] WebSocket поддержка
- [ ] GraphQL API
- [ ] Миграции БД
- [ ] Unit тесты
- [ ] Integration тесты

## 🤝 Development

### Project Structure
```
├── config/          # Конфигурация
├── database/        # Подключение к БД
├── handlers/        # HTTP handlers (контроллеры)
├── middleware/      # Middleware
├── models/          # Модели данных
├── routes/          # Маршруты API
├── services/        # Бизнес-логика (сервисы)
├── main.go          # Точка входа
├── Dockerfile       # Docker образ
├── Makefile         # Команды сборки
└── config.env       # Конфигурация
```

### Architecture

Проект следует принципам **Clean Architecture** с разделением на слои:

- **Handlers** - HTTP контроллеры, обрабатывают запросы и ответы
- **Services** - Бизнес-логика, работа с данными
- **Models** - Структуры данных и DTO
- **Database** - Слой доступа к данным
- **Routes** - Настройка маршрутов API
- **Middleware** - Промежуточное ПО (CORS, авторизация, логирование)

### Development Commands
```bash
make deps          # Установка зависимостей
make build         # Сборка
make run           # Запуск
make test          # Тесты
make clean         # Очистка
make swagger       # Генерация документации
make build-all     # Полная сборка
```

## 📄 License

MIT License

## 👥 Authors

HEPIC Development Team

## 📞 Support

- Email: support@hepic.local
- Документация: http://localhost:8080/api/v1/docs/
- Issues: GitHub Issues

---

**Версия:** v2.0.0  
**Go Version:** 1.21+  
**Echo Version:** v4.13.4
