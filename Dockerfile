FROM golang:1.13 AS build

WORKDIR /app
COPY . .

RUN go mod download
RUN go build ./...

FROM alpine:3.9.6 AS server

LABEL maintainer='Rodrigo Collavo <rjcollavo@gmail.com>'

WORKDIR /app
COPY --from=build /app/apio .

CMD ./apio