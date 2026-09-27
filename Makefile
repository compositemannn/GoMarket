include .env
export

env-up: # поднимаем окружение для приложения
	@docker compose up -d postgres

env-down: # останавливаем
	@docker compose down postgres

env-cleanup: # очистка окружения
	@read -p "Очистить все volume файлы окружения? Опасность утечки данных. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down postgres && \
		rm -rf out/pg_data \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

env-port-forwarder-up:
	@docker compose up -d port-forwarder

env-port-forwarder-down:
	@docker compose down port-forwarder

migrate-create: # up - для долгоживущих, run - для разовых сервисов, --rm - удалить контейнер migrate
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует необходимый параметр seq. Пример: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует необходимый параметр action. Пример: make migrate-action action=up 1"; \
		exit 1; \
	fi; \
	docker compose run --rm postgres-migrate \
		-path /migrations \
		-database postgres://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}?sslmode=disable \
		"${action}"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down
