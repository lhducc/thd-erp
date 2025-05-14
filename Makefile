.PHONY: dev-up dev-down

dev-up:
	docker compose -f ./deploys/docker-compose.dev.full.yml up -d
	sleep 4
	docker compose -f ./deploys/docker-compose.dev.full.yml ps

dev-down:
	docker compose -f ./deploys/docker-compose.dev.full.yml down
