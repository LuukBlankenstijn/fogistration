.PHONY: start stop build swagger migration-% migrate gen-db gen-domjudge-client

start:
	@docker compose up -d

# Stop and clean up
stop:
	@docker compose down

# Rebuild and start container
build:
	@docker build .

# Generate API documentation
swagger:
	@docker compose exec http swag init -g cmd/http-server/main.go

# Create en migration with name everything behind the dash
migration-%:
	@docker compose exec http mkdir -p ./internal/shared/database/migrations
	@docker compose exec http tern new $* --migrations ./internal/shared/database/migrations/ 

# Migrate up to the current version
migrate:
	# @docker compose exec http bash -c "source .env-http && tern migrate --migrations ./internal/shared/database/migrations/"
	@docker compose exec http tern migrate --migrations ./internal/shared/database/migrations/

migrate-down:
	@docker compose exec http tern migrate -d -1 --migrations ./internal/shared/database/migrations/

# generate the db code
gen-db:
	@go run cmd/schema-builder/main.go
	@docker compose exec http sqlc generate

gen-domjudge-client:
	@cd internal/domjudge/client && go generate
