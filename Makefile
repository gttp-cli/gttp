test:
	@echo "Running tests"
	go test -v ./...

server:
	go run . serve

image:
	docker build -t ghcr.io/gttp-cli/gttp .

up:
	docker compose up --build

down:
	docker compose down

