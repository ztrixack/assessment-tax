#!make

fmt:
	go fmt ./...

test:
	go test -v ./...

coverage:
	go test -cover -coverprofile=c.out
	go tool cover -html=c.out -o coverage.html

clean:
	go clean -i ./...