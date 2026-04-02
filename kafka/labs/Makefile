.PHONY: run-pub run-sub

destroy:
	docker compose -f ./docker/docker-compose.dev.yml down

setup: destroy
	docker compose -f ./docker/docker-compose.dev.yml up -d

run-pub:
	go run ./cmd/publisher

run-sub:
	go run ./cmd/subscriber