include .env
export $(shell sed 's/=.*//' .env)

.PHONY: psql migration test

psql:
	@docker exec -it bt-bot-db psql -U ${POSTGRES_USER} -d ${POSTGRES_NAME}

migration:
	@docker exec -it bt-bot-db psql -U ${POSTGRES_USER} -d ${POSTGRES_NAME} -f $(m)

test:
	docker exec -it bt-bot go test ./...
