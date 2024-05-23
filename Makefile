dev:
	cd ./cmd; air
dev-db:
	docker compose -f compose-dev.yaml up -d
prod:
	docker compose up --build
migrate:
	cd ./migrations; goose postgres postgres://postgres:pass@localhost:5432/data up
build-openapi:
	cd ./api; redocly build-docs ./openapi.yaml --output=index.html