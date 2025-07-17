#!make

# Docker image to run shell and go utility functions in
WORKER_IMAGE = golang:1.24-alpine
# Docker image to generate OAS3 specs
OAS3_GENERATOR_DOCKER_IMAGE = openapitools/openapi-generator-cli:latest-release

.PHONY: docs-%

docs-swagger:
	@echo "[SWAG] Generating OpenAPI 2.0 schema"
	@swag init -g main.go

docs-convert:
	@echo "[OAS3] Converting Swagger 2-to-3 (yaml)"
	@docker run --rm -v $(PWD)/services/jobsvc/docs:/work $(OAS3_GENERATOR_DOCKER_IMAGE) \
	  generate -i /work/swagger.yaml -o /work/v3 -g openapi-yaml --minimal-update
	@docker run --rm -v $(PWD)/services/jobsvc/docs/v3:/work $(WORKER_IMAGE) \
	  sh -c "rm -rf /work/.openapi-generator"
	@echo "[OAS3] Copying openapi-generator-ignore (json)"
	@docker run --rm -v $(PWD)/services/jobsvc/docs/v3:/work $(WORKER_IMAGE) \
	  sh -c "cp -f /work/.openapi-generator-ignore /work/openapi"
	@echo "[OAS3] Converting Swagger 2-to-3 (json)"
	@docker run --rm -v $(PWD)/services/jobsvc/docs:/work $(OAS3_GENERATOR_DOCKER_IMAGE) \
	  generate -s -i /work/swagger.json -o /work/v3/openapi -g openapi --minimal-update
	@echo "[OAS3] Cleaning up generated files"
	@docker run --rm -v $(PWD)/services/jobsvc/docs/v3:/work $(WORKER_IMAGE) \
	  sh -c "mv -f /work/openapi/openapi.json /work ; mv -f /work/openapi/openapi.yaml /work ; rm -rf /work/openapi"

docs-gen:
	@make docs-swagger
	@make docs-convert

docs-watch:
	@watch -n 10 swag init -g main.go
