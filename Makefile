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