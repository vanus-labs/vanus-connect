FROM golang:1.18 as compiler

ARG connector

COPY ./vance /tmp/vance
COPY ./cdk-go /tmp/cdk-go

RUN cd /tmp/vance/connectors/${connector} && \
    go build -v -o /tmp/vance/bin/${connector} ./cmd/main.go

FROM centos:8.4.2105

ARG connector

COPY --from=compiler /tmp/vance/bin/${connector} /etc/vance/bin/${connector}
COPY --from=compiler /tmp/vance/connectors/${connector}/run.sh /etc/vance/run.sh

RUN chmod a+x /etc/vance/bin/${connector}
RUN chmod a+x /etc/vance/run.sh

ENV EXECUTABLE_FILE=/etc/vance/bin/${connector}
ENV CONNECTOR_HOME=/etc/vance/${connector}
ENV CONNECTOR_CONFIG=/etc/vance/${connector}/config.yml
ENV CONNECTOR_SECRET=/etc/vance/${connector}/secert.yml

EXPOSE 8080

ENTRYPOINT ["/etc/vance/run.sh"]