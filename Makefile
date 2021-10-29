DB_CONTAINER?=rw_db
env?=local
remove-infras:
	docker-compose stop
	docker-compose rm -f

init: remove-infras
	docker-compose up -d lfw_db
	@echo "Waiting for database connection..."
	@while ! docker exec $(DB_CONTAINER) pg_isready -h localhost -p 5432 > /dev/null; do \
		sleep 1; \
	done

.PHONY: run
run:
	go run cmd/*.go

.PHONY: indexer
indexer:
	go run cmd/indexer/*.go

.PHONY: test
test:
	@PROJECT_PATH=$(shell pwd) go test -cover ./...

.PHONY: migrate
migrate:
	sql-migrate up -env=$(env)

.PHONY: gen-swagger
gen-swagger:
	@swag init -g ./cmd/server.go

.PHONY: gen-mock
gen-mock:
	@mockgen -source=./repo/repo.go -destination=./repo/mocks/repo.go
