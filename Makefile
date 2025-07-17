# Docker image to run shell and go utility functions in
WORKER_IMAGE = golang:1.24-alpine
# Docker image to generate OAS3 specs
OAS3_GENERATOR_DOCKER_IMAGE = openapitools/openapi-generator-cli:latest-release

.PHONY: swag up ci-swaggen build debug fmt test migrate-%

docs-jobs:
	@watch -n 10 swag init -g ./services/jobsvc/main.go -o ./services/jobsvc/docs

up:
	@docker compose up -d --build
	@terraform apply -auto-approve
	@docker exec -itd backend sh "go run main.go seed"

seed:
	@awslocal ssm put-parameter --name "api_key" --value $$TF_VAR_api_key \
		--type String

migrate-up:
	@migrate -path cassandra-migrations/ -database $$CASS_URL -verbose up

migrate-down:
	@migrate -path cassandra-migrations/ -database $$CASS_URL -verbose down
