FROM golang:1.16.5-alpine3.13

WORKDIR /

COPY . .
COPY docker.config.json ./config.json

RUN go build

CMD "./push-message-example"