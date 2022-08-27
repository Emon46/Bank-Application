postgres:
	docker run --name bank-postgres -p 5432:5432 -e POSTGRES_PASSWORD=secret1234 -d postgres:14.5

createdb:
	docker exec -it bank-postgres createdb --username=postgres --owner=postgres simple_bank

dropdb:
	docker exec -it bank-postgres dropdb --username=postgres simple_bank

migrateup:
	migrate --path db/migration --database "postgresql://postgres:secret1234@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate --path db/migration --database "postgresql://postgres:secret1234@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate --path db/migration --database "postgresql://postgres:secret1234@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate --path db/migration --database "postgresql://postgres:secret1234@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

fmt:
	go fmt ./...

test:
	go test -v -cover ./...
test-coverage:
	go test ./... -coverprofile=test-coverage.out -covermode count
	go tool cover -func=test-coverage.out | grep total | grep -Eo '[0-9]+\.[0-9]+'
	rm -rf test-coverage.out

sqlc:
	sqlc generate

server:
	go run main.go

mockdb:
	mockgen -destination db/mock/store.go -package mockdb github.com/emon46/bank-application/db/sqlc Store

.PHONY: fmt postgres createdb migrateup migratedown sqlc test server mockdb
