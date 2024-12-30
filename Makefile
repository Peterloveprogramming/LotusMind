# create docker contianer
postgres:
	docker run --name postgres17 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17-alpine
# create database inside the container
createdb:
	docker exec -it postgres17 createdb --username=root --owner=root meditation
# drop database inside the container
dropdb:
	docker exec -it postgres17 dropdb meditation
# create the tables in the database
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/meditation?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/meditation?sslmode=disable" -verbose down
sqlc:
	sqlc generate
testV:
	go test -v -cover ./...
test:
	go test ./...
.PHONY: postgres createdb dropdb migrateup migratedown sqlc testV


