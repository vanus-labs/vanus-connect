---
title: Amazon SQS
---

# Amazon SQS Source

## Introduction

The Amazon SQS Source is a [Vanus Connector][vc] which is designed to retrieve SQS messages transform them into CloudEvents.

For example, if the incoming message looks like:

```json
{
  "MessageId": "035e183b-275a-44de-95df-f212be1ed4ea",
  "Body": "Hello World"
}
```

Which is converted to:

```json
{
  "specversion": "1.0",
  "id": "035e183b-275a-44de-95df-f212be1ed4ea",
  "source": "cloud.aws.sqs.us-west-2.my-test-queue",
  "type": "com.amazonaws.sqs.message",
  "datacontenttype": "text/plain",
  "time": "2022-08-02T11:01:13.828+08:00",
  "data": "Hello World"
}
```

This section shows you how to use Amazon SQS Source to convert SQS message to a CloudEvent.

### Prerequisites

- Have a container runtime (i.e., docker).
- Have an AWS SQS queue.
- AWS IAM [Access Key][accesskey].
- AWS permissions for the IAM user:
  - sqs:GetQueueUrl
  - sqs:ReceiveMessage
  - sqs:DeleteMessage

### Create the config file

```shell
cat << EOF > config.yml
target: http://localhost:31081
aws:
  access_key_id: AKIAIOSFODNN7EXAMPLE
  secret_access_key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
sqs_arn: "arn:aws:sqs:us-west-2:843378899134:myQueue"
EOF
```

| Name                  | Required | Default | Description                         |
| :-------------------- | :------: | :-----: | :---------------------------------- |
| target                |   YES    |         | the target URL to send CloudEvents  |
| aws.access_key_id     |   YES    |         | the AWS IAM [Access Key][accesskey] |
| aws.secret_access_key |   YES    |         | the AWS IAM [Secret Key][accesskey] |
| sqs_arn               |   YES    |         | your SQS ARN                        |

The Amazon SQS Source tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-aws-sqs public.ecr.aws/vanus/connector/source-aws-sqs
```

### Test

Open a terminal and use the following command to run a Display sink, which receives and prints CloudEvents.

```shell
docker run -it --rm \
  -p 31081:8080 \
  --name sink-display public.ecr.aws/vanus/connector/sink-display
```

Make sure the `target` value in your config file is `http://localhost:31081` so that the Source can send CloudEvents to our Display Sink.

After running Display Sink, run the SQS Source

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-aws-sqs public.ecr.aws/vanus/connector/source-aws-sqs
```

Open [AWS SQS Console](https://us-west-2.console.aws.amazon.com/sqs/v2/home?region=us-west-2#/queues), select the queue and send a message.

Here is the sort of CloudEvent you should expect to receive in the Display Sink:

```json
{
  "specversion": "1.0",
  "id": "035e183b-275a-44de-95df-f212be1ed4ea",
  "source": "cloud.aws.sqs.us-west-2.my-test-queue",
  "type": "com.amazonaws.sqs.message",
  "datacontenttype": "text/plain",
  "time": "2022-08-02T11:01:13.828+08:00",
  "data": "Hello World"
}
```

### Clean

```shell
docker stop source-aws-sqs sink-display
```

## Run in Kubernetes

```shell
kubectl apply -f source-aws-sqs.yaml
```

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: source-aws-sqs
  namespace: vanus
data:
  config.yml: |-
    "target": "http://vanus-gateway.vanus:8080/gateway/quick_start"
    aws:
      access_key_id: AKIAIOSFODNN7EXAMPLE
      secret_access_Key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
    sqs_arn: "arn:aws:sqs:us-west-2:843378899134:myQueue"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: source-aws-sqs
  namespace: vanus
  labels:
    app: source-aws-sqs
spec:
  selector:
    matchLabels:
      app: source-aws-sqs
  replicas: 1
  template:
    metadata:
      labels:
        app: source-aws-sqs
    spec:
      containers:
        - name: source-aws-sqs
          image: public.ecr.aws/vanus/connector/source-aws-sqs
          imagePullPolicy: Always
          volumeMounts:
            - name: config
              mountPath: /vanus-connect/config
      volumes:
        - name: config
          configMap:
            name: source-aws-sqs
```

## Integrate with Vanus

This section shows how a source connector can send CloudEvents to a running [Vanus cluster](https://github.com/linkall-labs/vanus).

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

3. Update the target config of the Amazon SQS Source

```yaml
target: http://192.168.49.2:30001/gateway/quick-start
```

4. Run the Amazon SQS Source

```shell
kubectl apply -f source-aws-sqs.yaml
```

[vc]: https://www.vanus.ai/introduction/concepts#vanus-connect
[accesskey]: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html
