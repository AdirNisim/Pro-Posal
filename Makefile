-include _.env

SHELL=/bin/bash

.PHONY: run-tests
run-tests:
	go test -v ./...

.PHONY: run-db
run-db:
	docker-compose up -d --build

.PHONY: dao-gen
dao-gen: run-db
	sqlboiler psql

.PHONY: db-connect
db-connect:
	PGPASSWORD=Aa123456 psql -h localhost -p 5432 -U admin pro-posal

.PHONY: db-down
db-down:
	goose -dir=./migrations postgres "host=localhost user=admin dbname=pro-posal sslmode=disable password=Aa123456" down

.PHONY: db-reset
db-reset:
	goose -dir=./migrations postgres "host=localhost user=admin dbname=pro-posal sslmode=disable password=Aa123456" down-to 0


.PHONY:	db-up
db-up:
	goose -dir=./migrations postgres "host=localhost user=admin dbname=pro-posal sslmode=disable password=Aa123456" up