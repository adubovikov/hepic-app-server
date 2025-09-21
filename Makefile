# HEPIC App Server v2 Makefile

.PHONY: build run test clean deps swagger docker

# Переменные
BINARY_NAME=hepic-app-server-v2
DOCKER_IMAGE=hepic-app-server-v2
DOCKER_TAG=latest

# Сборка приложения
build:
	@echo "Сборка HEPIC App Server v2..."
	go build -o $(BINARY_NAME) .

# Запуск приложения
run: build
	@echo "Запуск HEPIC App Server v2..."
	./$(BINARY_NAME)

# Установка зависимостей
deps:
	@echo "Установка зависимостей..."
	go mod tidy
	go mod download

# Генерация Swagger документации
swagger:
	@echo "Генерация Swagger документации..."
	swag init -g main.go -o ./docs

# Запуск тестов
test:
	@echo "Запуск тестов..."
	go test -v ./...

# Очистка
clean:
	@echo "Очистка..."
	rm -f $(BINARY_NAME)
	rm -rf docs/

# Docker сборка
docker:
	@echo "Сборка Docker образа..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Docker запуск
docker-run: docker
	@echo "Запуск Docker контейнера..."
	docker run -p 8080:8080 --env-file config.env $(DOCKER_IMAGE):$(DOCKER_TAG)

# Установка swag
install-swagger:
	@echo "Установка swag..."
	go install github.com/swaggo/swag/cmd/swag@latest

# Полная сборка с документацией
build-all: deps swagger build
	@echo "Полная сборка завершена"

# Разработка
dev: deps
	@echo "Запуск в режиме разработки..."
	go run .

# Помощь
help:
	@echo "Доступные команды:"
	@echo "  build        - Сборка приложения"
	@echo "  run          - Запуск приложения"
	@echo "  deps         - Установка зависимостей"
	@echo "  swagger      - Генерация Swagger документации"
	@echo "  test         - Запуск тестов"
	@echo "  clean        - Очистка"
	@echo "  docker       - Сборка Docker образа"
	@echo "  docker-run   - Запуск Docker контейнера"
	@echo "  dev          - Запуск в режиме разработки"
	@echo "  build-all    - Полная сборка с документацией"
