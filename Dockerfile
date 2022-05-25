FROM golang:1.18.0-alpine3.15 as builder

WORKDIR app

RUN apk update \
  && apk add --no-cache git

COPY . .

RUN go mod download

ARG CGO_ENABLED=0
ARG GOOS=linux
ARG GOARCH=amd64
RUN go build \
    -o /go/app/main \
    -ldflags '-s -w'


FROM alpine:latest as runner

WORKDIR app

COPY --from=builder /go/app/main /app/main
COPY --from=builder /go/app/.env /app/.env
COPY --from=builder /go/app/resource/ /app/resource/

ENTRYPOINT "/app/main"