.PHONY: build run

build:
	CGO_ENABLED=0 go build -v

build-prod:
	CGO_ENABLED=0 go build -v -o investment-balancer-v3-prod -ldflags="-w -s -buildid=" -trimpath

run:
	./investment-balancer-v3

test:
	go test ./...

coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

coverage-open:
	xdg-open coverage.html
