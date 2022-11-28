FROM --platform=$BUILDPLATFORM  maven:3.8.5-openjdk-17 as builder
WORKDIR /workspace
COPY . .
RUN mvn clean package
FROM openjdk:17
WORKDIR /vance
COPY target/sink-http-1.0-SNAPSHOT-jar-with-dependencies.jar /vance
CMD ["java", "-jar", "./sink-http-1.0-SNAPSHOT-jar-with-dependencies.jar"]