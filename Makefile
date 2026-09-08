.PHONY: build test vet fmt lint check install clean

build:
	go build -o cca ./cmd/cca

test:
	go test ./... -race -cover

vet:
	go vet ./...

fmt:
	@out=$$(gofmt -l .); \
	if [ -n "$$out" ]; then \
		echo "gofmt needs to be run on:"; echo "$$out"; exit 1; \
	fi

lint:
	golangci-lint run

check: fmt vet test lint

install:
	go run ./cmd/cca install

clean:
	rm -f cca cca.exe
