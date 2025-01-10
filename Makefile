test:
	go test ./...

coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

coverage-open:
	xdg-open file://$(PWD)/coverage.html
