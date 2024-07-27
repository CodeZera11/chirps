build:
	@go build -o bin/chirps

run: build
	@./bin/chirps

test:
	go test -v ./...