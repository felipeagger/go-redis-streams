docker:
	@docker-compose down
	@docker-compose build
	@docker-compose up -d

dockerdown:
	@docker-compose down

build:
	@echo "---- Building Application ----"
	@go build -o consumer consumer/*.go
	@go build -o producer producer/*.go

consume:
	@echo "---- Running Consumer ----"
	@export REDIS_HOST=localhost
	@export STREAM=events
	#@export GROUP=GroupOne
	@go run consumer/*.go

run:
	@echo "---- Running Producer ----"
	@export REDIS_HOST=localhost
	@export STREAM=events
	@go run producer/*.go