FROM golang:1.17 as builder
WORKDIR /workspace

COPY vance-connectors/source-aws-billing  source
COPY cdk-go cdk-go
WORKDIR /workspace/source
ARG TARGETOS
ARG TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o bin/app cmd/main.go

FROM centos:8.4.2105
RUN mkdir -p /server
COPY --from=builder /workspace/source/bin/app /server/bin/app
WORKDIR /server
ENTRYPOINT ["bin/app"]

