FROM golang:1.13 AS builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN env CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"'

FROM alpine
COPY --from=builder /usr/src/app/do-api-proxy /
CMD ["/do-api-proxy"]
EXPOSE 1338/tcp
