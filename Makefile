.PHONY: build run build-prod compress test coverage coverage-open

build:
	CGO_ENABLED=0 go build -v

build-prod:
	CGO_ENABLED=0 go build -v -o investment-balancer-v3-prod -ldflags="-w -s -buildid=" -trimpath

run:
	./investment-balancer-v3

compress-prod:
	rm investment-balancer-v3-compressed
	upx --best -o ./investment-balancer-v3-compressed investment-balancer-v3

test:
	go test ./...

coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

coverage-open:
	xdg-open coverage.html
