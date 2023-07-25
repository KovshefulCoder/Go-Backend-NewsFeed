FROM golang:alpine AS builder
COPY . /build/

WORKDIR /build
ARG CGO_ENABLED=0
ARG GOOS=linux

RUN go build -installsuffix 'static' -o app server.go

FROM alpine:latest
COPY --from=builder /build/app .
EXPOSE 8082 8082
ENTRYPOINT ["./app"]
