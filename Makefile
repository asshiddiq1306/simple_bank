postgres:
	docker run --name postgres13 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:13-alpine

createdb:
	docker exec -it postgres13 createdb --username root --owner root simple_bank

dropdb:
	docker exec -it postgres13 dropdb simple_bank

migrate:
	migrate create -ext sql -dir db/migrations -seq init_schema

migrateup:
	migrate -path db/migrations -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" up

migratedown:
	migrate -path db/migrations -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" down

.PHONY: postgres createdb dropdb