.PHONY: generate run

generate:
	@mkdir -p internal/types internal/api
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config oapi-codegen-types.yaml api/types.yaml
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen -config oapi-codegen-server.yaml api/openapi.yaml
	go mod tidy

run:
	go run cmd/server/main.go
