.PHONY:

run_products_service:
	cd services/products/ && 	go run ./cmd/main.go

# ==============================================================================
# Docker Compose
docker-compose_infra_up:
	@echo Starting infrastructure docker-compose
	docker-compose -f deployments/docker-compose/docker-compose.infrastructure.yaml up --build

docker-compose_infra_down:
	@echo Stoping infrastructure docker-compose
	docker-compose -f deployments/docker-compose/docker-compose.infrastructure.yaml down

# ==============================================================================
# Docker

FILES := $(shell docker ps -aq)

docker_path:
	@echo $(FILES)

docker_down:
	docker stop $(FILES)
	docker rm $(FILES)

docker_clean:
	docker system prune -f

docker_logs:
	docker logs -f $(FILES)


# ==============================================================================
# Linters https://golangci-lint.run/usage/install/

run-linter:
	@echo Starting linters
	golangci-lint run ./...

# ==============================================================================
# Proto Identity Service

## go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

## grpc-server
proto_identities_get_user_by_id_service:
	@echo Generating identity_service proto
	protoc --go_out=./internal/services/identity-service/identity/grpc_server/protos --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=./internal/services/identity-service/identity/grpc_server/protos ./internal/services/identity-service/identity/grpc_server/protos/*.proto


## grpc-client
proto_identities_get_user_by_id_service:
	@echo Generating identity_service proto
	protoc --go_out=./internal/services/product-service/product/grpc_client/protos --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=./internal/services/product-service/product/grpc_client/protos ./internal/services/product-service/product/grpc_client/protos/*.proto

# ==============================================================================
# Swagger products Service  #https://github.com/swaggo/swag/issues/817

install_swag_products:
	cd internal/services/product-service/ && 	go get -u github.com/swaggo/swag/cmd/swag@latest

install_swag_identities:
	cd internal/services/identity-service/ && 	go get -u github.com/swaggo/swag/cmd/swag@latest

swagger_products:
	@echo Starting swagger generating
	swag init --parseDependency --parseInternal --parseDepth 1 -g ./cmd/main.go -d ./internal/services/product-service/ -o ./internal/services/product-service/docs

swagger_identities:
	@echo Starting swagger generating
	swag init --parseDependency --parseInternal --parseDepth 1 -g ./cmd/main.go -d ./internal/services/identity-service/ -o ./internal/services/identity-service/docs