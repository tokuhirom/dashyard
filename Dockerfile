## Stage 1: Build frontend
FROM node:20-alpine AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm ci
COPY frontend/ .
RUN npm run build

## Stage 2: Build backend
FROM golang:1.25-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist
RUN CGO_ENABLED=0 go build -o dashyard .

## Stage 3: Runtime
FROM alpine:3.20
RUN apk --no-cache add ca-certificates && \
    addgroup -S dashyard && adduser -S dashyard -G dashyard
WORKDIR /app
COPY --from=backend-builder /app/dashyard .
USER dashyard
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s CMD wget -qO- http://localhost:8080/ready || exit 1
ENTRYPOINT ["/app/dashyard"]
CMD ["serve", "--config", "/etc/dashyard/config.yaml", "--dashboards-dir", "/etc/dashyard/dashboards"]
