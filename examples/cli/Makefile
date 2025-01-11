.PHONY: build run build-prod compress test coverage coverage-open

build:
	CGO_ENABLED=0 go build -v

build-prod:
	CGO_ENABLED=0 go build -v -o cli-prod -ldflags="-w -s -buildid=" -trimpath

run:
	./cli

compress-prod:
	rm cli-compressed
	upx --best -o ./cli-compressed cli

test:
	go test ./...

coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

coverage-open:
	xdg-open coverage.html
