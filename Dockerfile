FROM golang:alpine AS builder
RUN mkdir /app
ADD ./ /app/
WORKDIR /app
RUN apk update && \
    apk add --no-cache alpine-sdk && \
    GOOS=linux GOARCH=amd64 go build -tags musl -a -o producer-enhancer .
FROM alpine
RUN apk add --no-cache tzdata bash && addgroup -S onms && adduser -S -G onms onms
COPY --from=builder /app/producer-enhancer /producer-enhancer
COPY ./docker-entrypoint.sh /
USER onms
LABEL maintainer="Alejandro Galue <agalue@opennms.org>" \
      name="OpenNMS Kafka Producer enhancer"
ENTRYPOINT [ "/docker-entrypoint.sh" ]
