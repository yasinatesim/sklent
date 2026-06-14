.PHONY: up down build api-build api-test api-vet web-dev web-build verify fmt

up:
	docker compose up --build

down:
	docker compose down

seed:
	docker compose run --rm api ./bin/seed

api-build:
	cd api && go build ./...

api-vet:
	cd api && go vet ./...

api-test:
	cd api && go test ./...

api: api-build api-vet api-test

web-dev:
	cd web && npm run dev

web-build:
	cd web && npm run build

verify:
	node e2e/verify.mjs

fmt:
	cd api && gofmt -w .
