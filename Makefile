IMAGE_NAME = weather-cep:1.0.0
CONTAINER_NAME = weather-cep
PORT = 3099


.PHONY: goserver
goserver:
	go run cmd/server/main.go


.PHONY: server up down stop test stopall _up _down _stop _stopall

# usage: make server {up|down|stop|stopall}
up down stop stopall:
	@true

server:
	@action="$(filter-out server,$(MAKECMDGOALS))"; \
	if [ -z "$$action" ]; then \
		echo "Usage: make server {up|down|stop|stopall}"; \
		exit 1; \
	fi; \
	$(MAKE) --no-print-directory _$$action

test: _up
	@echo "Running automated tests (builds and runs the app container)..."
	go test -tags e2e ./test/e2e/... -v -timeout 5m -count=1

_up:
	@if docker inspect -f '{{.State.Running}}' $(CONTAINER_NAME) 2>/dev/null | grep -q '^true$$'; then \
		echo "Container $(CONTAINER_NAME) is already running."; \
	else \
		echo "Starting the container..."; \
		if ! docker image inspect $(IMAGE_NAME) >/dev/null 2>&1; then \
			echo "Image $(IMAGE_NAME) not found, building..."; \
			docker build -t $(IMAGE_NAME) .; \
		else \
			echo "Using existing image $(IMAGE_NAME)."; \
		fi; \
		docker rm -f $(CONTAINER_NAME) >/dev/null 2>&1 || true; \
		docker run -d --name $(CONTAINER_NAME) -p $(PORT):8080 --env-file .env -e PORT=8080 $(IMAGE_NAME); \
		echo "Container $(CONTAINER_NAME) started and listening on port $(PORT)."; \
	fi

_down:
	@echo "Stopping the container (cleaning up)..."
	@if docker inspect -f '{{.State.Running}}' $(CONTAINER_NAME) 2>/dev/null | grep -q '^true$$'; then \
		docker stop $(CONTAINER_NAME); \
	fi
	@docker rm $(CONTAINER_NAME) >/dev/null 2>&1 || true
	@docker rmi -f $(IMAGE_NAME) >/dev/null 2>&1 || true

_stop:
	@echo "Stopping the container..."
	@docker rm -f $(CONTAINER_NAME) >/dev/null 2>&1 || true

_stopall:
	@echo "Stopping all containers..."
	@docker rm -f $(shell docker ps -aq) >/dev/null 2>&1 || true
