SHELL:=/bin/sh
TEST?=./...
TESTARGS=-test.v
FUNC=
TZ?=Asia/Tokyo
HOST=pi1

default: build

build: fmtcheck
	@docker compose build

run: fmtcheck
	docker compose up --build

test_up: test_down
	@COMPOSE_FILE="docker-compose.yml:docker-compose-dev.yml" \
		docker compose build
	@COMPOSE_FILE="docker-compose.yml:docker-compose-dev.yml" \
		docker compose run --rm uploader

test_down:
	@COMPOSE_FILE="docker-compose.yml:docker-compose-dev.yml" \
		docker compose down

################################
# Below commmands are to be used in container

test:
	CGO_ENABLED=0 gotest $(TEST) $(TESTARGS) -coverprofile=cover.out.tmp
	@cat cover.out.tmp | grep -v "mock_" > cover.out
	@go tool cover -html=cover.out -o cover.html
	@rm -f cover.out cover.out.tmp

hottest:
ifndef FUNC
	reflex -r '\.go$$' -- sh -c 'CGO_ENABLED=0 gotest $(TEST) $(TESTARGS)'
else
	reflex -r '\.go$$' -- sh -c 'CGO_ENABLED=0 gotest -run $(FUNC) $(TEST) $(TESTARGS)'
endif

generate:
	go generate ./...

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w .

fmtcheck:
	sh -c "'$(CURDIR)/../scripts/gofmtcheck.sh'"

sync:
	sh -c "'$(CURDIR)/../scripts/rsync.sh' $(HOST)"

.PHONY: default build run test_up test_down test hottest generate fmt fmtcheck sync
