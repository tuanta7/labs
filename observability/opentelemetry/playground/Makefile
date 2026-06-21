.PHONY: env-example up-dev up-signoz-dev up-db-dev

env-example:
	awk -F'=' 'BEGIN {OFS="="} \
    	/^[[:space:]]*#/ {print; next} \
    	/^[[:space:]]*$$/ {print ""; next} \
    	NF>=1 {gsub(/^[[:space:]]+|[[:space:]]+$$/, "", $$1); print $$1"="}' .env > .env.example
	echo ".env.example generated successfully."

up-signoz-dev:
	cd ./docker/signoz && \
	docker compose -f docker-compose.dev.yml up -d --build && \
	cd -

up-db-dev:
	cd ./docker && \
	docker compose -f docker-compose.dev.yml up -d --build && \
	cd -

down-db-dev:
	cd ./docker && \
	docker compose -f docker-compose.dev.yml down && \
	cd -

up-dev: up-signoz-dev up-db-dev

