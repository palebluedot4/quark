# syntax=docker/dockerfile:1

ARG GO_VERSION=1.26
ARG ALPINE_VERSION=3.22

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION}

RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download

COPY .air.toml .

EXPOSE 8080

CMD [ "air" ]
