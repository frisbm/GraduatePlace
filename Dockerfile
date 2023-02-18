FROM golang:1.19.1-alpine

RUN apk update && apk add --no-cache \
    bash \
    git \
    gcc \
    curl \
    musl-dev

WORKDIR /app

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b /usr/local/bin

CMD air
