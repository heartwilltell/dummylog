.PHONY: test
test:
	go test -cover ./...

.PHONY: test-race
test-race:
	go test -cover -race ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: linux
linux:
	GOOS=linux GOARCH=amd64 go build -o dummylog ./cmd

.PHONY: darwin
darwin:
	GOOS=darwin GOARCH=amd64 go build -o dummylog ./cmd

.PHONY: darwin
windows:
	GOOS=windows GOARCH=amd64 go build -o dummylog ./cmd