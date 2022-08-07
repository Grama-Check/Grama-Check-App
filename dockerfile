# syntax=docker/dockerfile:1

# Build stage

FROM golang:1.18.3-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN go build -o /grama-check

# Run Stage

FROM alpine:3.16

WORKDIR /app
COPY --from=builder grama-check .
COPY  public.pem . 
COPY  private.pem . 
COPY  app.env . 
COPY /public ./public

EXPOSE 9090

CMD [ "./grama-check" ]