# Build stage
FROM golang:1.23.2-alpine AS builder

WORKDIR /app

# Install necessary tools
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Скопировать исходный код прежде чем вызывать go mod tidy и go mod vendor
COPY . .

# Теперь, когда код уже в контейнере, можно подтянуть зависимости
RUN go mod tidy && go mod vendor

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose the port the app runs on
EXPOSE 3006

# Command to run the application
CMD ["./main", "api"]