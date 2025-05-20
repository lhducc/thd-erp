.PHONY: dev-up dev-down

dev-up:
	cp -rf example.env .env
	cp -rf frontend/example.env frontend/.env
	docker compose -f ./deploys/docker-compose.dev.full.yml up -d --build
	sleep 4
	docker compose -f ./deploys/docker-compose.dev.full.yml ps

dev-down:
	docker compose -f ./deploys/docker-compose.dev.full.yml down
