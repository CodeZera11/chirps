build:
	@go build -o bin/chirps

run: build
	@rm -r database.json && ./bin/chirps

test:
	go test -v ./... --debug