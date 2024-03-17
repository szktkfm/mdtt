.PHONY: all
all: build

.PHONY: build
build:
	go build -o test ./cmd/mdtt

install:
	go install ./cmd/mdtt

lint:
	golangci-lint run

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -f test
