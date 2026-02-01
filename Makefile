.PHONY: all frontend backend build test test-frontend test-e2e lint clean screenshots gen-prompt gen-prompt-up

all: build

frontend:
	cd frontend && npm run build

backend: frontend
	go build -o dashyard .

build: backend

test:
	go test ./...

test-frontend:
	cd frontend && npm run test

test-e2e:
	docker compose -f docker-compose.e2e.yaml up --build --exit-code-from e2e e2e; \
	rc=$$?; \
	docker compose -f docker-compose.e2e.yaml down; \
	exit $$rc

lint:
	golangci-lint run ./...
	cd frontend && npm run lint

screenshots:
	docker compose -f docker-compose.screenshots.yaml up --build --abort-on-container-exit screenshots
	docker compose -f docker-compose.screenshots.yaml down

gen-prompt:
	docker compose -f examples/real-world/docker-compose.yaml up -d prometheus otelcol traefik redis whoami dummyapp traffic-gen
	@echo "Waiting 60s for metrics to accumulate..."
	@sleep 60
	docker compose -f examples/real-world/docker-compose.yaml build gen-prompt
	docker compose -f examples/real-world/docker-compose.yaml run --rm gen-prompt
	docker compose -f examples/real-world/docker-compose.yaml down

gen-prompt-up:
	docker compose -f examples/real-world/docker-compose.yaml up --build

clean:
	rm -f dashyard dummyprom
	rm -rf frontend/dist
