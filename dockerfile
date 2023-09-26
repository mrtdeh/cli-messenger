FROM golang:alpine as build-env

ENV GO111MODULE=on

RUN apk update && apk add bash ca-certificates git gcc g++ libc-dev

RUN mkdir /api-channel
RUN mkdir -p /api-channel/proto 

WORKDIR /api-channel

COPY ./proto/service.pb.go /api-channel/proto
COPY ./main.go /api-channel

COPY go.mod .
COPY go.sum .

RUN go mod download

RUN go build -o api-channel .

CMD ./api-channel