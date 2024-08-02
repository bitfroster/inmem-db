FROM golang:1.22.4-alpine AS golang

WORKDIR /app

COPY . .

RUN go test -v