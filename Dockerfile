FROM golang:alpine as builder

ARG app="wedding-api"

WORKDIR /go/home

RUN apk update \
  && apk add --no-cache git
  # && apk add bash

COPY . .

RUN go mod download

WORKDIR /go/home/app/${app}

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
RUN go build \
  -o /go/home/app/${app}/main \
  -ldflags '-s -w'

FROM alpine:latest as runner
# FROM ubuntu:latest as runner

ARG app="wedding-api"

WORKDIR /home

COPY --from=builder /go/home/app/${app}/main /home/app/${app}/main
COPY --from=builder /go/home/.env /home/.env
COPY --from=builder /go/home/resource/ /home/resource/

WORKDIR /home/app/${app}

ENTRYPOINT "./main"