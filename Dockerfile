FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate GraphQL code
RUN go run github.com/99designs/gqlgen generate --config graphql/gqlgen.yml

# Generate Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o users-service .

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy binary
COPY --from=builder /app/users-service .

# Expose port
EXPOSE 3001

CMD ["./users-service"]