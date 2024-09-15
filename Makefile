build:
	@go build -o bin/url-shortner cmd/main.go 

run: build
	@./bin/ecom