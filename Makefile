.PHONY: run build 

BINARY := bin/recl
MAIN := ./cmd/recl


build:
	mkdir -p $(dir $(BINARY))
	go build -o $(BINARY) $(MAIN)

run:
	go run $(MAIN)
