FROM openjdk:8-jre-alpine
WORKDIR /vance
COPY target/aws-sns-source-1.0-SNAPSHOT-jar-with-dependencies.jar /vance
CMD ["java", "-jar", "./aws-sns-source-1.0-SNAPSHOT-jar-with-dependencies.jar"]
EXPOSE 8082