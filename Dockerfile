FROM golang:1.20.3-alpine as build
ENV CGO_ENABLED 0
COPY . /concourse-sonarqube-notifier

RUN mkdir -p /assets \
 && cd /concourse-sonarqube-notifier \
 && go test -v ./... \
 && go build -o /assets/in assets/in/main/in.go \
 && go build -o /assets/out assets/out/main/out.go \
 && go build -o /assets/check assets/check/main/check.go

FROM alpine:3.17.3 AS runtime
RUN apk add --no-cache ca-certificates
COPY --from=build assets/ /opt/resource/
RUN chmod +x /opt/resource/*
