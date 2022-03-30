FROM golang:1.17-alpine as builder
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY ./internal ./internal
COPY ./cmd ./cmd

RUN go build  -o bin/url-collector ./cmd/nasa/main.go

FROM alpine:latest
MAINTAINER Artur Obrzut <artek.obrzut@gmail.com>

ENV API_KEY DEMO_KEY
ENV PORT 8080
ENV CONCURRENT_REQUESTS 5
ENV GIN_MODE=release
WORKDIR /app
COPY --from=builder /app/bin/url-collector ./url-collector

RUN addgroup -S nasa \
&& adduser -D -S  -s /sbin/nologin -G nasa -u 10001 nasa \
&& chown -R nasa:nasa /app/url-collector \
&& apk add libcap \
&& setcap 'cap_net_bind_service=+ep' /app/url-collector

EXPOSE ${PORT}
USER 10001

CMD [ "/app/url-collector" ]