.PHONY: all
all: build

.PHONY: build
build:
	go build -o test ./cmd/mdtt

install:
	go install ./cmd/mdtt

.PHONY: test
test:
	go test ./...

.PHONY: clean
clean:
	rm -f test
