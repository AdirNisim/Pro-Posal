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
