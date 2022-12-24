FROM golang:1.18.5 as builder

ARG connector

COPY . /build/vance

WORKDIR /build/vance/connectors/${connector}

RUN go mod tidy
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -o /build/vance/bin/${connector} ./cmd/main.go

FROM centos:8.4.2105

ARG connector

COPY --from=builder /build/vance/bin/${connector} /vance/bin/${connector}

ENV CONNECTOR=${connector}
ENV EXECUTABLE_FILE=/vance/bin/${connector}
ENV CONNECTOR_HOME=/vance
ENV CONNECTOR_CONFIG=/vance/config/config.yml
ENV CONNECTOR_SECRET=/vance/secret/secret.yml
ENV LOG_LEVEL=INFO

RUN echo '#!/bin/sh' >> /vance/run.sh
RUN echo $EXECUTABLE_FILE >> /vance/run.sh
RUN chmod a+x /vance/bin/${connector}
RUN chmod a+x /vance/run.sh

EXPOSE 8080

ENTRYPOINT ["/vance/run.sh"]
