#go test -p 1 -covearprofile=coverage.out ./dataaccess/...
run:
	go run main.go
clean:
	rm roam
build:
	go build
cover:
	go test -p 1 -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
lint:
	golangci-lint run
fix:
	golangci-lint run --fix
generate:
	mockery --all