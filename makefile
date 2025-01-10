# Define variables
DOCKER_COMPOSE_DEV = docker-compose.dev.yml
DOCKER_COMPOSE_PROD = docker-compose.prod.yml

# Development targets
dev-up:
	docker-compose -f $(DOCKER_COMPOSE_DEV) up --build

dev-down:
	docker-compose -f $(DOCKER_COMPOSE_DEV) down

dev-logs:
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
	docker-compose ps

clean:
	docker-compose -f $(DOCKER_COMPOSE_DEV) down -v
	docker-compose -f $(DOCKER_COMPOSE_PROD) down -v
