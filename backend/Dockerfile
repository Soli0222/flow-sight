# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies and swag tool
RUN apk add --no-cache git
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate Swagger documentation
RUN swag init -g cmd/main.go -o docs

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -a -installsuffix cgo -o main ./cmd/main.go

# Final stage
FROM alpine:3.18

RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/docs ./docs

# Expose port
EXPOSE 8080

# Command to run
CMD ["./main"]
