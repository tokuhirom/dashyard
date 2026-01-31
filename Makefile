.PHONY: all frontend backend build test test-frontend test-e2e lint clean dev-frontend dev-backend dev-dummyprom screenshots gen-prompt gen-prompt-up

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
	cd frontend && npx playwright test

lint:
	golangci-lint run ./...
	cd frontend && npm run lint

dev-frontend:
	cd frontend && npm run dev

dev-backend:
	go run . serve --config examples/config.yaml

dev-dummyprom:
	go run ./cmd/dummyprom

screenshots:
	docker compose -f docker-compose.screenshots.yaml up --build --abort-on-container-exit screenshots
	docker compose -f docker-compose.screenshots.yaml down

gen-prompt:
	docker compose -f docs/gen-prompt/docker-compose.yaml up -d prometheus otelcol traefik redis whoami traffic-gen
	@echo "Waiting 60s for metrics to accumulate..."
	@sleep 60
	docker compose -f docs/gen-prompt/docker-compose.yaml build gen-prompt
	docker compose -f docs/gen-prompt/docker-compose.yaml run --rm gen-prompt
	docker compose -f docs/gen-prompt/docker-compose.yaml down

gen-prompt-up:
	docker compose -f docs/gen-prompt/docker-compose.yaml up --build

clean:
	rm -f dashyard dummyprom
	rm -rf frontend/dist
