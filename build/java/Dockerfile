FROM maven:3.8.5-openjdk-17 as builder

ARG TARGETOS
ARG TARGETARCH
ARG connector

COPY . /build/vance

WORKDIR /build/vance/connectors/${connector}

RUN microdnf install libatomic
RUN mvn clean package -Dmaven.javadoc.skip=true
# DON'T REMOVE SPACE BEHIND OF '/build/executable.jar'
RUN ls -alh target/*jar-with-dependencies.jar | awk '{system("cp " $9 " /build/executable.jar") }'

FROM openjdk:17

ARG connector

COPY --from=builder /build/executable.jar /vance/bin/${connector}.jar

ENV CONNECTOR=${connector}
ENV EXECUTABLE_FILE=/vance/bin/${connector}.jar
ENV CONNECTOR_HOME=/vance
ENV CONNECTOR_CONFIG=/vance/config/config.yml
ENV CONNECTOR_SECRET=/vance/secret/secret.yml
ENV CONNECTOR_SECRET_ENABLE=false

RUN echo '#!/bin/sh' >> /vance/run.sh
RUN echo 'java -jar $EXECUTABLE_FILE' >> /vance/run.sh
RUN chmod a+x /vance/run.sh

EXPOSE 8080

ENTRYPOINT ["/vance/run.sh"]
