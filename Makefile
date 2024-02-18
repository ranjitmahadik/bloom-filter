build:
	@go build -o ./cmd/bloom

run: build
	@docker compose up -d
	@./cmd/bloom
	@make cleanup

cleanup:
	@docker compose down
	@rm -rf ./cmd