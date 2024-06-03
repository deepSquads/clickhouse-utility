FROM golang:1.22-alpine3.19 as builder

WORKDIR /app
ARG COMMAND_NAME
ARG COMMAND_PATH
ARG TAGS

COPY . /app/

RUN cat /app/.gitignore

RUN apk add --update --no-cache git

RUN ./.build/build.sh $COMMAND_PATH $COMMAND_NAME $TAGS

FROM alpine:3.19

ARG COMMAND_NAME
ARG COMMAND_PATH
ARG TAGS

WORKDIR /app

RUN apk add --update --no-cache ca-certificates

COPY --from=builder /app/.build/target/$COMMAND_NAME $COMMAND_NAME

ENV COMMAND_NAME $COMMAND_NAME
