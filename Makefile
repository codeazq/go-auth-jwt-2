build: 
	@go build -o bin/go-fintech-bank

run: build
	@./bin/go-fintech-bank

test: 
	@go test -v ./...