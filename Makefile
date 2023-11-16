build:
	@go build -o bin/demoQuizApp

run: build
	@./bin/demoQuizApp
