---
title: Amazon S3
---

# Amazon S3 Source
This document provides a brief introduction to the Amazon S3 Source. It's also designed to guide you through the
process of running an Amazon S3 Source Connector.

## Introduction
The Amazon S3 Source connector subscribes to event notifications from an Amazon S3 bucket. To enable notifications, you
should list events that you want Amazon S3 to publish in your configuration file. Ensure you also provide the Amazon Resource
Name (ARN) value of an Amazon SQS queue to receive notifications. Our S3 Source will create an SQS queue for you if you didn't
specify one, which requires extra permissions for your AWS user account.

## S3 Event Structure
The [event notification][s3-events] sent from Amazon S3 is in the JSON format.

Here is an example of a message Amazon S3 sends to publish a `s3:ObjectCreated:Put` event.
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
###
The Amazon S3 Source will transform the S3 event above into a CloudEvent with the following structure:
```json
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
The process to convert AWS S3 events into CloudEvents conforms to the [CloudEvents S3 Adapter specification][s3-adapter].

## Features
- **At least once delivery**: Events sent from the S3 Source are designed to be delivered at least once.

## Limitations
- **Limited to one S3 bucket**: Each S3 Source subscribes to event notifications from one Amazon S3 bucket.

## IAM Policy for Amazon S3 Source
- s3:PutBucketNotification
- sqs:GetQueueUrl
- sqs:GetQueueAttributes
- sqs:SetQueueAttributes
- sqs:ListQueues
- sqs:CreateQueue (only needed when you didn't specify your SQS queue)
---
## Quick Start
This quick start will guide you through the process of running an Amazon S3 Source connector.

### Prerequisites
- A container runtime (i.e., docker).
- An Amazon [S3 bucket][s3-bucket].
- A Properly settled [IAM] policy for your AWS user account.
- An AWS account configured with [Access Keys][access-keys].

### Set S3 Source Configurations
You can specify your configs by either setting environments variables or mounting a config.json to
`/vance/config/config.json` when running the connector.

Here is an example of a configuration file for the Amazon S3 Source.
```shell
$ vim config.json
{
  "v_target": "http://host.docker.internal:8081",
  "region": "us-west-2",
  "s3_bucket_arn": "arn:aws:s3:::mybucket",
  "s3_events": ["s3:ObjectCreated:*","s3:ObjectRemoved:*"]
}
```

| Configs       | Description                                                                                                                | Example                                     | Required|
|:--------------|:---------------------------------------------------------------------------------------------------------------------------|:--------------------------------------------|:--------|
| v_target      | `v_target` specifies the target URL the Amazon S3 Source will send CloudEvents to.                                         | "http://host.docker.internal:8081"          | **YES** |
| s3_bucket_arn | `s3_bucket_arn` is your bucket ARN.                                                                                        | "arn:aws:s3:::mybucket"                     | **YES** |
| s3_events     | `s3_events` is an json array consisting of [types of s3 events][s3-events-types] you're interested in.                     | ["s3:ObjectCreated:*","s3:ObjectRemoved:*"] | **YES** |
| region        | `region` describes where the SQS queue will be created at. This field is only required when you didn't specify your sqsArn.| "us-west-2"                                 | **NO**  |
| sqs_arn       | `sqs_arn` is the ARN of your SQS queue. The S3 Source will create a SQS queue located at `region` if `sqs_arn` is omitted. | "arn:aws:sqs:us-west-2:12345678910:myqueue" | **NO**  |

### Set S3 Source Secrets

Users should set their sensitive data Base64 encoded in a secret file. And mount that secret file to
`/vance/secret/secret.json` when running the connector.

#### Encode your Sensitive Data
Replace `MY_SECRET` with your sensitive data to get the Base64-based string.

```shell
$ echo -n MY_SECRET | base64
TVlfU0VDUkVU
```

Here is an example of a secret file for the S3 Source.
```shell
$ vim secret.json
{
  "awsAccessKeyID": "TVlfU0VDUkVU",
  "awsSecretAccessKey": "U2VjcmV0QWNjZXNzS2V5"
}
```
| Secrets            | Description                                                          | Example                       |Required|
|:-------------------|:---------------------------------------------------------------------|:------------------------------|:-------|
| awsAccessKeyID     | `awsAccessKeyID` is the Access key ID of your aws credential.        | "BASE64VALUEOFYOURACCESSKEY=" |**YES** |
| awsSecretAccessKey | `awsSecretAccessKey` is the Secret access key of youraws credential. | "BASE64VALUEOFYOURSECRETKEY=" |**YES** |

### Run the Amazon S3 Source with Docker

Create your `config.json` and `secret.json`, and mount them to specific paths to run the S3 source using following command.

In order to send events to Display Connector, set the `v_target` value as `http://host.docker.internal:8081` in your config.json file.

```shell
docker run -v $(pwd)/secret.json:/vance/secret/secret.json -v $(pwd)/config.json:/vance/config/config.json --rm vancehub/source-aws-s3
```

### Verify the Amazon S3 Source

You can verify if the Amazon S3 Source works properly by running the Display Sink and by uploading a file to the S3 bucket.
> docker run -p 8081:8081 --rm vancehub/sink-display

:::tip
Set the v_target as http://host.docker.internal:8081
:::

Here is an example output of the Display Connector when I upload a "cat.jpg" to the "s3-my-test-bucket".
```shell
[01:41:19:558] [INFO] - com.linkall.sink.display.DisplaySink.lambda$start$0(DisplaySink.java:21) - receive a new event, in total: 1
[01:41:19:577] [INFO] - com.linkall.sink.display.DisplaySink.lambda$start$0(DisplaySink.java:23) - {
  "id" : "E1Y3YATHFSHQHDCC.rq4Gxw46+dNp5X8SQUF4Ckur54WAA415OX2u8+inlqCBKVlRgFeaxRlcFGf/vEo/Uylc2xmQR0",
  "source" : "aws:s3.us-west-2.s3-my-test-bucket",
  "specversion" : "V1",
  "type" : "com.amazonaws.s3.ObjectCreated:Put",
  "datacontenttype" : "application/json",
  "subject" : "cat.jpg",
  "time" : "2022-09-27T01:33:42.334Z",
  "data" : {
    "s3SchemaVersion" : "1.0",
    "configurationId" : "vance-s3-notification",
    "bucket" : {
      "name" : "s3-my-test-bucket",
      "ownerIdentity" : {
        "principalId" : "A189UASLUCKECHAEF"
      },
      "arn" : "arn:aws:s3:::s3-my-test-bucket"
    },
    "object" : {
      "key" : "cat.jpg",
      "size" : 122667,
      "eTag" : "850ba6f83a3a0b216e8a95393630f0ca",
      "sequencer" : "00633252F64A47824B"
    }
  }
}
```

[s3-events]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/notification-content-structure.html
[s3-adapter]: https://github.com/cloudevents/spec/blob/main/cloudevents/adapters/aws-s3.md
[s3-events-types]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/notification-how-to-event-types-and-destinations.html
[s3-bucket]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/UsingBucket.html
[iam]: https://docs.aws.amazon.com/IAM/latest/UserGuide/introduction.html?icmpid=docs_iam_console
[access-keys]: https://docs.aws.amazon.com/general/latest/gr/aws-sec-cred-types.html#access-keys-and-secret-access-keys