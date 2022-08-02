FROM openjdk:8-jre-alpine
WORKDIR /vance
COPY target/aws-sqs-source-1.0-SNAPSHOT-jar-with-dependencies.jar /vance
CMD ["java", "-jar", "./aws-sqs-source-1.0-SNAPSHOT-jar-with-dependencies.jar"]