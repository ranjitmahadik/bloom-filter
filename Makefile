build:
	@go build -o ./cmd/bloom

run: build
	docker compose down
	docker compose up -d
	@./cmd/bloom