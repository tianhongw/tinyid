FROM golang:1.12-alpine AS build_base

RUN apk add --no-cache git make bash

WORKDIR /opt/tinyid

COPY go.mod .

RUN go mod download

COPY . .

RUN make build

CMD ["/opt/tinyid/build/tinyid", "-c", "/opt/tinyid/config/conf.toml"]
