.PHONY: dev-fe-up dev-fe-down dev-be-up dev-be-down

dev-fe-up:
	./scripts/dev-fe.sh up

dev-fe-down:
	./scripts/dev-fe.sh down

dev-be-up:
	./scripts/dev-be.sh up

dev-be-down:
	./scripts/dev-be.sh down
