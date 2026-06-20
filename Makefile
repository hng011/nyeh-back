dev:
	@echo "=== Updating Swagger Docs"
	@swag init -g cmd/api/main.go

	@echo "=== Executing swag fmt"
	@swag fmt

	@echo "\n\n=== Running the API :/"
	go run cmd/api/main.go
