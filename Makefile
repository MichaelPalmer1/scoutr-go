test:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...

cover: test
	go tool cover -func=coverage.txt
