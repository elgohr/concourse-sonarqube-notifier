FROM golang as build
ENV CGO_ENABLED 0
COPY . /concourse-sonarqube-notifier

RUN mkdir -p /assets \
 && cd /concourse-sonarqube-notifier \
 && go get -u github.com/maxbrunsfeld/counterfeiter \
 && go generate ./... \
 && go test -v ./... \
 && go build -o /assets/in assets/in/main/in.go \
 && go build -o /assets/out assets/out/main/out.go \
 && go build -o /assets/check assets/check/main/check.go

FROM alpine AS runtime
RUN apk add --no-cache ca-certificates
COPY --from=build assets/ /opt/resource/
RUN chmod +x /opt/resource/*
