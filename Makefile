build:
	docker-compose build grpc-task

run:
	docker-compose up grpc-task

test:
	go test -v ./...

lint:
	golangci-lint run