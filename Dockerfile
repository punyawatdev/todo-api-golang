# Stage 1: Build the binary
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy file go.mod และ go.sum for caching dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build App (--> cmd/api/main.go)
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# Stage 2: Final image
FROM alpine:latest
RUN adduser -D appuser
USER appuser

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080
CMD ["./main"]