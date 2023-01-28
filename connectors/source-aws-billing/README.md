---
title: Amazon Billing
---

# Amazon Billing Source

## Introduction

The Amazon Billing Source is a [Vanus Connector][vc] which aims to convert billing data
to a CloudEvent. The Amazon Billing Source use [AWS Cost Explorer][awsbill] api pull yesterday data
at fix time.

The billing data is converted to:

```json
{
  "specversion": "1.0",
  "id": "4395ffa3-f6de-443c-bf0e-bb9798d26a1d",
  "source": "cloud.aws.billing",
  "type": "aws.service.daily",
  "datacontenttype": "application/json",
  "time": "2022-06-14T07:05:55.777689Z",
  "data": {
    "date": "2023-01-15",
    "service": "CloudWatch Events",
    "amortizedCost": {
      "amount": "0.0009189721",
      "unit": "USD"
    },
    "blendedCost": {
      "amount": "0.0009189721",
      "unit": "USD"
    },
    "netAmortizedCost": {
      "amount": "0.0009189721",
      "unit": "USD"
    },
    "netUnblendedCost": {
      "amount": "0.0009189721",
      "unit": "USD"
    },
    "unblendedCost": {
      "amount": "0.0009189721",
      "unit": "USD"
    }
  }
}
```

## Quick Start

This section shows how Amazon Billing Source convert billing data to a CloudEvent.

### Prerequisites

- Have a container runtime (i.e., docker).
- AWS IAM [Access Key][accessKey].
- AWS permissions `ce:GetCostAndUsage` for the IAM user.

### Create the config file

```shell
cat << EOF > config.yml
target: http://localhost:31081
secret:
  access_key_id: AKIAIOSFODNN7EXAMPLE
  secret_access_key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
EOF
```

| Name              | Required  | Default                            | Description                                               |
|:------------------|:----------|:-----------------------------------|:----------------------------------------------------------|
| target            | required  |                                    | the target URL will send CloudEvents to                   |
| endpoint          | optional  | https://ce.us-east-1.amazonaws.com | the AWS cost explorer api endpoint                        |
| pull_hour         | optional  | 2                                  | Amazon Billing Source pull AWS billing data at which hour |
| access_key_id     | required  |                                    | the AWS IAM [Access Key][accessKey]                       |
| secret_access_key | required  |                                    | the AWS IAM [Secret Key][accessKey]                       |

The Amazon Billing Source tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-aws-billing public.ecr.aws/vanus/connector/source-aws-billing
```

### Test

Open a terminal and use the following command to run a Display sink, which receives and prints CloudEvents.

```shell
docker run -it --rm \
  -p 31081:8080 \
  --name sink-display public.ecr.aws/vanus/connector/sink-display
```

Make sure the `target` value in your config file is `http://localhost:31081` so that the Source can send CloudEvents to our Display Sink.

then you view events like: 

```json
{
  "specversion": "1.0",
  "id": "4395ffa3-f6de-443c-bf0e-bb9798d26a1d",
  "source": "cloud.aws.billing",
  "type": "aws.service.daily",
  "datacontenttype": "application/json",
  "time": "2022-06-14T07:05:55.777689Z",
  "data": {
    "date": "2023-01-15",
    "service": "CloudWatch Events",
    "amortizedCost": {
      "amount": "0.0009189721",
      "unit": "USD"
    },
    "blendedCost": {
      "amount": "0.0009189721",
      "unit": "USD"
    },
    "netAmortizedCost": {
      "amount": "0.0009189721",
      "unit": "USD"
    },
    "netUnblendedCost": {
      "amount": "0.0009189721",
      "unit": "USD"
    },
    "unblendedCost": {
      "amount": "0.0009189721",
      "unit": "USD"
    }
  }
}
```

### Clean

```shell
docker stop source-aws-billing sink-display
```

## Run in Kubernetes

```shell
kubectl apply -f source-aws-billing.yaml
```

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: source-aws-billing
  namespace: vanus
data:
  config.yml: |-
    target: "http://vanus-gateway.vanus:8080/gateway/quick_start"
    secret:
      access_key_id: AKIAIOSFODNN7EXAMPLE
      secret_access_key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: source-aws-billing
  namespace: vanus
  labels:
    app: source-aws-billing
spec:
  selector:
    matchLabels:
      app: source-aws-billing
  replicas: 1
  template:
    metadata:
      labels:
        app: source-aws-billing
    spec:
      containers:
        - name: source-aws-billing
          image: public.ecr.aws/vanus/connector/source-aws-billing
          imagePullPolicy: Always
          volumeMounts:
            - name: source-aws-billing-config
              mountPath: /vanus-connector/config
      volumes:
        - name: source-aws-billing-config
          configMap:
            name: source-aws-billing
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

3. Update the target config of the Amazon Billing Source
```yaml
target: http://192.168.49.2:30001/gateway/quick-start
```

4. Run the Amazon Billing Source
```shell
docker run --network=host \
  --rm \
  -v ${PWD}:/vanus-connect/config \
  --name source-aws-billing public.ecr.aws/vanus/connector/source-aws-billing
```

[vc]: https://www.vanus.dev/introduction/concepts#vanus-connect
[awsbill]: https://docs.aws.amazon.com/aws-cost-management/latest/APIReference/API_GetCostAndUsage.html
[accessKey]: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html
