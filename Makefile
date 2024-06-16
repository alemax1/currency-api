ifeq ($(wildcard .env),)
    include .example.env
else
    include .env
endif

DATABASE_URL:="postgresql://${POSTGRES_USERNAME}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable"

ifeq (migrate-create,$(firstword $(MAKECMDGOALS)))
  RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
  $(eval $(RUN_ARGS):;@:)
endif

run:
	go run cmd/currencies-api/main.go -c .example.env 

migrate-create:
	migrate create -ext sql -dir migrations/postgres -seq ${RUN_ARGS}

migrate-up:
	migrate -path migrations/postgres -database ${DATABASE_URL} -verbose up

migrate-down:
	migrate -path migrations/postgres -database ${DATABASE_URL} -verbose down 1

migrate-force-previous:
	@version=$$(psql $(DATABASE_URL) -c "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1;" -qt); \
	version=$$((version-1)); \
	migrate -path migrations/postgres -database $(DATABASE_URL) force $${version}

swagger-gen:
	swag init -g ./cmd/currencies-api/main.go

compose:
	docker-compose up -d