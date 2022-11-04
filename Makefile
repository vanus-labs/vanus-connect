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

build-dev-java-image:
	docker build -f build/java/Dockerfile \
		--platform linux/amd64 \
		--build-arg connector=${CONNECTOR} \
		--build-arg version=dev  .

build-go-image:
	docker build $(DOCKER_BUILD_ARG) -t ${DOCKER_REPO}/${CONNECTOR}:${IMAGE_TAG} -f build/go/Dockerfile \
		--build-arg connector=${CONNECTOR} .

push-go-image:
	docker buildx build $(DOCKER_BUILD_ARG) -t ${DOCKER_REPO}/${CONNECTOR}:${IMAGE_TAG} -f build/go/Dockerfile \
		--platform ${DOCKER_PLATFORM} \
		--build-arg connector=${CONNECTOR} \
		--push .
