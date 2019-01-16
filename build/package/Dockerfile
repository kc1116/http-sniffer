FROM golang:1.10 AS builder

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && chmod +x /usr/local/bin/dep

RUN mkdir -p /go/src/http-sniffer
WORKDIR /go/src/http-sniffer

COPY . .

RUN dep ensure && \
    make build
    
CMD ./build/bin/http-sniffer