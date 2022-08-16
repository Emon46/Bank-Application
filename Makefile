postgres:
	docker run --name bank-postgres -p 5432:5432 -e POSTGRES_PASSWORD=secret1234 -d postgres:14.5
createdb:
	docker exec -it bank-postgres createdb --username=postgres --owner=postgres simple_bank
dropdb:
	docker exec -it bank-postgres dropdb --username=postgres simple_bank
migrateup:
	migrate --path db/migration --database "postgresql://postgres:secret1234@localhost:5432/simple_bank?sslmode=disable" -verbose up
migratedown:
	migrate --path db/migration --database "postgresql://postgres:secret1234@localhost:5432/simple_bank?sslmode=disable" -verbose down
.PHONY: postgres createdb migrateup migratedown
