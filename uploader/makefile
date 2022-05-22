SHELL:=/bin/sh
TEST?=./...
TESTARGS=-test.v
FUNC=
TZ?=Asia/Tokyo
HOST=pi1

default: build

build: fmtcheck
	docker build --target build-stage -t dev-img --build-arg TZ=$(TZ) .
	docker build -t prod-img --build-arg TZ=$(TZ) .

run: build
	docker run -ti --rm --net=host \
		-v $(PWD)/tmp/data:/var/data \
    -v ${PWD}/crt/ca.crt:${CA_PATH} \
    -v ${PWD}/crt/client.crt:${CERT_PATH} \
    -v ${PWD}/crt/client.key:${KEY_PATH}
		--name gps_uploader prod-img:latest

test_up: test_down
	docker run -ti --rm --net=host \
		-v $(PWD)/:/opt/gps_uploader/ \
		-v $(PWD)/tmp/data:/var/data \
    -v ${PWD}/crt/ca.crt:${CA_PATH} \
    -v ${PWD}/crt/client.crt:${CERT_PATH} \
    -v ${PWD}/crt/client.key:${KEY_PATH}
		--name gps_uploader dev-img:latest

test_down:
	docker rm -f build-server

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
