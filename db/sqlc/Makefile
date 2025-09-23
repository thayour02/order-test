postgres:
	docker run --name postgres21 -p 5434:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres21 createdb --username=root --owner=root orders

dropdb:
	docker exec -it postgres21 dropdb  orders
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/orders?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/order?sslmode=disable" -verbose down
sqlc:
	sqlc generate
.PHONY: postgres createdb dropdb migrateup migratedown sqlc