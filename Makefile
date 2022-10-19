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
# PPROF

pprof_heap:
	go tool pprof -http :8006 http://localhost:6060/debug/pprof/heap?seconds=10

pprof_cpu:
	go tool pprof -http :8006 http://localhost:6060/debug/pprof/profile?seconds=10

pprof_allocs:
	go tool pprof -http :8006 http://localhost:6060/debug/pprof/allocs?seconds=10

# ==============================================================================
# Proto products Service

proto_products_product_kafka_messages:
	@echo Generating products kafka messages proto
	protoc --go_out=./services/products/internal/products/contracts/grpc/kafka_messages --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=./services/products/internal/products/contracts/grpc/kafka_messages api_docs/products/protobuf/products/kafka_messages/product_kafka_messages.proto

proto_products_product_service:
	@echo Generating product_service client proto
	protoc --go_out=./services/products/internal/products/contracts/grpc/service_clients --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=./services/products/internal/products/contracts/grpc/service_clients api_docs/products/protobuf/products/service_clients/product_service_client.proto

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