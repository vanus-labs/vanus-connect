---
title: Amazon S3
---

# Amazon S3 Sink

## Introduction

The Amazon S3 (Simple Storage Service) Sink is a [Vanus Connector][vc] that aims to cache incoming CloudEvents in local
files and upload them to AWS S3. The files uploaded to AWS S3 will be **partitioned by time**. S3 Sink
supports **`HOURLY` and `DAILY`** partitioning.

### Features

- **Upload scheduled interval**: S3 Sink supports **scheduled periodic check** for closing and uploading files to S3.
  When the **time interval** reaches the threshold, the file will be directly closed and uploaded, regardless of whether
  the file is full. The interval defaults to **60 seconds**.
- **Flush size**: When file reaches **flush size**, it will be uploaded automatically. The flush size defaults to **
  1000**.
- **Partition**: CloudEvents uploaded by S3 Sink are stored in S3 **partitioned**. S3 Sink create storage path in S3
  according to current time. Time-based partitioning options are daily or hourly.

### Limitations

- **valid schema**: Data S3 Sink received must conform **[CloudEvents Schema Registry][ce-schema]**.
- **Limitation of rate for receiving CloudEvents**: S3 Sink Connector will check the flush size and upload scheduled
  interval one per second. Therefore, the number of CloudEvents sent to S3 Sink per second should be less than flush
  size.

## Quickstart

This quick start will guide you through the process of running an Amazon S3 Sink connector.

### Prerequisites

- Have a container runtime (i.e., docker).
- An Amazon S3 bucket.
- AWS IAM [Access Key][accessKey].
- AWS permissions for the IAM user:
    - s3:PutObject

### Create the config file

```shell
cat << EOF > config.yml
port: 8080
aws:
  access_key_id: AKIAIOSFODNN7EXAMPLE
  secret_access_key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
region: "us-west-2"
bucket: "mybucket"
scheduled_interval: 10
EOF
```

| Name                  | Required | Default | Description                                                                                                                                                                                                                                                                                                    |
|:----------------------|:--------:|:--------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| port                  |    NO    | 8080    | the port which S3 Sink listens on                                                                                                                                                                                                                                                                              |
| aws.access_key_id     |   YES    |         | the AWS IAM [Access Key][accessKey]                                                                                                                                                                                                                                                                            |                                                                                         |
| aws.secret_access_key |   YES    |         | the AWS IAM [Secret Key][accessKey]                                                                                                                                                                                                                                                                            |                                                                                        |
| region                |   YES    |         | the S3 bucket region                                                                                                                                                                                                                                                                                           |
| bucket                |   YES    |         | the S3 bucket name                                                                                                                                                                                                                                                                                             |
| flush_size            |    NO    | 1000    | the number of CloudEvents cached to the local file before S3 Sink upload the file                                                                                                                                                                                                                              |
| scheduled_interval    |    NO    | 60      | the maximum time interval between S3 Sink closing and uploading files which unit is second.                                                                                                                                                                                                                    |
| time_interval         |    NO    | HOURLY  | the partitioning interval of files have been uploaded to the S3. S3 Sink supports `HOURLY` and `DAILY` time interval. For example, when `timeInterval` is `HOURLY`, files uploaded between 3 pm and 4 pm will be partitioned to one path, while files uploaded after 4 pm will be partitioned to another path  |

The Amazon S3 Sink tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the
position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm \
  -p 31080:8080 \
  -v ${PWD}:/vanus-connect/config \
  --name sink-aws-s3 public.ecr.aws/vanus/connector/sink-aws-s3
```

### Test

Open a terminal and use following command to send a CloudEvent to the Sink.

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
     "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "quickstart",
    "specversion": "1.0",
    "type": "quickstart",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "data": {
        "myData": "simulation event data"
    }
}'
```

Open [AWS S3 Console](https://s3.console.aws.amazon.com), select the bucket and verify the file has uploaded.

### Clean resource

```shell
docker stop sink-aws-s3
```

## Run in Kubernetes

```shell
kubectl apply -f sink-aws-s3.yaml
```

```yaml
apiVersion: v1
kind: Service
metadata:
  name: sink-aws-s3
  namespace: vanus
spec:
  selector:
    app: sink-aws-s3
  type: ClusterIP
  ports:
    - port: 8080
      name: sink-aws-s3
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: sink-aws-s3
  namespace: vanus
data:
  config.yml: |-
    port: 8080
    aws:
      access_key_id: AKIAIOSFODNN7EXAMPLE
      secret_access_Key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
    region: "us-west-2"
    bucket: "mybucket"
    scheduled_interval: 10
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sink-aws-s3
  namespace: vanus
  labels:
    app: sink-aws-s3
spec:
  selector:
    matchLabels:
      app: sink-aws-s3
  replicas: 1
  template:
    metadata:
      labels:
        app: sink-aws-s3
    spec:
      containers:
        - name: sink-aws-s3
          image: public.ecr.aws/vanus/connector/sink-aws-s3
          imagePullPolicy: Always
          ports:
            - containerPort: 8080
                name: http
          volumeMounts:
            - name: config
              mountPath: /vanus-connector/config
      volumes:
        - name: config
          configMap:
            name: sink-aws-s3
```

## Integrate with Vanus

This section shows how a sink connector can receive CloudEvents from a
running [Vanus cluster](https://github.com/linkall-labs/vanus).

1. Run the sink-aws-s3.yaml

```shell
kubectl apply -f sink-aws-s3.yaml
```

2. Create an eventbus

```shell
vsctl eventbus create --name quick-start
```

3. Create a subscription (the sink should be specified as the sink service address or the host name with its port)

```shell
vsctl subscription create \
  --name quick-start \
  --eventbus quick-start \
  --sink 'http://sink-aws-s3:8080'
```

[vc]: https://www.vanus.dev/introduction/concepts#vanus-connect
[accessKey]: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html
[ce-schema]: https://github.com/cloudevents/spec/blob/main/schemaregistry/spec.md
