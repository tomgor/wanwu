include .env
include .env.image.${WANWU_ARCH}

LDFLAGS := -X main.buildTime=$(shell date +%Y-%m-%d,%H:%M:%S) \
			-X main.buildVersion=${WANWU_VERSION} \
			-X main.gitCommitID=$(shell git --git-dir=./.git rev-parse HEAD) \
			-X main.gitBranch=$(shell git --git-dir=./.git for-each-ref --format='%(refname:short)->%(upstream:short)' $(shell git --git-dir=./.git symbolic-ref -q HEAD)) \
			-X main.builder=$(shell git config user.name)

build-bff-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -ldflags "$(LDFLAGS)" -o ./bin/amd64/ ./cmd/bff-service

build-bff-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -mod vendor -ldflags "$(LDFLAGS)" -o ./bin/arm64/ ./cmd/bff-service

build-iam-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -ldflags "$(LDFLAGS)" -o ./bin/amd64/ ./cmd/iam-service

build-iam-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -mod vendor -ldflags "$(LDFLAGS)" -o ./bin/arm64/ ./cmd/iam-service

build-model-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -ldflags "$(LDFLAGS)" -o ./bin/amd64/ ./cmd/model-service

build-model-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -mod vendor -ldflags "$(LDFLAGS)" -o ./bin/arm64/ ./cmd/model-service

build-mcp-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -ldflags "$(LDFLAGS)" -o ./bin/amd64/ ./cmd/mcp-service

build-mcp-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -mod vendor -ldflags "$(LDFLAGS)" -o ./bin/arm64/ ./cmd/mcp-service

build-knowledge-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -ldflags "$(LDFLAGS)" -o ./bin/amd64/ ./cmd/knowledge-service

build-knowledge-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -mod vendor -ldflags "$(LDFLAGS)" -o ./bin/arm64/ ./cmd/knowledge-service

build-rag-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -ldflags "$(LDFLAGS)" -o ./bin/amd64/ ./cmd/rag-service

build-rag-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -mod vendor -ldflags "$(LDFLAGS)" -o ./bin/arm64/ ./cmd/rag-service

build-app-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -ldflags "$(LDFLAGS)" -o ./bin/amd64/ ./cmd/app-service

build-app-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -mod vendor -ldflags "$(LDFLAGS)" -o ./bin/arm64/ ./cmd/app-service

build-operate-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -ldflags "$(LDFLAGS)" -o ./bin/amd64/ ./cmd/operate-service

build-operate-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -mod vendor -ldflags "$(LDFLAGS)" -o ./bin/arm64/ ./cmd/operate-service

build-assistant-amd64:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod vendor -ldflags "$(LDFLAGS)" -o ./bin/amd64/ ./cmd/assistant-service

build-assistant-arm64:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -mod vendor -ldflags "$(LDFLAGS)" -o ./bin/arm64/ ./cmd/assistant-service

create-docker-net:
	docker network create ${WANWU_DOCKER_NETWORK}

check:
	go vet ./...
	go fmt ./...
	docker run --rm -t -v $(PWD):/app -w /app golangci/golangci-lint:v1.64.8 bash -c 'golangci-lint run -v --timeout 3m'

doc:
	docker run --name golang-swag --privileged=true --rm -v $(PWD):/app -w /app crpi-6pj79y7ddzdpexs8.cn-hangzhou.personal.cr.aliyuncs.com/gromitlee/golang:1.24.6-bookworm-swag1.16.6 bash -c 'make doc-swag'

doc-swag:
	# swag version v1.16.4
	# v1
	swag fmt  -g guest.go -d internal/bff-service/server/http/handler/v1
	swag init -g guest.go -d internal/bff-service/server/http/handler/v1 -o docs/v1 --md docs --pd
	# openapi
	swag fmt  -g openapi.go -d internal/bff-service/server/http/handler/openapi
	swag init -g openapi.go -d internal/bff-service/server/http/handler/openapi -o docs/openapi --pd
	# callback
	swag fmt  -g callback.go -d internal/bff-service/server/http/handler/callback
	swag init -g callback.go -d internal/bff-service/server/http/handler/callback -o docs/callback --pd
	# openurl
	swag fmt  -g openurl.go -d internal/bff-service/server/http/handler/openurl
	swag init -g openurl.go -d internal/bff-service/server/http/handler/openurl -o docs/openurl --pd

docker-image-backend:
	docker build -f Dockerfile.backend --build-arg WANWU_ARCH=${WANWU_ARCH} -t wanwulite/wanwu-backend:${WANWU_VERSION}-$(shell git rev-parse --short HEAD)-${WANWU_ARCH} .

docker-image-frontend:
	docker build -f Dockerfile.frontend --build-arg WANWU_ARCH=${WANWU_ARCH} -t wanwulite/wanwu-frontend:${WANWU_VERSION}-$(shell git rev-parse --short HEAD)-${WANWU_ARCH} .

docker-image-rag:
	docker build -f Dockerfile.rag --build-arg WANWU_ARCH=${WANWU_ARCH} -t wanwulite/rag:${WANWU_VERSION}-$(shell git rev-parse --short HEAD)-${WANWU_ARCH} .

docker-image-agent:
	docker build -f Dockerfile.agent --build-arg WANWU_ARCH=${WANWU_ARCH} -t wanwulite/agent:${WANWU_VERSION}-$(shell git rev-parse --short HEAD)-${WANWU_ARCH} .

docker-image-agent-base:
	docker build -f Dockerfile.agent-base --build-arg WANWU_ARCH=${WANWU_ARCH} -t wanwulite/agent-base:${WANWU_VERSION}-$(shell git rev-parse --short HEAD)-${WANWU_ARCH} .

grpc-protoc:
	protoc --proto_path=. --go_out=paths=source_relative:api --go-grpc_out=paths=source_relative:api proto/*/*.proto

i18n-jsonl:
	go test ./pkg/i18n -run TestI18nConvertXlsx2Jsonl

init:
	go mod tidy
	go mod vendor

pb:
	docker run --name golang-grpc --privileged=true --rm -v $(PWD):/app -w /app crpi-6pj79y7ddzdpexs8.cn-hangzhou.personal.cr.aliyuncs.com/gromitlee/golang:1.24.6-bookworm-protoc29.4-gengo1.34.1-gengrpc1.5.1-gengw2.20.0-genapi2.20.0 bash -c 'make grpc-protoc'

# --- mysql ---
run-mysql:
	docker compose -f docker-compose.yaml \
		--env-file .env.image.${WANWU_ARCH} \
		--env-file .env \
		up -d mysql

stop-mysql:
	docker compose -f docker-compose.yaml \
		--env-file .env.image.${WANWU_ARCH} \
		--env-file .env \
		down mysql

# --- mysql-setup ---
run-mysql-setup:
	docker compose -f docker-compose.yaml \
		--env-file .env.image.${WANWU_ARCH} \
		--env-file .env \
		up mysql-setup

stop-mysql-setup:
	docker compose -f docker-compose.yaml \
		--env-file .env.image.${WANWU_ARCH} \
		--env-file .env \
		down mysql-setup

# --- redis ---
run-redis:
	docker compose -f docker-compose.yaml \
		--env-file .env.image.${WANWU_ARCH} \
		--env-file .env \
		up -d redis

stop-redis:
	docker compose -f docker-compose.yaml \
		--env-file .env.image.${WANWU_ARCH} \
		--env-file .env \
		down redis

# --- minio ---
run-minio:
	docker compose -f docker-compose.yaml \
		--env-file .env.image.${WANWU_ARCH} \
		--env-file .env \
		up -d minio

stop-minio:
	docker compose -f docker-compose.yaml \
		--env-file .env.image.${WANWU_ARCH} \
		--env-file .env \
		down minio

# --- kafka ---
run-kafka:
	docker compose -f docker-compose.yaml \
		--env-file .env.image.${WANWU_ARCH} \
		--env-file .env \
		up -d kafka

stop-kafka:
	docker compose -f docker-compose.yaml \
		--env-file .env.image.${WANWU_ARCH} \
		--env-file .env \
		down kafka

# --- elastic-setup ---
run-es-setup:
	docker compose -f docker-compose.yaml \
		--env-file .env.image.${WANWU_ARCH} \
		--env-file .env \
		up -d es-setup

stop-es-setup:
	docker compose -f docker-compose.yaml \
		--env-file .env.image.${WANWU_ARCH} \
		--env-file .env \
		down es-setup

# --- elastic ---
run-es:
	docker compose -f docker-compose.yaml \
		--env-file .env.image.${WANWU_ARCH} \
		--env-file .env \
		up -d es

stop-es:
	docker compose -f docker-compose.yaml \
		--env-file .env.image.${WANWU_ARCH} \
		--env-file .env \
		down es