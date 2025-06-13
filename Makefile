.PHONY: dev-up dev-down uat-up uat-down

dev-up:
	cp -rf example.env .env
	cp -rf frontend/example.env frontend/.env
	docker compose -f ./deploys/docker-compose.dev.full.yml up -d --build
	sleep 4
	docker compose -f ./deploys/docker-compose.dev.full.yml ps

	# zip node_modules for no-internet in uat and prod env
	docker exec erp-frontend tar -czf /tmp/node_modules.tar.gz -C /app node_modules
	docker cp erp-frontend:/tmp/node_modules.tar.gz ./frontend/

dev-down:
	docker exec erp-frontend rm -f /tmp/node_modules.tar.gz
	docker compose -f ./deploys/docker-compose.dev.full.yml down

uat-up:
	cp -rf example.env .env
	cp -rf frontend/example.env frontend/.env
	tar -xzvf ./frontend/node_modules.tar.gz -C ./frontend/
	docker compose -f ./deploys/docker-compose.uat.yml up -d --build
	sleep 4
	docker compose -f ./deploys/docker-compose.uat.yml ps

uat-down:
	rm -rf ./frontend/node_modules
	docker compose -f ./deploys/docker-compose.uat.phat.yml down
