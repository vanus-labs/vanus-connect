#!/usr/bin/env bash

DOCKER_REPO=$1
IMAGE_TAG=$2
DOCKER_PLATFORM=$3

connector_list=(
  "sink-display"
  "sink-doris"
  "sink-elasticsearch"
  "sink-email"
  "sink-feishu"
  "sink-k8s"
  "sink-mongodb"
  "sink-slack"
  "sink-tencentcloud-scf"
  "source-alicloud-billing"
  "source-aws-billing"
  "source-http"
  "source-tencentcloud-cos"
)

for((i=0;i<${#connector_list[@]};i++))
do
    docker buildx build -t "${DOCKER_REPO}/${connector_list[i]}:${IMAGE_TAG}" -f build/go/Dockerfile \
        --platform "${DOCKER_PLATFORM}" \
        --build-arg connector="${connector_list[i]}" \
        --push .
done;


