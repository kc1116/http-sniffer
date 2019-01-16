OUT := http-sniffer
PKG := ./cmd
DOCKERFILE := ./build/package/Dockerfile

all: run

docker-image:
	docker build -f ${DOCKERFILE} -t ${OUT} .
	
deps:
	go get ${PKG}

build:
	go build -i -v -o ./build/bin/${OUT} ${PKG}

vet:
	@go vet ${PKG_LIST}

clean:
	-@rm ${OUT} ${OUT}-v*

.PHONY: run build docker-image vet deps