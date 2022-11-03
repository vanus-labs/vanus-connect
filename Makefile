
build-dev-java-image:
	docker build -f build/java/Dockerfile \
		--platform linux/amd64 \
		--build-arg connector=${CONNECTOR} \
		--build-arg version=dev  .

build-dev-go-image:
	docker build -f build/go/Dockerfile \
		--platform linux/amd64 \
		--build-arg connector=${CONNECTOR} \
		--build-arg version=dev  .