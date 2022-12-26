---
title: AWS SQS
---

# AWS-SQS Source 

## Overview

A [Vance Connector][vc] which retrieves SQS messages, transform them into CloudEvents and deliver CloudEvents to the target URL.

## User Guidelines

### Connector Introduction

If the original SQS message looks like:

```json
{
  "MessageId": "035e183b-275a-44de-95df-f212be1ed4ea",
  "Body": "Hello World"
}
```

A transformed CloudEvent looks like:

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

## AWS-SQS Source Configs

Users can specify their configs by either setting environments variables or mount a config.json to
`/vance/config/config.json` when they run the connector. Find examples of setting configs [here][config].

### Config Fields of the AWS-SQS Source

| Configs   | Description                                                                     | Example                 | Required                 |
|:----------|:--------------------------------------------------------------------------------|:------------------------|:------------------------|
| v_target  | `v_target` is used to specify the target URL HTTP Source will send CloudEvents to. | "http://localhost:8081" |**YES** |
| sqs_arn    | `sqs_arn` is the arn of your SQS queue.  | "arn:aws:sqs:us-west-2:12345678910:myqueue"                   |**YES** |

## AWS-SQS Source Secrets

Users should set their sensitive data Base64 encoded in a secret file. And mount your local secret file to `/vance/secret/secret.json` when you run the connector.

### Encode your sensitive data

```shell
$ echo -n ABCDEFG | base64
QUJDREVGRw==
```

Replace 'ABCDEFG' with your sensitive data.

### Set your local secret file

```shell
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


## AWS-SQS Source Image

> docker.io/vancehub/source-aws-sqs

### Run the SQS-source image in a container

Mount your local config file and secret file to specific positions with `-v` flags.

```shell
docker run -v $(pwd)/secret.json:/vance/secret/secret.json -v $(pwd)/config.json:/vance/config/config.json --rm docker.io/vancehub/source-aws-sqs
```

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

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md
[config]: https://github.com/linkall-labs/vance-docs/blob/main/docs/connector.md