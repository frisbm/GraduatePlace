FROM golang:1.19.1-alpine

RUN apk update && apk add --no-cache \
    bash \
    make \
    openssh-client \
    git \
    curl \
    gcc

WORKDIR /app

RUN curl -fLo install.sh https://raw.githubusercontent.com/cosmtrek/air/master/install.sh \
    && chmod +x install.sh && sh install.sh && cp ./bin/air /bin/air

CMD air
