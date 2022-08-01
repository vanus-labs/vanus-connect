FROM openjdk:8-jre-alpine
WORKDIR /vance
COPY target/source-github-1.0-SNAPSHOT-jar-with-dependencies.jar /vance
CMD ["java", "-jar", "./source-github-1.0-SNAPSHOT-jar-with-dependencies.jar"]
EXPOSE 8081