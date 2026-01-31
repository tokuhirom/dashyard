.PHONY: all frontend backend build test test-e2e lint clean dev-frontend dev-backend dev-dummyprom screenshots gen-prompt

all: build

frontend:
	cd frontend && npm run build

backend: frontend
	go build -o dashyard .

build: backend

test:
	go test ./...

test-e2e:
	cd frontend && npx playwright test

lint:
	golangci-lint run ./...

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
	@echo "Starting dummyprom..."
	@go run ./cmd/dummyprom & DUMMYPROM_PID=$$!; \
	sleep 1; \
	go run . gen-prompt http://localhost:9090 -o examples/gen-prompt-example.md; \
	kill $$DUMMYPROM_PID 2>/dev/null; \
	echo "Generated examples/gen-prompt-example.md"

clean:
	rm -f dashyard dummyprom
	rm -rf frontend/dist
