FROM --platform=$BUILDPLATFORM golang:1.17-alpine AS builder
WORKDIR /workspace
ARG CONNECTOR

COPY vance/connectors/$CONNECTOR vance/connectors/$CONNECTOR
COPY cdk-go cdk-go
WORKDIR /workspace/vance/connectors/$CONNECTOR
RUN go mod download

ARG TARGETOS
ARG TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /workspace/vance/bin/$CONNECTOR cmd/main.go

FROM alpine:3.15.4
WORKDIR /vance
ARG CONNECTOR
COPY --from=builder /workspace/vance/bin/$CONNECTOR /vance/bin/connector

ENV CONNECTOR_CONFIG=/vance/config/config.yaml

ENTRYPOINT ["/vance/bin/connector"]
