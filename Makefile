test:
	go test ./...
install:
	go install ./...
fmt:
	goimports -w .
	gofmt -w .
