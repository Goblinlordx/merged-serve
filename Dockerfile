# syntax=docker/dockerfile:1
FROM golang:1.17-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /merged-serve



###

FROM alpine:latest

COPY --from=build /merged-serve /merged-serve

ARG USER=nonroot

RUN adduser -D nonroot \
  && chmod 755 /merged-serve

USER $USER

EXPOSE 8080

ENTRYPOINT [ "./merged-serve" ]
