FROM golang:1.24.5-alpine AS builder

WORKDIR /app
COPY . . 

RUN go build -o myapp


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/myapp .
COPY config.json .

RUN adduser -D appuser
RUN chown -R appuser:appuser /app

USER appuser

ENTRYPOINT ["./myapp"]
