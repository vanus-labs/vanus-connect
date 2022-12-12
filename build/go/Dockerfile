FROM golang:1.18.5 as builder
WORKDIR /workspace
ARG TARGETOS
ARG TARGETARCH
ARG connector
COPY . /build/vance
RUN cd /build/vance/connectors/${connector} && \
    go mod tidy && \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -o /build/vance/bin/${connector} ./cmd/main.go

FROM centos:8.4.2105

ARG connector

WORKDIR /vance

COPY --from=builder /build/vance/bin/${connector} /vance/bin/${connector}

ENV CONNECTOR=${connector}
ENV EXECUTABLE_FILE=/vance/bin/${connector}
ENV CONNECTOR_HOME=/vance
ENV CONNECTOR_CONFIG=/vance/config/config.yml
ENV CONNECTOR_SECRET=/vance/secret/secret.yml
ENV LOG_LEVEL=INFO

RUN touch /vance/run.sh && \
    echo '#!/bin/sh' >> /vance/run.sh && \
    echo $EXECUTABLE_FILE >> /vance/run.sh && \
    chmod a+x /vance/bin/${connector} && \
    chmod a+x /vance/run.sh && \
    touch /vance/config/config.yml

EXPOSE 8080

ENTRYPOINT ["/vance/run.sh"]
