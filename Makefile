postgres:
	docker run --name postgresdb -p 5000:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:alpine3.16

createdb:
	docker exec -it postgresdb createdb --username=root --owner=root postgres

dropdb:
	docker exec -it postgresdb dropdb persons

migrateup:
	migrate -path db/migration -database "postgresql://jhivan:25May2001@grama-check-db.postgres.database.azure.com/postgres?sslmode=require" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://jhivan:25May2001@grama-check-db.postgres.database.azure.com/postgres?sslmode=require" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...
	
.PHONY:  postgres createdb dropdb migrateup migratedown sqlc