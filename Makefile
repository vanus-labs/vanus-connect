WORKDIR=$(shell pwd)
#os linux or darwin
GOOS ?= linux
#arch amd64 or arm64
GOARCH ?= amd64
GO_BUILD= GO111MODULE=on CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -trimpath
DOCKER_REGISTRY ?= public.ecr.aws
DOCKER_REPO ?= ${DOCKER_REGISTRY}/vanus/connector
IMAGE_TAG ?= latest
DOCKER_BUILD_ARG= --build-arg TARGETARCH=$(GOARCH) --build-arg TARGETOS=$(GOOS)
DOCKER_PLATFORM ?= linux/amd64,linux/arm64

push-java-image:
	docker buildx build -t ${DOCKER_REPO}/${CONNECTOR}:${IMAGE_TAG} -f build/java/Dockerfile \
    		--platform ${DOCKER_PLATFORM} \
    		--build-arg connector=${CONNECTOR} \
    		--push .

# make push-go-image CONNECTOR=source-http
# make push-go-image DOCKER_REGISTRY=linkall.tencentcloudcr.com CONNECTOR=sink-feishu
push-go-image:
	docker buildx build $(DOCKER_BUILD_ARG) -t ${DOCKER_REPO}/${CONNECTOR}:${IMAGE_TAG} -f build/go/Dockerfile \
		--platform ${DOCKER_PLATFORM} \
		--build-arg connector=${CONNECTOR} \
		--push .
