FROM alpine:latest
RUN apk --no-cache add ca-certificates wget
COPY users-api /app/
WORKDIR /app
EXPOSE 3001
EXPOSE 50051
CMD ["./users-api"]