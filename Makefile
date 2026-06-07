BINARY ?= intervals-mcp

.PHONY: build run clean

build:
	go build -o $(BINARY) ./cmd

run:
	go run ./cmd

clean:
	rm -f $(BINARY)