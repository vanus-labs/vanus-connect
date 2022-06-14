FROM golang:1.17 as builder
WORKDIR /workspace

COPY vance/connectors/sink-elasticsearch  sink
COPY cdk-go cdk-go
WORKDIR /workspace/sink
ARG TARGETOS
ARG TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o bin/sink cmd/main.go

FROM centos:8.4.2105
RUN mkdir -p /vance
COPY --from=builder /workspace/sink/bin/sink /vance/bin/sink
WORKDIR /vance
ENTRYPOINT ["bin/sink"]
