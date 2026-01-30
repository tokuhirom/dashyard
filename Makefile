.PHONY: all frontend backend build test clean dev-frontend dev-backend dev-dummyprom

all: build

frontend:
	cd frontend && npm run build

backend: frontend
	go build -o dashyard .

build: backend

test:
	go test ./...

dev-frontend:
	cd frontend && npm run dev

dev-backend:
	go run . -config examples/config.yaml

dev-dummyprom:
	go run ./cmd/dummyprom

clean:
	rm -f dashyard dummyprom
	rm -rf frontend/dist
