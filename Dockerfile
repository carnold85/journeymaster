FROM golang:1.23.4 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o journeymastet

FROM alpine:latest
COPY --from=builder /app/journeymastet /usr/local/bin/journeymastet
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/journeymastet"]