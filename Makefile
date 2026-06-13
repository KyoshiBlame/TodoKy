include .env
export

PROJECT_ROOT := $(CURDIR)
export PROJECT_ROOT

.PHONY: env-up env-down env-clean-up

env-up:
	docker compose up -d todoky-postgres

env-down:
	docker compose down

env-clean-up:
	@printf "Очистить все файлы окружения? Возможна потеря данных. [y/N]: "; \
	read ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down && \
		rm -rf out/pgdata && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует параметр seq. Пример команды: make migrate-create seq=example"; \
		exit 1; \
	fi
	@docker compose run --rm todoky-postgres-migrate \
		create -ext sql -dir /migrations -seq "$(seq)"


test-target:
	@echo "value: $(var)" 