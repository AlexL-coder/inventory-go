# Define variables
DOCKER_COMPOSE_DEV = docker-compose.dev.yml
DOCKER_COMPOSE_PROD = docker-compose.prod.yml
GOLANGCI_LINT = golangci-lint
INSTALL_LINT = go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
GOSEC = gosec
INSTALL_GOSEC = go install github.com/securego/gosec/v2/cmd/gosec@latest

# Default target
.DEFAULT_GOAL := help

# Lint targets
lint-install: ## Install golangci-lint if not already installed
	@if not exist "$(GOLANGCI_LINT)" ( \
		echo "golangci-lint not found. Installing..."; \
		$(INSTALL_LINT); \
	) else ( \
		echo "golangci-lint already installed."; \
	)

lint: ## Run golangci-lint
	echo "Running golangci-lint..."
	$(GOLANGCI_LINT) run --verbose ./...

lint-fix: ## Automatically fix linting issues
	$(GOLANGCI_LINT) run --fix ./...

# Security Scan
gosec-install: ## Install gosec if not already installed
	@if [ -z "$(GOSEC)" ]; then \
		echo "gosec not found. Installing..."; \
		$(INSTALL_GOSEC); \
	else \
		echo "gosec already installed."; \
	fi

gosec-scan: ## Run gosec security scan
	$(GOSEC) ./...

gosec-report: ## Run gosec and generate an HTML report
	$(GOSEC) -fmt=html -out=gosec-report.html ./...
	@echo "Security report generated at gosec-report.html"

# Development targets
dev-up: lint  gosec-scan gosec-report ## Lint the code and start the development environment
	@echo "Starting development environment..."
	docker-compose -f $(DOCKER_COMPOSE_DEV) up --build

dev-down: ## Stop the development environment
	@echo "Stopping development environment..."
	docker-compose -f $(DOCKER_COMPOSE_DEV) down

dev-logs: ## Show logs for development environment
	@echo "Fetching logs for development environment..."
	docker-compose -f $(DOCKER_COMPOSE_DEV) logs -f


# Production targets
prod-up:
	docker-compose -f $(DOCKER_COMPOSE_PROD) up --build -d

prod-down:
	docker-compose -f $(DOCKER_COMPOSE_PROD) down

prod-logs:
	docker-compose -f $(DOCKER_COMPOSE_PROD) logs -f

# General targets
ps:
	@echo "Listing running containers..."
	docker-compose ps

clean:
	@echo "Cleaning up development and production environments..
	docker-compose -f $(DOCKER_COMPOSE_DEV) down -v
	docker-compose -f $(DOCKER_COMPOSE_PROD) down -v

# Help target
help: ## Show this help message
	@echo "Available targets:"
	@type $(MAKEFILE_LIST) | findstr "##"