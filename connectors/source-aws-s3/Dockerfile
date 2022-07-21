FROM openjdk:8-jre-alpine
WORKDIR /vance
COPY target/aws-s3-source-1.0-SNAPSHOT-jar-with-dependencies.jar /vance
CMD ["java", "-jar", "./aws-s3-source-1.0-SNAPSHOT-jar-with-dependencies.jar"]