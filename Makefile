build:
	@go build -o bin/api main.go
run: build
	@./bin/api
seed:
	@go run scripts/seed.go