FROM golang:1.10 AS builder

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && chmod +x /usr/local/bin/dep

RUN mkdir -p /go/src/http-sniffer
WORKDIR /go/src/http-sniffer

COPY . .

RUN dep ensure && \
    make build

RUN cp ./build/bin/http-sniffer /usr/local/bin/http-sniffer
CMD http-sniffer