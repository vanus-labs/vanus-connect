#!/usr/bin/env bash
sex -ex

DOCKER_REPO=$1
IMAGE_TAG=$2
DOCKER_PLATFORM=$3

connector_list=(
    "source-mongodb"
    "source-mysql"
)

for((i=0;i<${#connector_list[@]};i++))
do
    docker buildx build -t "${DOCKER_REPO}/${connector_list[i]}:${IMAGE_TAG}" -f build/java/Dockerfile \
        --platform "${DOCKER_PLATFORM}" \
        --build-arg connector="${connector_list[i]}" \
        --push .
done;

