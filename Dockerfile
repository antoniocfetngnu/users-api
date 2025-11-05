FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy ALL source code including generated GraphQL files ✅
COPY . .

# Just build, no generation needed ✅
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/graphql/schema.graphql ./graphql/
RUN mkdir -p /root/data
EXPOSE 8080
CMD ["./main"]