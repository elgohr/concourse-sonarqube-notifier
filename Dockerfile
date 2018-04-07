FROM golang:alpine as build
ENV CGO_ENABLED 0
COPY . /go/src/github.com/concourse-sonarqube-notifier

RUN mkdir -p /assets \
 && go build -o /assets/in github.com/concourse-sonarqube-notifier/assets/in/main \
 && go build -o /assets/out github.com/concourse-sonarqube-notifier/assets/out/main \
 && go build -o /assets/check github.com/concourse-sonarqube-notifier/assets/check/main

FROM alpine AS runtime
RUN apk add --no-cache ca-certificates
COPY --from=build assets/ /opt/resource/
RUN chmod +x /opt/resource/*
