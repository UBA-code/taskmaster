.PHONY: build run clean test

BINARY=./bin/taskmaster

build:
	@mkdir -p bin
	go build -o $(BINARY) ./cmd/taskmaster

run: build
	$(BINARY) example.yml

clean:
	rm -rf bin/
