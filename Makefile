GOVERSION ?= $(shell go list -m -f '{{.GoVersion}}')

.PHONY: build
build: clean
	@CGO_ENABLED=0 go build -o ./bin/ ./cmd/*

.PHONY: build/container
build/container:
	@docker buildx build . -t slewiskelly/ock --build-arg GOVERSION=$(GOVERSION)

.PHONY: clean
clean:
	@rm -rf ./bin

.PHONY: test
test:
	@go test -race ./...

.PHONY: vet
vet:
	@go vet ./...
