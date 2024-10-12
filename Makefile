.PHONY: install
install:
	go install github.com/rubenv/sql-migrate/...@latest


.PHONY: migrate-up
migrate-up:
	@echo "Migrating up..."
	sql-migrate up

.PHONY: migrate-down
migrate-down:
	@echo "Migrating down..."
	sql-migrate down

.PHONY: docker-up
docker-up:
	@echo "Starting docker compose..."
	docker compose up -d

.PHONY: docker-down
docker-down:
	@echo "Stopping docker compose..."
	docker compose down

.PHONY: image-build
image-build:
	@echo "Building a docker image..."
	@bash scripts/image-build.sh $(TAG)

.PHONY: run
run:
	@bash scripts/run.sh
