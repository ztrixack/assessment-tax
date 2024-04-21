#!make

fmt:
	go fmt ./...

test:
	go test -v ./...

test-integration:
	docker-compose -f docker-compose.test.yaml up --build --abort-on-container-exit --exit-code-from assessment-tax
	docker-compose -f docker-compose.test.yaml down

coverage:
	go test -cover -coverprofile=c.out -v ./...
	go tool cover -html=c.out -o coverage.html

clean:
	go clean -i ./...
