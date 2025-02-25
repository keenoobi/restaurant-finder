.PHONY: all build elastic-up elastic-down load-data run start stop clean test help

# Переменные
BIN_DIR=bin
CMD_DIR=cmd
WEB_APP_BIN=$(BIN_DIR)/web-app
ELASTIC_APP_BIN=$(BIN_DIR)/elastic-app
CONFIG_PATH=internal/config/config.yaml
DOCKER_COMPOSE_FILE=docker-compose.yml

# Цели по умолчанию
all: build

# Сборка проекта
build:
	@echo "Сборка проекта..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(WEB_APP_BIN) ./$(CMD_DIR)/web
	@go build -o $(ELASTIC_APP_BIN) ./$(CMD_DIR)/elastic-app
	@echo "Сборка завершена."

# Запуск Elasticsearch через Docker Compose
elastic-up:
	@echo "Запуск Elasticsearch..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up -d
	@echo "Elasticsearch запущен."

# Остановка Elasticsearch
elastic-down:
	@echo "Остановка Elasticsearch..."
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down
	@echo "Elasticsearch остановлен."

# Загрузка данных в Elasticsearch
load-data: build elastic-up
	@echo "Загрузка данных в Elasticsearch..."
	@$(ELASTIC_APP_BIN) -config $(CONFIG_PATH)
	@echo "Данные загружены."

# Запуск веб-приложения
run: build
	@echo "Запуск веб-приложения..."
	@$(WEB_APP_BIN) -config $(CONFIG_PATH)

# Запуск всех сервисов (Elasticsearch + веб-приложение)
start: elastic-up run

# Остановка всех сервисов
stop: elastic-down
	@echo "Все сервисы остановлены."

# Очистка проекта
clean:
	@echo "Очистка проекта..."
	@rm -rf $(BIN_DIR)
	@echo "Очистка завершена."

# Запуск тестов
test:
	@echo "Запуск тестов..."
	@go test ./...
	@echo "Тесты завершены."

# Помощь (список доступных команд)
help:
	@echo "Доступные команды:"
	@echo "  make build         - Сборка проекта"
	@echo "  make elastic-up    - Запуск Elasticsearch через Docker Compose"
	@echo "  make elastic-down  - Остановка Elasticsearch"
	@echo "  make load-data     - Загрузка данных в Elasticsearch"
	@echo "  make run           - Запуск веб-приложения"
	@echo "  make start         - Запуск всех сервисов (Elasticsearch + веб-приложение)"
	@echo "  make stop          - Остановка всех сервисов"
	@echo "  make clean         - Очистка проекта (удаление скомпилированных файлов)"
	@echo "  make test          - Запуск тестов"
	@echo "  make help          - Показать это сообщение"