---
title: Amazon Billing
---

# Amazon Billing Source

## Introduction

The Amazon Billing Source is a [Vanus Connector][vc] which aims to convert billing data
to a CloudEvent. The Amazon Billing Source uses [AWS Cost Explorer][awsbill] api and pulls 
billing data from the previous day at a fixed time.

The billing data is converted to:

```json
{
  "specversion": "1.0",
  "id": "026046e2-3cb0-4116-895e-c77877072dd2",
  "source": "cloud.aws.billing",
  "type": "aws.service.daily",
  "datacontenttype": "application/json",
  "time": "2023-01-28T06:11:10.012579049Z",
  "data": {
    "date": "2023-01-27",
    "service": "Amazon Relational Database Service",
    "amortizedCost": {
      "amount": "0.2672917174",
      "unit": "USD"
    },
    "blendedCost": {
      "amount": "0.2672917174",
      "unit": "USD"
    },
    "netAmortizedCost": {
      "amount": "0.2672917174",
      "unit": "USD"
    },
    "netUnblendedCost": {
      "amount": "0.2672917174",
      "unit": "USD"
    },
    "unblendedCost": {
      "amount": "0.2672917174",
      "unit": "USD"
    }
  }
}
```

## Quick Start

This section shows how Amazon Billing Source converts billing data to a CloudEvent.

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

| Name              | Required | Default                            | Description                                               |
|:------------------|:---------|:-----------------------------------|:----------------------------------------------------------|
| target            | YES      |                                    | the target URL to send CloudEvents                        |
| endpoint          | NO       | https://ce.us-east-1.amazonaws.com | the AWS cost explorer api endpoint                        |
| pull_hour         | NO       | 2                                  | specify the hour at which the billing data will be pulled |
| access_key_id     | YES      |                                    | the AWS IAM [Access Key][accessKey]                       |
| secret_access_key | YES      |                                    | the AWS IAM [Secret Key][accessKey]                       |

The Amazon Billing Source tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-aws-billing public.ecr.aws/vanus/connector/source-aws-billing
```

### Test

Before starting Amazon billing Source use the following command to run a Display sink, which will receive and prints the incoming CloudEvents.

```shell
docker run -it --rm \
  -p 31081:8080 \
  --name sink-display public.ecr.aws/vanus/connector/sink-display
```

Make sure the `target` value in your config file is `http://localhost:31081` so that the Source can send CloudEvents to our Display Sink.

Here is the sort of CloudEvent you should expect to receive in the Display Sink: 

```json
{
  "specversion": "1.0",
  "id": "026046e2-3cb0-4116-895e-c77877072dd2",
  "source": "cloud.aws.billing",
  "type": "aws.service.daily",
  "datacontenttype": "application/json",
  "time": "2023-01-28T06:11:10.012579049Z",
  "data": {
    "date": "2023-01-27",
    "service": "Amazon Relational Database Service",
    "amortizedCost": {
      "amount": "0.2672917174",
      "unit": "USD"
    },
    "blendedCost": {
      "amount": "0.2672917174",
      "unit": "USD"
    },
    "netAmortizedCost": {
      "amount": "0.2672917174",
      "unit": "USD"
    },
    "netUnblendedCost": {
      "amount": "0.2672917174",
      "unit": "USD"
    },
    "unblendedCost": {
      "amount": "0.2672917174",
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
kubectl apply -f source-aws-billing.yaml
```

[vc]: https://www.vanus.dev/introduction/concepts#vanus-connect
[awsbill]: https://docs.aws.amazon.com/aws-cost-management/latest/APIReference/API_GetCostAndUsage.html
[accessKey]: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html
