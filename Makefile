.PHONY: build run

build:
	CGO_ENABLED=0 go build -v

run:
	./investment-balancer-v3
