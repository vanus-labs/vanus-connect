---
title: AWS S3
---

# AWS-S3 Source 

## Overview

A [Vance Connector][vc] which retrieves S3 events, transform them into CloudEvents and deliver CloudEvents to the target URL.

## User Guidelines

### Connector Introduction

The AWS-S3 Source is a [Vance Connector][vc] which designed to retrieve S3 events from a specific bucket and 
transform them into CloudEvents based on [CloudEvents Adapter specification][ceas].

This connector allows users to specify a SQS queue to receive S3 event notification messages. 
It will automatically create a SQS queue if you don't specify yours.

The original S3 events looks like:

```json
{  
   "Records":[  
      {  
         "eventVersion":"2.1",
         "eventSource":"aws:s3",
         "awsRegion":"us-west-2",
         "eventTime":"1970-01-01T00:00:00.000Z",
         "eventName":"ObjectCreated:Put",
         "userIdentity":{  
            "principalId":"AIDAJDPLRKLG7UEXAMPLE"
         },
         "requestParameters":{  
            "sourceIPAddress":"127.0.0.1"
         },
         "responseElements":{  
            "x-amz-request-id":"C3D13FE58DE4C810",
            "x-amz-id-2":"FMyUVURIY8/IgAtTv8xRjskZQpcIZ9KG4V5Wp6S7S/JRWeUWerMUE5JgHvANOjpD"
         },
         "s3":{  
            "s3SchemaVersion":"1.0",
            "configurationId":"testConfigRule",
            "bucket":{  
               "name":"mybucket",
               "ownerIdentity":{  
                  "principalId":"A3NL1KOZZKExample"
               },
               "arn":"arn:aws:s3:::mybucket"
            },
            "object":{  
               "key":"HappyFace.jpg",
               "size":1024,
               "eTag":"d41d8cd98f00b204e9800998ecf8427e",
               "versionId":"096fKKXTRTtl3on89fVO.nfljtsv6qko",
               "sequencer":"0055AED6DCD90281E5"
            }
         }
      }
   ]
}
```

A transformed CloudEvent looks like:

``` json
{
  "id" : "C3D13FE58DE4C810.FMyUVURIY8/IgAtTv8xRjskZQpcIZ9KG4V5Wp6S7S/JRWeUWerMUE5JgHvANOjpD",
  "source" : "aws:s3.us-west-2.mybucket",
  "specversion" : "V1",
  "type" : "com.amazonaws.s3.ObjectCreated:Put",
  "datacontenttype" : "application/json",
  "subject" :	"HappyFace.jpg",
  "time" : "1970-01-01T00:00:00.000Z",
  "data" : {
    "s3":{  
            "s3SchemaVersion":"1.0",
            "configurationId":"testConfigRule",
            "bucket":{  
               "name":"mybucket",
               "ownerIdentity":{  
                  "principalId":"A3NL1KOZZKExample"
               },
               "arn":"arn:aws:s3:::mybucket"
            },
            "object":{  
               "key":"HappyFace.jpg",
               "size":1024,
               "eTag":"d41d8cd98f00b204e9800998ecf8427e",
               "versionId":"096fKKXTRTtl3on89fVO.nfljtsv6qko",
               "sequencer":"0055AED6DCD90281E5"
            }
         }
  }
}
```

## AWS-S3 Source Configs

Users can specify their configs by either setting environments variables or mount a config.json to
`/vance/config/config.json` when they run the connector. Find examples of setting configs [here][config].

### Config Fields of the AWS-S3 Source

| Configs   | Description                                                                     | Example                 | Required                 |
|:----------|:--------------------------------------------------------------------------------|:------------------------|:------------------------|
| v_target  | `v_target` is used to specify the target URL HTTP Source will send CloudEvents to. | "http://localhost:8081" |**YES** |
| s3_bucket_arn    | `s3_bucket_arn` is your bucket arn. | "arn:aws:s3:::mybucket"                  |**YES** |
| s3_events    | `s3_events` is an json array consisting of [s3 events](https://docs.aws.amazon.com/AmazonS3/latest/userguide/notification-how-to-event-types-and-destinations.html) you're interested in. | ["s3:ObjectCreated:*","s3:ObjectRemoved:*"]                  |**YES** |
| region    | `region` describes where the SQS queue will be created at. This field is only required when you didn't specify your sqsArn.| "us-west-2"                  |**NO** |
| sqs_arn    | `sqs_arn` is the arn of your SQS queue. The AWS-S3 Source will create a queue located at `region` if this field is omitted.| "arn:aws:sqs:us-west-2:12345678910:myqueue"                  |**NO** |

## AWS-S3 Source Secrets

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
| awsSecretAccessKey    | `awsSecretAccessKey` is the Secret access key of youraws credential. | "BASE64VALUEOFYOURSECRETKEY="                  |**YES** |


## AWS-S3 Source Image

> docker.io/vancehub/source-aws-s3

### Run the S3-source image in a container

Mount your local config file and secret file to specific positions with `-v` flags.

```shell
docker run -v $(pwd)/secret.json:/vance/secret/secret.json -v $(pwd)/config.json:/vance/config/config.json --rm docker.io/vancehub/source-aws-s3
```

## Local Development

You can run the source codes of the AWS-S3 Source locally as well.

### Building via Maven

```shell
$ cd connectors/source-aws-s3
$ mvn clean package
```

### Running via Maven

```shell
$ mvn exec:java -Dexec.mainClass="com.linkall.source.aws.s3.Entrance"
```

⚠️ NOTE: For better local development and test, the connector can also read configs from `main/resources/config.json`. So, you don't need to 
declare any environment variables or mount a config file to `/vance/config/config.json`. Same logic applies to `main/resources/secret.json` as well.

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md
[config]: https://github.com/linkall-labs/vance-docs/blob/main/docs/connector.md
[ceas]: https://github.com/cloudevents/spec/blob/main/cloudevents/adapters/aws-s3.md