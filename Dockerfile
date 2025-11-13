FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY users-api /app/
WORKDIR /app
EXPOSE 3001
CMD ["./users-api"]
