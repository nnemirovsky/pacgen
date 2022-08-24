FROM golang:alpine AS build
WORKDIR /app
COPY . .
RUN apk add --no-cache build-base
RUN go mod download
RUN make build


FROM alpine:latest
VOLUME /app/data
WORKDIR /app
COPY --from=build /app/bin/* .
COPY entrypoint.sh .
COPY migrations migrations
RUN apk add --no-cache curl
USER 1001:1001
CMD ["./server"]
