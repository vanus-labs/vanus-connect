---
title: Amazon S3
---

# Amazon S3 Sink
This document provides a brief introduction
to the Amazon S3 Sink. It's also designed to guide you through the
process of running an Amazon S3 Sink Connector.

## Introduction
The Amazon S3 (Simple Storage Service) Sink Connector help you **convert received CloudEvents into JSON format and upload**
them to AWS S3. To enable the Amazon S3 service, you should configure your **AWS region and bucket name**. When CloudEvents
arrive, S3 Sink will **cache** them in local files. Then S3 Sink will upload local files according to the **configured upload
strategy**. Files uploaded to AWS S3 will be **partitioned by time**. S3 Sink supports **`HOURLY` and `DAILY`** partitioning.

## Features
- **Upload scheduled interval**: S3 Sink supports **scheduled periodic check** for closing and uploading files to S3. When
  the **time interval** reaches the threshold, the file will be directly closed and uploaded, regardless of whether the file
  is full. The interval defaults to **60 seconds**.
- **Flush size**: When file reaches **flush size**, it will be uploaded automatically. The flush size defaults to **1000**.
- **Partition**: CloudEvents uploaded by S3 Sink are stored in S3 **partitioned**. S3 Sink create storage path in S3 according
  to current time. Time-based partitioning options are daily or hourly.

## Limitations
- **valid schema**: Data S3 Sink received must conform **[CloudEvents Schema Registry][ce-schema]**.
- **Limitation of rate for receiving CloudEvents**: S3 Sink Connector will check the flush size and upload scheduled interval
  one per second. Therefore, the number of CloudEvents sent to S3 Sink per second should be less than flush size.

## IAM policy for Amazon S3 Sink
- s3:List:ListAllMyBuckets
- s3:List:ListBucket
- s3:List:ListMultipartUploadParts
- s3:List:ListBucketMultipartUploads
- s3:Read:GetBucketLocation
- s3:Read:GetObject
- s3:Write:AbortMultipartUpload
- s3:Write:CreateBucket
- s3:Write:DeleteObject
- s3:Write:PutBucketOwnershipControls
- s3:Write:PutObject
 ---
## Quick Start
This quick start will guide you through the process of running an Amazon S3 Sink connector.

### Prerequisites
- A container runtime (i.e., docker).
- An Amazon [S3 bucket][s3-bucket].
- A Properly settled [IAM] policy for your AWS user account.
- An AWS account configured with [Access Keys][access-keys].

### Set S3 Sink Configurations
You can specify your configs by either setting environments
variables or mounting a config.json to `/vance/config/config.json`
when running the connector.

Here is an example of a configuration file for the Amazon S3 Sink.\
  ```json 
  $ vim config.json
  {
    "v_port": "8080",
    "region": "us-west-2",
    "bucketName": "${bucketName}",
    "flushSize": "1000",
    "scheduledInterval": "60",
    "timeInterval": "HOURLY"
  }
  ```

|  Configs    |                                                                                                                                                     Description    																                                                                                                                                                     |  Example    			  |  Required    |
  |  :----:     |:---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------:|  :----:     			  |  :----:      |
|  v_port   |                                                                                                                                           `v_port` Specifies the port S3 Sink will listen to                                                                                                                                            |  "8080"  |  YES  		 |
|  region     |                                                                                                                                    `region` Describes in which aws region S3 bucket is created. 				                                                                                                                                    |  "us-west-2"	                  |  YES         |
|  bucketName     |                                                                                                                                                     Unique name of S3 bucket.					                                                                                                                                                      |  "aws-s3-notify"	                  |  YES         |
|  flushSize     |                                                                                                              `flushSize` Refers to the number of CloudEvents cached to the local file before S3 Sink committing the file.                                                                                                               |  "1000"	     |  NO, defaults to 1000         |
|  scheduledInterval     |                                                                                                        `scheduledInterval` Refers to the maximum time interval between S3 Sink closing and uploading files which unit is second.                                                                                                        |  "60"	     |  NO, defaults to 60         |
|  timeInterval     | `timeInterval` Refers to the partitioning interval of files have been uploaded to the S3. S3 Sink supports `HOURLY` and `DAILY` time interval. For example, when `timeInterval` is `HOURLY`, files uploaded between 3 pm and 4 pm will be partitioned to one path, while files uploaded after 4 pm will be partitioned to another path. |  "HOURLY"	     |  YES        |

### Set S3 Sink Secrets
Users should set their sensitive data Base64 encoded in a secret file.
And mount your local secret file to `/vance/secret/secret.json` when you run the connector.

#### Encode your Sensitive Data
Replace `MY_SECRET` with your sensitive data to get the Base64-based string.

  ```shell
  $ echo -n MY_SECRET | base64
  TVlfU0VDUkVU
  ```

Create a local secret file with the Base64-based string you just got from your secrets.
  ```json
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

### Run the Amazon S3 Sink with Docker
Create your config.json and secret.json, and mount them to
specific paths to run the MySQL Sink using the following command.

>  docker run -v $(pwd)/secret.json:/vance/secret/secret.json -v $(pwd)/config.json:/vance/config/config.json -p 8080:8080 --rm vancehub/sink-aws-s3


### Verify the Amazon S3 Sink

To verify the S3 Sink, you should send CloudEvents to the address of S3 Sink. Try to use the  following `curl` command.

  ```shell 
  > curl -X POST 
  > -d '{"specversion":"0.3","id":"b25e2717-a470-45a0-8231-985a99aa9416","type":"com.github.pull.create","source":"https://github.com/cloudevents/spec/pull/123","time":"2019-07-04T17:31:00.000Z","datacontenttype":"application/json","data":{"much":"wow"}}' 
  > -H'Content-Type:application/cloudevents+json' 
  > http://localhost:8081 
  ``` 

:::tip
Note that the last line contains the address and port to send the event.
:::


After S3 Sink Connector, you can see following logs:
  ```shell 
  [root@iZ8vbhcwtixrzhsn023sghZ s3sink]# docker run -v $(pwd)/secret.json:/vance/secret/secret.json -v $(pwd)/config.json:/vance/config/config.json -p 8080:8080 --rm sink-aws-s3
  [09:02:19:001] [INFO] - com.linkall.vance.common.config.ConfigPrinter.printVanceConf(ConfigPrinter.java:10) - vance configs-v_target: v_target
  [09:02:19:010] [INFO] - com.linkall.vance.common.config.ConfigPrinter.printVanceConf(ConfigPrinter.java:11) - vance configs-v_port: 8080
  [09:02:19:010] [INFO] - com.linkall.vance.common.config.ConfigPrinter.printVanceConf(ConfigPrinter.java:12) - vance configs-v_config_path: /vance/config/config.json
  [09:02:19:011] [INFO] - com.linkall.sink.aws.AwsHelper.checkCredentials(AwsHelper.java:14) - ====== Check aws Credential start ======
  [09:02:19:018] [INFO] - com.linkall.sink.aws.AwsHelper.checkCredentials(AwsHelper.java:17) - ====== Check aws Credential end ======
  [09:02:20:583] [INFO] - com.linkall.vance.core.http.HttpServerImpl.lambda$listen$5(HttpServerImpl.java:127) - HttpServer is listening on port: 8080
  [02:39:20:943] [INFO] - com.linkall.sink.aws.S3Sink.lambda$start$0(S3Sink.java:96) - receive a new event, in total: 1
  [02:40:14:141] [INFO] - com.linkall.sink.aws.S3Sink.uploadFile(S3Sink.java:162) - [upload file <eventing-0000001> completed
  ```

[ce-schema]: https://github.com/cloudevents/spec/blob/main/schemaregistry/spec.md