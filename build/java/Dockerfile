FROM --platform=$BUILDPLATFORM maven:3.8.5-openjdk-17 as builder

ARG TARGETOS
ARG TARGETARCH
ARG connector

COPY . /build/vanus-connect

WORKDIR /build/vanus-connect/connectors/${connector}

RUN microdnf install libatomic
RUN mvn clean package -Dmaven.javadoc.skip=true
# DON'T REMOVE SPACE BEHIND OF '/build/executable.jar'
RUN ls -alh target/*jar-with-dependencies.jar | awk '{system("cp " $9 " /build/executable.jar") }'

FROM openjdk:17

ARG connector

COPY --from=builder /build/executable.jar /vanus-connect/bin/${connector}.jar

ENV CONNECTOR=${connector}
ENV EXECUTABLE_FILE=/vanus-connect/bin/${connector}.jar
ENV CONNECTOR_HOME=/vanus-connect
ENV CONNECTOR_CONFIG=/vanus-connect/config/config.yml
ENV CONNECTOR_SECRET=/vanus-connect/secret/secret.yml
ENV CONNECTOR_SECRET_ENABLE=false

RUN echo '#!/bin/sh' >> /vanus-connect/run.sh
RUN echo 'java -jar $EXECUTABLE_FILE' >> /vanus-connect/run.sh
RUN chmod a+x /vanus-connect/run.sh
RUN mkdir /vanus-connect/config
RUN touch /vanus-connect/config/config.yml

EXPOSE 8080

ENTRYPOINT ["/vanus-connect/run.sh"]
