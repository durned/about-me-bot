# Define the docker-compose command
DC=docker-compose
PROJECT_DIR := ./about-me-bot
CMD_DIR := $(PROJECT_DIR)/cmd

# Define commands
.PHONY: local deployed stop clean

local: stop
	@echo "Starting MongoDB and running Go application locally..."
	$(DC) up -d mongo
	cd $(CMD_DIR) && go run main.go

deployed:
	@echo "Starting MongoDB and Go application in Docker containers..."
	$(DC) up --build

stop:
	@echo "Stopping containers..."
	$(DC) down || true

clean:
	@echo "Stopping and cleaning the containers..."
	$(DC) down -v || true