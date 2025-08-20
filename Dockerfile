FROM golang:1.24.3 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG APP_PORT
ARG DB_HOST
ARG DB_PORT
ARG DB_USER
ARG DB_PASSWORD
ARG DB_NAME

ENV APP_PORT=$APP_PORT \
    DB_HOST=$DB_HOST \
    DB_PORT=$DB_PORT \
    DB_USER=$DB_USER \
    DB_PASSWORD=$DB_PASSWORD \
    DB_NAME=$DB_NAME

RUN go build -o main ./cmd/main.go

FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE ${APP_PORT}
CMD ["./main"]
