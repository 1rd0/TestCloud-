upDocker:
	docker compose up
migrate_create:
	migrate create -ext sql -dir  migration/ init
migrate_up:
	migrate -path  migration/ -database "postgresql://user:secret@localhost:5430/postgres_db?sslmode=disable" -verbose up