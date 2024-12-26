GOBASE 		:= $(shell pwd)
GOBIN		:= $(GOBASE)/bin

build-server:
	@go build -o $(GOBIN)/server/app ./cmd/server/main.go

run-server: build-server
	@./bin/server/app

build-client:
	@go build -o $(GOBIN)/client/app ./cmd/client/main.go

run-client: build-client
	@./bin/client/app

clean:
	@rm -rf $(GOBIN)

.PHONY: build-server run-server build-client run-client clean
