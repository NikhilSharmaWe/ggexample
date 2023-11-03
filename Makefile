build:
	@go build -o bin/ggexample

run: build
	@./bin/ggexample
