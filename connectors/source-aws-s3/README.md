---
title: Amazon S3
---

# Amazon S3 Source

## Introduction

The Amazon S3 Source is a [Vanus Connector][vc] which is designed to retrieve S3 events from a specific bucket and
transform them into CloudEvents based on [CloudEvents Adapter specification][ceas].

This connector allows users to specify a SQS queue to receive S3 event notification messages.
It will automatically create a SQS queue if you don't specify yours.

The original S3 events looks like:

```json
{
  "Records": [
    {
      "eventVersion": "2.1",
      "eventSource": "aws:s3",
      "awsRegion": "us-west-2",
      "eventTime": "1970-01-01T00:00:00.000Z",
      "eventName": "ObjectCreated:Put",
      "userIdentity": {
        "principalId": "AIDAJDPLRKLG7UEXAMPLE"
      },
      "requestParameters": {
        "sourceIPAddress": "127.0.0.1"
      },
      "responseElements": {
        "x-amz-request-id": "C3D13FE58DE4C810",
        "x-amz-id-2": "FMyUVURIY8/IgAtTv8xRjskZQpcIZ9KG4V5Wp6S7S/JRWeUWerMUE5JgHvANOjpD"
      },
      "s3": {
        "s3SchemaVersion": "1.0",
        "configurationId": "testConfigRule",
        "bucket": {
          "name": "mybucket",
          "ownerIdentity": {
            "principalId": "A3NL1KOZZKExample"
          },
          "arn": "arn:aws:s3:::mybucket"
        },
        "object": {
          "key": "HappyFace.jpg",
          "size": 1024,
          "eTag": "d41d8cd98f00b204e9800998ecf8427e",
          "versionId": "096fKKXTRTtl3on89fVO.nfljtsv6qko",
          "sequencer": "0055AED6DCD90281E5"
        }
      }
    }
  ]
}
```

which is converted to:

```json
{
  "id": "C3D13FE58DE4C810.FMyUVURIY8/IgAtTv8xRjskZQpcIZ9KG4V5Wp6S7S/JRWeUWerMUE5JgHvANOjpD",
  "source": "aws:s3.us-west-2.mybucket",
  "specversion": "V1",
  "type": "com.amazonaws.s3.ObjectCreated:Put",
  "datacontenttype": "application/json",
  "subject": "HappyFace.jpg",
  "time": "1970-01-01T00:00:00.000Z",
  "data": {
    "s3": {
      "s3SchemaVersion": "1.0",
      "configurationId": "testConfigRule",
      "bucket": {
        "name": "mybucket",
        "ownerIdentity": {
          "principalId": "A3NL1KOZZKExample"
        },
        "arn": "arn:aws:s3:::mybucket"
      },
      "object": {
        "key": "HappyFace.jpg",
        "size": 1024,
        "eTag": "d41d8cd98f00b204e9800998ecf8427e",
        "versionId": "096fKKXTRTtl3on89fVO.nfljtsv6qko",
        "sequencer": "0055AED6DCD90281E5"
      }
    }
  }
}
```

## Quick Start

This section will show you how Amazon S3 Source converts S3 events to a CloudEvent.

### Prerequisites

- Have a container runtime (i.e., docker).
- Have an AWS S3 bucket.
- AWS IAM [Access Key][accesskey].
- AWS permissions for the IAM user:
  - s3:PutBucketNotification
  - sqs:ListQueues
  - sqs:GetQueueUrl
  - sqs:ReceiveMessage
  - sqs:GetQueueAttributes
  - sqs:CreateQueue
  - sqs:SetQueueAttributes
  - sqs:DeleteMessage

### Create the config file

```shell
cat << EOF > config.yml
target: http://localhost:31081
aws:
  access_key_id: AKIAIOSFODNN7EXAMPLE
  secret_access_key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
s3_bucket_arn: "arn:aws:s3:::<buckeName>"
s3_events: ["s3:ObjectCreated:*","s3:ObjectRemoved:*"]
region: "us-west-2"
EOF
```

| Name                  | Required | Default | Description                                                                                                                  |
| :-------------------- | :------: | :-----: | :--------------------------------------------------------------------------------------------------------------------------- |
| target                |   YES    |         | the target URL to send CloudEvents                                                                                           |
| aws.access_key_id     |   YES    |         | the AWS IAM [Access Key][accesskey]                                                                                          |
| aws.secret_access_key |   YES    |         | the AWS IAM [Secret Key][accesskey]                                                                                          |
| s3_bucket_arn         |   YES    |         | your S3 bucket arn, example: "arn:aws:s3:::mybucket"                                                                         |
| s3_events             |   YES    |         | it is an array consisting of [s3 events][s3event] you're interested in. example: ["s3:ObjectCreated:*","s3:ObjectRemoved:*"] |
| region                |    NO    |         | it describes where the SQS queue will be created at. This field is only required when you didn't specify your sqsArn.        |
| sqs_arn               |    NO    |         | it is the arn of your SQS queue. The Amazon S3 Source will create a queue located at `region` if this field is omitted.      |

The Amazon S3 Source tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-aws-s3 public.ecr.aws/vanus/connector/source-aws-s3
```

### Test

Open a terminal and use the following command to run a Display sink, which receives and prints CloudEvents.

```shell
docker run -it --rm \
  -p 31081:8080 \
  --name sink-display public.ecr.aws/vanus/connector/sink-display
```

Make sure the `target` value in your config file is `http://localhost:31081` so that the Source can send CloudEvents to our Display Sink.

Open [AWS S3 Console](https://s3.console.aws.amazon.com), select the bucket and upload a file.

Here is the sort of CloudEvent you should expect to receive in the Display Sink:

```json
{
  "id": "C3D13FE58DE4C810.FMyUVURIY8/IgAtTv8xRjskZQpcIZ9KG4V5Wp6S7S/JRWeUWerMUE5JgHvANOjpD",
  "source": "aws:s3.us-west-2.mybucket",
  "specversion": "V1",
  "type": "com.amazonaws.s3.ObjectCreated:Put",
  "datacontenttype": "application/json",
  "subject": "HappyFace.jpg",
  "time": "1970-01-01T00:00:00.000Z",
  "data": {
    "s3": {
      "s3SchemaVersion": "1.0",
      "configurationId": "testConfigRule",
      "bucket": {
        "name": "mybucket",
        "ownerIdentity": {
          "principalId": "A3NL1KOZZKExample"
        },
        "arn": "arn:aws:s3:::mybucket"
      },
      "object": {
        "key": "HappyFace.jpg",
        "size": 1024,
        "eTag": "d41d8cd98f00b204e9800998ecf8427e",
        "versionId": "096fKKXTRTtl3on89fVO.nfljtsv6qko",
        "sequencer": "0055AED6DCD90281E5"
      }
    }
  }
}
```

### Clean

```shell
docker stop source-aws-s3 sink-display
```

## Run in Kubernetes

```shell
kubectl apply -f source-aws-s3.yaml
```

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: source-aws-s3
  namespace: vanus
data:
  config.yml: |-
    "target": "http://vanus-gateway.vanus:8080/gateway/quick_start"
    aws:
      access_key_id: AKIAIOSFODNN7EXAMPLE
      secret_access_Key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
    s3_bucket_arn: "arn:aws:s3:::mybucket"
    s3_events: ["s3:ObjectCreated:*","s3:ObjectRemoved:*"]
    region: "us-west-2"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: source-aws-s3
  namespace: vanus
  labels:
    app: source-aws-s3
spec:
  selector:
    matchLabels:
      app: source-aws-s3
  replicas: 1
  template:
    metadata:
      labels:
        app: source-aws-s3
    spec:
      containers:
        - name: source-aws-s3
          image: public.ecr.aws/vanus/connector/source-aws-s3
          imagePullPolicy: Always
          volumeMounts:
            - name: config
              mountPath: /vanus-connect/config
      volumes:
        - name: config
          configMap:
            name: source-aws-s3
```

## Integrate with Vanus

This section shows how a source connector can send CloudEvents to a running [Vanus cluster](https://github.com/vanus-labs/vanus).

### Prerequisites

- Have a running K8s cluster
- Have a running Vanus cluster
- Vsctl Installed

1. Export the VANUS_GATEWAY environment variable (the ip should be a host-accessible address of the vanus-gateway service)

```shell
export VANUS_GATEWAY=192.168.49.2:30001
```

2. Create an eventbus

```shell
vsctl eventbus create --name quick-start
```

3. Update the target config of the Amazon S3 Source

```yaml
target: http://192.168.49.2:30001/gateway/quick-start
```

4. Run the Amazon S3 Source

```shell
kubectl apply -f source-aws-s3.yaml
```

[vc]: https://docs.vanus.ai/introduction/concepts#vanus-connect
[ceas]: https://github.com/cloudevents/spec/blob/main/cloudevents/adapters/aws-s3.md
[s3event]: https://docs.aws.amazon.com/AmazonS3/latest/userguide/notification-how-to-event-types-and-destinations.html
[accesskey]: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html
