.PHONY: build test clean install lint

build:
	go build -o weathercli cmd/weathercli/main.go

install:
	go install ./cmd/weathercli

test:
	go test -v ./...

coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run

clean:
	rm -f weathercli coverage.out coverage.html
	go clean
