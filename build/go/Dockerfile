FROM --platform=$BUILDPLATFORM golang:1.20.5 as builder

ARG connector

COPY . /build/vanus-connect

WORKDIR /build/vanus-connect/connectors/${connector}

RUN go mod tidy

ARG TARGETOS
ARG TARGETARCH
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -o /build/vanus-connect/bin/${connector} ./cmd/main.go

#FROM centos:8.4.2105
FROM ubuntu:22.10
# install ca-certificates
RUN apt update -y && apt upgrade -y && apt install -y ca-certificates && update-ca-certificates

ARG connector

COPY --from=builder /build/vanus-connect/bin/${connector} /vanus-connect/bin/${connector}

ENV CONNECTOR=${connector}
ENV EXECUTABLE_FILE=/vanus-connect/bin/${connector}
ENV CONNECTOR_HOME=/vanus-connect
ENV CONNECTOR_CONFIG=/vanus-connect/config/config.yml
ENV CONNECTOR_SECRET=/vanus-connect/secret/secret.yml
ENV LOG_LEVEL=INFO

RUN echo '#!/bin/sh' >> /vanus-connect/run.sh
RUN echo 'exec $EXECUTABLE_FILE' >> /vanus-connect/run.sh
RUN chmod a+x /vanus-connect/bin/${connector}
RUN chmod a+x /vanus-connect/run.sh
RUN mkdir /vanus-connect/config
RUN touch /vanus-connect/config/config.yml

EXPOSE 8080

ENTRYPOINT ["/vanus-connect/run.sh"]
