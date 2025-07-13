# Docker image to run shell and go utility functions in
WORKER_IMAGE = golang:1.24-alpine
# Docker image to generate OAS3 specs
OAS3_GENERATOR_DOCKER_IMAGE = openapitools/openapi-generator-cli:latest-release

.PHONY: swag up ci-swaggen build debug fmt test

swag:
	cd backend && watch -n 10 swag init

up:
	@docker compose up -d --build
	@terraform apply -auto-approve
	@docker exec -itd backend sh "go run main.go seed"

seed:
	@awslocal ssm put-parameter --name "api_key" --value $$TF_VAR_api_key \
		--type String

