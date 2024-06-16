ifeq ($(wildcard .env),)
    include .example.env
else
    include .env
endif

DATABASE_URL:="postgresql://${POSTGRES_USERNAME}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable"
BINARY_NAMES="golang-migrate|migrate"
BINARY_NAMES="swag|swag"

ifeq (migrate-create,$(firstword $(MAKECMDGOALS)))
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(RUN_ARGS):;@:)
endif

.checkdeps:
	@echo "Checking dependencies..."
	
	@if ! [ -x "$$(command -v brew)" ]; then \
		echo "brew is not installed."; \
		echo "Installing brew..."; \
		/bin/bash -c "$$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"; \
	else \
		echo "brew is already installed."; \
	fi; \

	@for BINARY in $(BINARY_NAMES); do \
		PACKAGE_NAME=$$(echo $$BINARY | cut -d '|' -f 1); \
		BINARY_NAME=$$(echo $$BINARY | cut -d '|' -f 2); \
		if ! [ -x "$$(command -v $$BINARY_NAME)" ]; then \
			echo "$$BINARY_NAME is not installed."; \
			echo "Installing $$PACKAGE_NAME..."; \
			brew install $$PACKAGE_NAME; \
		else \
			echo "$$BINARY_NAME is already installed."; \
		fi; \
	done

run:
	go run cmd/currencies-api/main.go -c .example.env 

migrate-create: .checkdeps
	migrate create -ext sql -dir migrations/postgres -seq ${RUN_ARGS}

migrate-up: .checkdeps
	migrate -path migrations/postgres -database ${DATABASE_URL} -verbose up

migrate-down: .checkdeps
	migrate -path migrations/postgres -database ${DATABASE_URL} -verbose down 1

migrate-force-previous: .checkdeps
	@version=$$(psql $(DATABASE_URL) -c "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1;" -qt); \
	version=$$((version-1)); \
	migrate -path migrations/postgres -database $(DATABASE_URL) force $${version}

swagger-gen: .checkdeps
	swag init -g cmd/currencies-api/main.go

compose:
	docker-compose up -d