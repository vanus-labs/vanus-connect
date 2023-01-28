FROM golang:1.18.5 as builder

ARG connector

COPY . /build/vanus-connector

WORKDIR /build/vanus-connector/connectors/${connector}

RUN go mod tidy
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -o /build/vanus-connector/bin/${connector} ./cmd/main.go

FROM centos:8.4.2105

ARG connector

COPY --from=builder /build/vanus-connector/bin/${connector} /vanus-connector/bin/${connector}

ENV CONNECTOR=${connector}
ENV EXECUTABLE_FILE=/vanus-connector/bin/${connector}
ENV CONNECTOR_HOME=/vanus-connector
ENV CONNECTOR_CONFIG=/vanus-connector/config/config.yml
ENV CONNECTOR_SECRET=/vanus-connector/secret/secret.yml
ENV LOG_LEVEL=INFO

RUN echo '#!/bin/sh' >> /vanus-connector/run.sh
RUN echo $EXECUTABLE_FILE >> /vanus-connector/run.sh
RUN chmod a+x /vanus-connector/bin/${connector}
RUN chmod a+x /vanus-connector/run.sh

EXPOSE 8080

ENTRYPOINT ["/vanus-connector/run.sh"]
