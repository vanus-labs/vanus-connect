---
title: Amazon SQS
---

# AWS-SQS Source
This document provides a brief introduction of the SQS Source.
It is also designed to guide you through the process of running an
SQS Source Connector.

## Introduction
A [Vance Connector][vc] which retrieves SQS messages, transform them into CloudEvents
and deliver CloudEvents to the target URL.

## SQS Event Structure

For example, if the incoming message looks like:
```json
{
  "MessageId": "035e183b-275a-44de-95df-f212be1ed4ea",
  "Body": "Hello World"
}
```
###
The Amazon SQS Source will transform the SQS message above into a CloudEvent
with the following structure:
``` json
{
  "id" : "035e183b-275a-44de-95df-f212be1ed4ea",
  "source" : "cloud.aws.sqs.us-west-2.my-test-queue",
  "specversion" : "V1",
  "type" : "com.amazonaws.sqs.message",
  "datacontenttype" : "text/plain",
  "time" : "2022-08-02T11:01:13.828+08:00",
  "data" : "Hello World"
}
```

---
## Quick Start
This quick start will guide you through the process of running an SQS Source Connector.

### Set SQS Source Configurations
You can specify your configs by either setting environments
variables or mounting a config.json to `/vance/config/config.json`
when running the Connector.

Here is an example of a configuration file for the SQS Source.
```json
config.json
{
  "v_target": "http://localhost:8081",
  "sqs_arn": "arn:aws:sqs:us-west-2:12345678910:myqueue"
}
```

| Configs   | Description                                                                     | Example                 | Required                 |
|:----------|:--------------------------------------------------------------------------------|:------------------------|:------------------------|
| v_target  | `v_target` is used to specify the target URL HTTP Source will send CloudEvents to. | "http://localhost:8081" |**YES** |
| sqs_arn    | `sqs_arn` is the arn of your SQS queue.  | "arn:aws:sqs:us-west-2:12345678910:myqueue"                   |**YES** |

### AWS-SQS Source Secrets
Users should set their sensitive data Base64 encoded in a secret file.
And mount your local secret file to `/vance/secret/secret.json` when you run the Connector.

#### Encode your sensitive data
Replace MY_SECRET with your sensitive data to get the Base64-based string.

```shell
$ echo -n MY_SECRET | base64
QUJDREVGRw==
```
Here is an example of a Secret file for the SQS Source.
```shell
$ cat secret.json
{
  "awsAccessKeyID": "TVlfU0VDUkVUTVlfU0VDUkVU",
  "awsSecretAccessKey": "TVlfU0VDUkVUTVlfU0VDUkVU"
}
```
#### Secret Fields of the SQS Source

| Secrets   | Description                                                                     | Example                 | Required                 |
|:----------|:--------------------------------------------------------------------------------|:------------------------|:------------------------|
| awsAccessKeyID  | `awsAccessKeyID` is the Access key ID of your aws credential. | "BASE64VALUEOFYOURACCESSKEY=" |**YES** |
| awsSecretAccessKey    | `awsSecretAccessKey` is the Secret access key of your aws credential. | "BASE64VALUEOFYOURSECRETKEY="                  |**YES** |

### Run the SQS Source with Docker
Create your config.json and secret.json, and mount them to
specific paths to run the SQS Source using the following command.

> docker run -v $(pwd)/secret.json:/vance/secret/secret.json -v $(pwd)/config.json:/vance/config/config.json --rm vancehub/soure-aws-sqs
docker pull




[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md
[config]: https://github.com/linkall-labs/vance-docs/blob/main/docs/connector.mdub.com/linkall-labs/vance-docs/blob/main/docs/connector.md