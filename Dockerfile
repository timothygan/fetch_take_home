FROM golang:1.23 AS builder
LABEL authors="timgan"

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go build -o app cmd/server/main.go
CMD ["./app"]