BINARY_NAME=sentinel
GO=go

.PHONY: all build clean fmt vet test test-race test-cover ci run docker-run

all: build

build:
	$(GO) build -o $(BINARY_NAME) ./cmd/sentinel

run: build
	./$(BINARY_NAME)

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

test:
	$(GO) test ./...

test-race:
	$(GO) test -race ./...

test-cover:
	$(GO) test -cover ./...

ci: fmt vet test test-race test-cover

clean:
	$(GO) clean
	rm -f $(BINARY_NAME)

docker-run:
	docker compose up --build
