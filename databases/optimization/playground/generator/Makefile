.PHONY: run migrate drop

run:
	docker compose -f ../docker-compose.dev.yml up -d
	go run .

migrate:
	go run .

drop:
	docker compose -f ../docker-compose.dev.yml down