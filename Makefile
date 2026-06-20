include .env
export

MIGRATIONS_PATH=db/migrations

.PHONY: run build test test-coverage migrate-up migrate-down migrate-force docker-up docker-down docker-build seed clean

run:
	go run ./cmd/server

##unit tests
test:
	go test ./... -v

migrate-up:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URI)" up

migrate-down:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URI)" down 1

migrate-force:
	migrate -path $(MIGRATIONS_PATH) -database "$(DB_URI)" force $(VERSION)

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

## rebuild docker images without starting
docker-build:
	docker compose build

## Seed the database manually
seed:
	SEED_DB=true go run ./cmd/server

