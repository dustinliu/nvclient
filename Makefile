APP = nvclient

.PHONY: build run clean test release
build: test
	@go mod tidy
	@go build -o $(APP) main.go

test:
	@go test -timeout 10s ./...


clean:
	@go clean
	@go clean -testcache

release:
	@goreleaser release
