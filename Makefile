.PHONY: install
install:
	go install github.com/rubenv/sql-migrate/...@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@v1.8.12
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	# go install github.com/golang/mock/mockgen@latest
	go install go.uber.org/mock/mockgen@latest
	


.PHONY: migrate-up
migrate-up:
	@echo "Migrating up..."
	@bash scripts/migrate_up.sh

.PHONY: migrate-down
migrate-down:
	@echo "Migrating down..."
	@bash scripts/migrate_down.sh

.PHONY: migrate-clean
migrate-clean:
	@echo "Migrating clean..."
	@bash scripts/migrate_clean.sh


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
	@bash scripts/image_build.sh $(TAG)

.PHONY: run
run:
	@bash scripts/run.sh

.PHONY: run-debug
run-debug:
	@bash scripts/run_debug.sh



.PHONY: docs-gen
docs-gen:
	swag init -g cmd/main/main.go --parseDependency --parseInternal -o ./generated/docs

.PHONY: lint
lint:
	golangci-lint run --fix


.PHONY: build
build:
	@bash scripts/build.sh $(version)
