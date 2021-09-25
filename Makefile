.PHONY: test
test:
	go test -cover ./...

.PHONY: test-race
test-race:
	go test -cover -race ./...

.PHONY: test-d
test-d:
	go test -v -cover ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: lint-d
lint-d:
	golangci-lint run -v ./...

.PHONY: linux
linux:
	GOOS=linux GOARCH=amd64 go build -o dummylog ./cmd

.PHONY: linux-d
linux-d:
	GOOS=linux GOARCH=amd64 go build -v -a -o dummylog ./cmd

.PHONY: darwin
darwin:
	GOOS=darwin GOARCH=amd64 go build -o dummylog ./cmd

.PHONY: darwin
windows:
	GOOS=windows GOARCH=amd64 go build -o dummylog ./cmd

.PHONY: release
release:
	rm -rf ./release

	GOOS=linux GOARCH=amd64 go build -o ./release/dummylog-linux-amd64/dummylog ./cmd
	tar -czvf ./release/dummylog-linux-amd64.tar.gz ./release/dummylog-linux-amd64
	rm -rf ./release/dummylog-linux-amd64

	GOOS=darwin GOARCH=amd64 go build -o ./release/dummylog-darwin-amd64/dummylog ./cmd
	tar -czvf ./release/dummylog-darwin-amd64.tar.gz ./release/dummylog-darwin-amd64
	rm -rf ./release/dummylog-darwin-amd64

	GOOS=windows GOARCH=amd64 go build -o ./release/dummylog-windows-amd64/dummylog ./cmd
	tar -czvf ./release/dummylog-windows-amd64.tar.gz ./release/dummylog-windows-amd64
	rm -rf ./release/dummylog-windows-amd64

