APP_NAME = rssagg

run: build
	@./bin/$(APP_NAME)

build:
	@go build -o ./bin/$(APP_NAME) ./cmd/$(APP_NAME)/
