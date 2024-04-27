#!make

fmt:
	go fmt ./...

test:
	go test -v ./...

test-integration:
	docker compose -f docker-compose.test.yaml down
	docker compose -f docker-compose.test.yaml up --build --abort-on-container-exit --exit-code-from assessment-tax; \
	docker compose -f docker-compose.test.yaml down

coverage:
	go test -cover -coverprofile=report.out -v ./...
	go tool cover -html=report.out -o coverage.html

clean:
	go clean -i ./...

gen-doc:
	swag fmt
	swag init

run-dev:
	docker compose -f docker-compose.dev.yaml down
	docker compose -f docker-compose.dev.yaml up --build --renew-anon-volumes --abort-on-container-exit --exit-code-from assessment-tax; \
	docker compose -f docker-compose.dev.yaml down

run-release:
	docker compose -f docker-compose.release.yaml down
	docker compose -f docker-compose.release.yaml up --build -d

run:
	docker compose -f docker-compose.yaml up -d
	PORT="8080" DATABASE_URL=host="localhost port=5432 user=postgres password=postgres dbname=ktaxes sslmode=disable" ADMIN_USERNAME="adminTax" ADMIN_PASSWORD="admin!" go run main.go