postgres:
	docker run --name postgres13 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:13-alpine

createdb:
	docker exec -it postgres13 createdb --username root --owner root simple_bank

dropdb:
	docker exec -it postgres13 dropdb simple_bank

migrate:
	migrate create -ext sql -dir db/migrations -seq init_schema

migrateup:
	migrate -path db/migrations -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" up

migrateup1:
	migrate -path db/migrations -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" up 1

migratedown:
	migrate -path db/migrations -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" down

migratedown1:
	migrate -path db/migrations -database "postgres://root:secret@localhost:5432/simple_bank?sslmode=disable" down 1

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/asshiddiq1306/simple_bank/db/sql Store

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: postgres createdb dropdb