.PHONY: all frontend backend build test test-e2e lint clean dev-frontend dev-backend dev-dummyprom

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

clean:
	rm -f dashyard dummyprom
	rm -rf frontend/dist
