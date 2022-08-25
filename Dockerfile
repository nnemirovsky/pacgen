FROM golang:alpine AS build
WORKDIR /app
COPY . .
RUN apk add --no-cache build-base
RUN go mod download
RUN make build



FROM alpine:latest

LABEL maintainer="Nikita Nemirovsky <vaze.legend@gmail.com>"

VOLUME /app/data
ENV APP_PORT=8080
ENV APP_LOG_LEVEL=debug
ENV APP_USER=admin
ENV APP_PASSWORD=admin

WORKDIR /app
COPY --from=build /app/bin/* .
COPY migrations migrations

RUN apk add --no-cache curl && \
    adduser -D -u 1001 www

USER www:www

CMD ["./server"]
