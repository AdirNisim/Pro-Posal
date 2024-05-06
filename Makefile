-include _.env

SHELL=/bin/bash

.PHONY: run-db
run-db:
	docker-compose up -d --build

.PHONY: dao-gen
dao-gen: run-db
	sqlboiler psql

.PHONY: db-connect
db-connect:
	PGPASSWORD=Aa123456 psql -h localhost -p 5432 -U admin pro-posal

.PHONY: db-migrate-up
db-migrate-up:
	goose -dir=./migrations postgres "user=admin dbname=pro-posal sslmode=disable password=Aa123456" up

.PHONY: db-migrate-down
db-migrate-down:
	goose -dir=./migrations postgres "user=admin dbname=pro-posal sslmode=disable password=Aa123456" down

.PHONY: db-migrate-reset
db-migrate-reset:
	goose -dir=./migrations postgres "user=admin dbname=pro-posal sslmode=disable password=Aa123456" down-to 0
