postgres:
	docker run --name simple-bank -p 5432:5432 -d -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret postgres

createdb:
	docker exec -it simple-bank createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it simple-bank dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1


sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/mhg14/simplebank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock migrateup1 migratedown1