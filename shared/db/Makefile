# Default DB connection (kann beim Aufruf überschrieben werden)
# postgres://USER:PASSWORD@HOST:PORT/DBNAME?sslmode=disable
DB_URL ?= postgres://postgres:example@localhost:5432/postgres?sslmode=disable

.PHONY: migrate-up migrate-down migrate-create sqlc generate

migrate-up:       ## Führe alle Migrations aus
	migrate -database "$(DB_URL)" -path ./migrations  up

migrate-down:     ## Rolle letzte Migration zurück
	migrate -path migrations -database "$(DB_URL)" down 1

migrate-create:   ## Neue Migration anlegen (Name bei Prompt) migrate create -ext sql -dir db/migrations -seq create_users_table
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir ./migrations -seq $$name

# --- SQLC ---
sqlc:             ## Generiere Go Code aus SQL
	sqlc generate -f sqlc/sqlc.yaml

# --- Full generation workflow ---
generate: sqlc migrate-up  ## SQLC + Migration together
	@echo "✅ Schema & Code aktuell"
