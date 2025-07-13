#!make

# Docker image to run shell and go utility functions in
WORKER_IMAGE = golang:1.24-alpine
# Docker image to generate OAS3 specs
OAS3_GENERATOR_DOCKER_IMAGE = openapitools/openapi-generator-cli:latest-release

MIGRATION_NAME ?= $(shell bash -c 'read -p "Migration name: " migration_name; echo $$migration_name')

.PHONY: fmt lint coverage docs-% db-%


fmt:
	@go mod tidy -v
	@go fmt ./...

lint:
	@echo "[GOCI] Linting project files"
	@golangci-lint run --tests=false

coverage:
	@echo "[TEST] Generating coverage report"
	@xgo test -v -cover -coverpkg ./... -coverprofile=cover.out ./...
	@go test -coverprofile=c.out
	@go tool cover -html=c.out -o coverage.html

docs-swagger:
	@echo "[SWAG] Generating OpenAPI 2.0 schema"
	@swag init -g main.go

docs-convert:
	@echo "[OAS3] Converting Swagger 2-to-3 (yaml)"
	@docker run --rm -v $(PWD)/docs:/work $(OAS3_GENERATOR_DOCKER_IMAGE) \
	  generate -i /work/swagger.yaml -o /work/v3 -g openapi-yaml --minimal-update
	@docker run --rm -v $(PWD)/docs/v3:/work $(WORKER_IMAGE) \
	  sh -c "rm -rf /work/.openapi-generator"
	@echo "[OAS3] Copying openapi-generator-ignore (json)"
	@docker run --rm -v $(PWD)/docs/v3:/work $(WORKER_IMAGE) \
	  sh -c "cp -f /work/.openapi-generator-ignore /work/openapi"
	@echo "[OAS3] Converting Swagger 2-to-3 (json)"
	@docker run --rm -v $(PWD)/docs:/work $(OAS3_GENERATOR_DOCKER_IMAGE) \
	  generate -s -i /work/swagger.json -o /work/v3/openapi -g openapi --minimal-update
	@echo "[OAS3] Cleaning up generated files"
	@docker run --rm -v $(PWD)/docs/v3:/work $(WORKER_IMAGE) \
	  sh -c "mv -f /work/openapi/openapi.json /work ; mv -f /work/openapi/openapi.yaml /work ; rm -rf /work/openapi"

docs-gen:
	@make docs-swagger
	@make docs-convert

db-clean:
	@echo "[ATLS] Clearing DB"
	@atlas schema clean --env gorm --url $$DB_URL

db-push:
	@echo "[ATLS] Pushing DB schema to Atlas cloud"
	@atlas migrate push license-server --env gorm --dev-url $$DB_URL

db-apply:
	@echo "[ATLS] Applying schema to DB"
	@atlas schema apply --env gorm --url $$DB_URL --dev-url $$DB_URL --auto-approve

db-schema:
	@echo "[ATLS] Opening DB schema viewer"
	@atlas schema inspect -w --env gorm --url $$DB_URL

db-refresh:
	@echo "[ATLS] Recreating DB"
	@make db-clean
	@make db-apply

db-lint:
	@echo "[ATLS] Linting DB"
	@atlas migrate lint --dev-url $$DB_URL -w

db-diff:
	@echo "[ATLS] Showing schema diff"
	@atlas schema diff --env gorm --from $$DB_URL --to "./migrations?format=atlas&version="

db-migration:
	@echo "[ATLS] Generating migration and applying to DB"
	@make db-clean
	@atlas migrate diff $(MIGRATION_NAME) --dev-url $$DB_URL --env gorm
	@make db-push
	@make db-apply
