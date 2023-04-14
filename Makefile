build:
	docker-compose build youtube-thumbnails-downloader

run:
	docker-compose up youtube-thumbnails-downloader

test:
	go test -v ./...

lint:
	golangci-lint run