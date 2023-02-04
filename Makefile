.PHONY: help bash psql logs
.PHONY: build up down run test

# set default goal to help
.DEFAULT_GOAL := help

## ==============================================
# start set service name

SERVICES := zookeeper kafka receiver screen-shot-service scrapy-splash screen-shot-db screen-shot-api apache-server

.PHONY: $(SERVICES)

$(SERVICES)		: %: ; @:

ARGS1 := $(wordlist 2,2,$(MAKECMDGOALS))

# check if SERVICE in arguments[2] ex.(make logs receiver)
ifneq "$(filter $(ARGS1),$(SERVICES))" ""
SERVICE := $(ARGS1)
endif

# end set service name
## ==============================================


DOCKER_COMPOSE := docker-compose
DOCKER_EXEC := $(DOCKER_COMPOSE) exec
DOCKER_RUN := $(DOCKER_COMPOSE) run --rm


### * make help                                             		Print this help
help: Makefile
	@sed -n 's/^###//p' $<


### * make bash SERVICE_NAME						Log into a service bash
bash:
	$(DOCKER_EXEC) $(SERVICE) bash

### * make logs SERVICE_NAME						Get the logs of a service
logs:
	$(DOCKER_COMPOSE) logs -f $(SERVICE)

### * make test								services tests locally 
test:
	cd ./receiver/tests; go test -v 

psql:
	$(DOCKER_COMPOSE) exec screen-shot-db psql screenshotdb postgres

### * make up_service                                       Up and run specific service
up_service:
	$(DOCKER_COMPOSE) up --build $(SERVICE)

### * make build								Build all docker containers
build:
	- $(DOCKER_COMPOSE) build $(SERVICE)

### * make up								Compose up all docker containers
up:build
	#Creates the required external network
	docker network create -d bridge --subnet 1.1.1.0/24 --gateway 1.1.1.1 host_machine || echo "Network already created."
	$(DOCKER_COMPOSE) up $(SERVICE)

### * make down								Compose down all docker containers
down:
	$(DOCKER_COMPOSE) down $(SERVICE)

### * make run								Down then build then up
run: down build up
