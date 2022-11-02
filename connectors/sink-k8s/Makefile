WORKDIR=$(shell pwd)
#os linux or darwin
GOOS ?= linux
#arch amd64 or arm64
GOARCH ?= amd64
GO_BUILD= GO111MODULE=on CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -trimpath
DOCKER_REGISTRY ?= docker.io
DOCKER_REPO ?= ${DOCKER_REGISTRY}/vancehub
IMAGE_TAG ?= latest
DOCKER_BUILD_ARG= --build-arg TARGETARCH=$(GOARCH) --build-arg TARGETOS=$(GOOS)
DOCKER_PLATFORM ?= linux/amd64,linux/arm64

clean :
	rm -rf bin

build:
	$(GO_BUILD) -o bin/sink cmd/main.go

docker-push:
	docker buildx build --platform ${DOCKER_PLATFORM} -t ${DOCKER_REPO}/sink-k8s:${IMAGE_TAG} -f Dockerfile . --push

docker-build:
	docker build -t ${DOCKER_REPO}/sink-k8s:${IMAGE_TAG} $(DOCKER_BUILD_ARG) -f Dockerfile .