#psql -h localhost -p 5432  -U postgres -W -c "\c simple_bank"
#docker run --name postgres13 -p 5432:5432 -e PG_USER=root -e PG_PASSWORD=mypassword -d postgres:latest

postgres:
	psql "dbname=simple_bank host=localhost user=${PSQLUSER} password=${PSQLPASS} port=5432 sslmode=disable"


createdb:
	docker exec -it postgres createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres dropdb simple_bank

sqlc:
	sqlc generate


migrateDown:
	migrate -path db/migration -database "postgresql://${PSQLUSER}:${PSQLPASS}@localhost:5432/simple_bank?sslmode=disable" -verbose down


migrateUp:
	migrate -path db/migration -database "postgresql://${PSQLUSER}:${PSQLPASS}@localhost:5432/simple_bank?sslmode=disable" -verbose up

test:
	go test -v -cover ./...

.PHONY: postgres createdb dropdb migrateDown migrateUp sqlc