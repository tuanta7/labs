.PHONY: setup deploy verify check-ram

setup:
	# https://signoz.io/docs/install/docker/
	curl -fsSL https://signoz.io/foundry.sh | bash

check-ram:
	docker info --format 'OS: {{.OperatingSystem}} | CPU: {{.NCPU}} | Memory: {{.MemTotal}} bytes'
	@mem=$$(docker info --format '{{.MemTotal}}'); \
	min=$$((6*1024*1024*1024)); \
	if [ "$$mem" -lt "$$min" ]; then \
		echo "Docker has $$((mem/1024/1024/1024))GB RAM allocated, need at least 6GB. Increase it in Docker Desktop > Settings > Resources."; \
		exit 1; \
	fi

deploy: check-ram
	# Open ports: 8080 (SigNoz UI), 4317 and 4318 (OTLP ingestion)
	# and 8000 if enable the MCP server.
	foundryctl cast -f casting.yaml

verify:
	docker ps