# AWS SNS Source

## Overview
AWS SNS(Simple Notification Service) Source is a Vance Source Connector which subscribe to the SNS topic and receive messages published to the topic, and then transform them into CloudEvents and deliver them to the target URL. 
Push is adopted by AWS SNS to deliver messages from SNS topics to the endpoints. Therefore, AWS SNS Source Connector should subscribe the SNS topics and start an endpoint to receive messages from the SNS topics. Application-to-application subscribers supported by SNS include http/https endpoints, Amazon Kinesis Data Firehose, Amazon SQS and AWS Lambda. AWS SNS Source Connector of Vance support http/https endpoints to receive messages. 

## User Guidelines

## Adapter
Adapter convert events into CloudEvents.
AWS SNS may send a subscription confirmation, notification, or unsubscribe confirmation message to your HTTP/HTTPS endpoints.
Each section below describes how to transform SNS messages events into CloudEvents.
Subscription confirmation
### Subscription Confirmation

| CloudEvents Attribute | Value                                           |
| :-------------------- | :---------------------------------------------- |
| `id`                  | "x-amz-sns-message-id" value |
| `source`              | "x-amz-sns-topic-arn" value |
| `specversion`         | `1.0`                                           |
| `type`                | "x-amz-sns-message-type" value    |
| `datacontenttype`     | `application/json`         |
| `dataschema`          | Omit                                            |
| `subject`             | Omit                        |
| `time`                | "Timestamp" value                               |
| `data`                | HTTP payload                                       |

### Notification

| CloudEvents Attribute | Value                                           |
| :-------------------- | :---------------------------------------------- |
| `id`                  | "x-amz-sns-message-id" value |
| `source`              | "x-amz-sns-subscription-arn" value |
| `specversion`         | `1.0`                                           |
| `type`                | "x-amz-sns-message-type" value    |
| `datacontenttype`     | `application/json`         |
| `dataschema`          | Omit                                            |
| `subject`             | "Subject" value (if present)                    |
| `time`                | "Timestamp" value                               |
| `data`                | HTTP payload                                       |

### Unsubscribe Confirmation

| CloudEvents Attribute | Value                                           |
| :-------------------- | :---------------------------------------------- |
| `id`                  | "x-amz-sns-message-id" value |
| `source`              | "x-amz-sns-subscription-arn" value |
| `specversion`         | `1.0`                                           |
| `type`                | "x-amz-sns-message-type" value    |
| `datacontenttype`     | `application/json`         |
| `dataschema`          | Omit                                            |
| `subject`             | Omit                    |
| `time`                | "Timestamp" value                               |
| `data`                | HTTP payload                                       |

Original SNS message pushed to http/https endpoints looks like:
```HTTP
 POST / HTTP/1.1
x-amz-sns-message-type: Notification
x-amz-sns-message-id: 22b80b92-fdea-4c2c-8f9d-bdfb0c7bf324
x-amz-sns-topic-arn: arn:aws:sns:us-west-2:123456789012:MyTopic
x-amz-sns-subscription-arn: arn:aws:sns:us-west-2:123456789012:MyTopic:c9135db0-26c4-47ec-8998-413945fb5a96
Content-Length: 773
Content-Type: text/plain; charset=UTF-8
Host: myhost.example.com
Connection: Keep-Alive
User-Agent: Amazon Simple Notification Service Agent

{
  "Type" : "Notification",
  "MessageId" : "22b80b92-fdea-4c2c-8f9d-bdfb0c7bf324",
  "TopicArn" : "arn:aws:sns:us-west-2:123456789012:MyTopic",
  "Subject" : "My First Message",
  "Message" : "Hello world!",
  "Timestamp" : "2012-05-02T00:54:06.655Z",
  "SignatureVersion" : "1",
  "Signature" : "EXAMPLEw6JRN...",
  "SigningCertURL" : "https://sns.us-west-2.amazonaws.com/SimpleNotificationService-f3ecfb7224c7233fe7bb5f59f96de52f.pem",
  "UnsubscribeURL" : "https://sns.us-west-2.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-west-2:123456789012:MyTopic:c9135db0-26c4-47ec-8998-413945fb5a96"
}
```
A transformed CloudEvent looks like:
```JSON
CloudEvent:{
  "id" : "86f1c45f-4b82-5337-8bd6-9f289c2690da",
  "source" : "arn:aws:sns:us-west-2:843378899134:TestTopic:70fdfbd7-805a-4dbd-83a2-66cffa859a27",
  "specversion" : "V1",
  "type" : "com.amazonaws.sns.Notification",
  "datacontenttype" : "application/json",
  "subject" : "display",
  "time" : "2022-08-18T06:00:04.638Z",
  "data" : {
    "Type" : "Notification",
    "MessageId" : "86f1c45f-4b82-5337-8bd6-9f289c2690da",
    "TopicArn" : "arn:aws:sns:us-west-2:843378899134:TestTopic",
    "Subject" : "display",
    "Message" : "a message send to display",
    "Timestamp" : "2022-08-18T06:00:04.638Z",
    "SignatureVersion" : "1",
    "Signature" : "AVzjPMTtXaRQz1XA/8kA2ouc4BVKNQYPHzAK5nF3ehUz44YydQmNv2VzQZ6in+riHtv6PgdWz9QsnYjDSItO9rVL+k6ngQBIagHCOkf+hL307J8fsIIRMZsYIHsv7+JNmhjQfsM3BelDJh+OBj6gJPmTKz+PtOXiPAxDKHijmOqN4Eq8TIjTkHX4YrcDAS9ZmM77mtvRlbpVVFNaYTtfYX8sYktuEiFZSfrBInWY74nP47Aqb+s7vi0soRDS6ACCd64k1U4X3r5VDs/t3NLrU++OCfhuE1GvJe+2WVxCvZbkY6K/nL/Q96mZMYT7iNRSma0mDQ89PQOclgbM58FVog==",
    "SigningCertURL" : "https://sns.us-west-2.amazonaws.com/SimpleNotificationService-56e67fcb41f6fec09b0196692625d385.pem",
    "UnsubscribeURL" : "https://sns.us-west-2.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-west-2:843378899134:TestTopic:70fdfbd7-805a-4dbd-83a2-66cffa859a27"
  }
}
```
## AWS SNS Source Connector Configs
Users can specify their configs by either setting environments variables or mount a config.json to /vance/config/config.json when they run the connector. Find examples of setting configs [here](https://github.com/linkall-labs/vance-docs/blob/main/docs/connector.md). 
AWS SNS Source Connector requires users to provide SNS TopicArn, target url, host url, port listened of endpoints and the protocol(chosen in http and https).
### Config Fields of the GitHub Source
|  Configs    |  Description    																  |  Example    			  |  Required    |
|  :----:     |  :----:         																  |  :----:     			  |  :----:      |
|  v_target   |  v_target is used to specify the target URL HTTP Source will send CloudEvents to  |  "http://localhost:8081"  |  YES  		 |
|  v_port     |  the port of http/http endpoints to receive SNS messages					  |  "8080"	                  |  YES         |
|  v_host     |  the url of http/https endpoints				  |  "http://xxx.xxx.xxx:8082"	                  |  YES         |
|  TopicArn     |  the arn of the SNS topic					  |  "arn:aws:sns:us-west-2:843378899134:Testxxxx"	                  |  YES         |
|  protocol     |  the protocol used to subscribe SNS topic					  |  "http"	                  |  YES         |
The config.json looks like:
```JSON
{
  "TopicArn": "arn:aws:sns:us-west-2:843378899134:Testxxx",
  "v_target": "http://8.142.xxx.xx:8080",
  "v_port": "8082",
  "v_host": "http://8.142.xxx.xx:8082",
  "protocol": "http"
}
```
## AWS SNS Source Connector Secrets
Users should set their sensitive data Base64 encoded in a secret file. And mount your local secret file to /vance/secret/secret.json when you run the connector.
### Encode your sensitive data
```Bash
$ echo -n ABCDEFG | base64
QUJDREVGRw==
```
Replace 'ABCDEFG' with your sensitive data.
### Set your local secret file
```Bash
$ cat secret.json
{
   "awsAccessKeyID": "${awsAccessKeyID}",
   "awsSecretAccessKey": "${awsSecretAccessKey}"
}
```
| Secrets   | Description                                                                     | Example                 | Required                 |
|:----------|:--------------------------------------------------------------------------------|:------------------------|:------------------------|
| awsAccessKeyID  | `awsAccessKeyID` is the Access key ID of your aws credential. | "BASE64VALUEOFYOURACCESSKEY=" |**YES** |
| awsSecretAccessKey    | `awsSecretAccessKey` is the Secret access key of your aws credential. | "BASE64VALUEOFYOURSECRETKEY="                  |**YES** |
## AWS SNS Source Connector Image
>    
### Run the SNS-source image in a container
Touch the config.json and secret.json in a directory, and mount your local config file and secret file to specific positions with -v flags.
```Bash
docker run -v $(pwd)/secret.json:/vance/secret/secret.json -v $(pwd)/config.json:/vance/config/config.json -p 8082:8082 aws-sns-source
```
You can start a CloudEvents Display Sink Connector of Vance, pull the image from:
```Bash
docker.io/vancehub/display
```
Run the image with -p flags to configure the explosure port, and you can see CloudEvents delivered from AWS SNS Source Connector.

## Local Development

You can run the source codes of the AWS-SQS Source locally as well.

### Building via Maven

```shell
$ cd connectors/source-aws-sqs
$ mvn clean package
```

### Running via Maven

```shell
$ mvn exec:java -Dexec.mainClass="com.linkall.source.aws.sqs.Entrance"
```

⚠️ NOTE: For better local development and test, the connector can also read configs from `main/resources/config.json`. So, you don't need to 
declare any environment variables or mount a config file to `/vance/config/config.json`. Same logic applies to `main/resources/secret.json` as well.