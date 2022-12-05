FROM --platform=$BUILDPLATFORM golang:1.17-alpine AS builder
WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG TARGETOS
ARG TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o bin/connector cmd/main.go

FROM alpine:3.15.4
WORKDIR /vance
COPY --from=builder /workspace/bin/connector /vance/bin/connector

ENV CONNECTOR_CONFIG=/vance/config/config.yaml
ENV CONNECTOR_SECRET=/vance/secret/secret.yaml

ENTRYPOINT ["/vance/bin/connector"]
