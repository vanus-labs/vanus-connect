---
title: Alibaba Cloud Billing
---

# Alibaba Cloud Billing Source

## Introduction

The Alibaba Cloud Billing Source is a [Vanus Connector][vc] which aims to convert billing data to CloudEvents. The
Alibaba Cloud Billing Source use [Alibaba Cloud Billing][alibill] api and pulls billing data from the previous day at a
fixed time.

The billing data is converted to:

```json
{
  "specversion": "1.0",
  "id": "bd64e9e0-cd46-43f1-95fa-2008b6b49e85",
  "source": "cloud.alibaba.billing",
  "type": "alibaba.billing.daily",
  "datacontenttype": "application/json",
  "time": "2022-06-14T07:01:55.277687Z",
  "data": {
    "AdjustAmount": 0,
    "BillAccountID": "123456",
    "BillAccountName": "aliyun23456",
    "BillingDate": "2022-06-13",
    "BizType": "",
    "CashAmount": 0,
    "CostUnit": "",
    "Currency": "CNY",
    "DeductedByCashCoupons": 0,
    "DeductedByCoupons": 0,
    "DeductedByPrepaidCard": 0,
    "InvoiceDiscount": 0,
    "OutstandingAmount": 0,
    "OwnerId": "123456",
    "OwnerName": "aliyun23456",
    "PaymentAmount": 0.23,
    "PipCode": "disk",
    "PretaxAmount": 0,
    "PretaxGrossAmount": 0.2352,
    "ProductCode": "disk",
    "ProductName": "",
    "SubscriptionType": "PayAsYouGo"
  }
}
```

## Quick Start

This section shows how Alibaba Cloud Billing Source converts billing data to a CloudEvent.

### Prerequisites

- Have a container runtime (i.e., docker).
- Alibaba Cloud RAM [Access Key][accessKey].
- Alibaba Cloud RAM permissions `bssapi:QueryAccountBill` for the RAM user.

### Create the config file

```shell
cat << EOF > config.yml
target: http://localhost:31081
secret:
  access_key_id: LTAI5t7eneojm6aeRKncB7V9
  secret_access_key: 26Fno8k0F7jtL2hUIrm4JTRF16jPLL
EOF
```

| Name              | Required | Default               | Description                                                                |
|:------------------|:--------:|:----------------------|:---------------------------------------------------------------------------|
| target            |   YES    |                       | the target URL to send CloudEvents                                         |
| endpoint          |    NO    | business.aliyuncs.com | the Alibaba Cloud business api endpoint                                    |
| pull_hour         |    NO    | 2                     | specify the hour at which the billing data will be pulled, value is [1,23] |
| pull_zone         |    NO    | UTC                   | pull billing data hour time zone                                           |
| access_key_id     |   YES    |                       | the Alibaba Cloud [access Key][accessKey]                                  |
| secret_access_Key |   YES    |                       | the Alibaba Cloud [secret Key][accessKey]                                  |

The Alibaba Cloud Billing Source tries to find the config file at `/vanus-connect/config/config.yml` by default. You can
specify the position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-alicloud-billing public.ecr.aws/vanus/connector/source-alicloud-billing
```

### Test

Open a terminal and use the following command to run a Display sink, which receives and prints CloudEvents.

```shell
docker run -it --rm \
  -p 31081:8080 \
  --name sink-display public.ecr.aws/vanus/connector/sink-display
```

Make sure the `target` value in your config file is `http://localhost:31081` so that the Source can send CloudEvents to
our Display Sink.

Here is the sort of CloudEvent you should expect to receive in the Display Sink:

```json
{
  "specversion": "1.0",
  "id": "bd64e9e0-cd46-43f1-95fa-2008b6b49e85",
  "source": "cloud.alibaba.billing",
  "type": "alibaba.billing.daily",
  "datacontenttype": "application/json",
  "time": "2022-06-14T07:01:55.277687Z",
  "data": {
    "AdjustAmount": 0,
    "BillAccountID": "123456",
    "BillAccountName": "aliyun23456",
    "BillingDate": "2022-06-13",
    "BizType": "",
    "CashAmount": 0,
    "CostUnit": "",
    "Currency": "CNY",
    "DeductedByCashCoupons": 0,
    "DeductedByCoupons": 0,
    "DeductedByPrepaidCard": 0,
    "InvoiceDiscount": 0,
    "OutstandingAmount": 0,
    "OwnerId": "123456",
    "OwnerName": "aliyun23456",
    "PaymentAmount": 0.23,
    "PipCode": "disk",
    "PretaxAmount": 0,
    "PretaxGrossAmount": 0.2352,
    "ProductCode": "disk",
    "ProductName": "",
    "SubscriptionType": "PayAsYouGo"
  }
}
```

### Clean

```shell
docker stop source-alicloud-billing sink-display
```

## Run in Kubernetes

```shell
kubectl apply -f source-alicloud-billing.yaml
```

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: source-alicloud-billing
  namespace: vanus
data:
  config.yml: |-
    target: "http://vanus-gateway.vanus:8080/gateway/quick_start"
    secret:
      access_key_id: LTAI5t7eneojm6aeRKncB7V9
      secret_access_key: 26Fno8k0F7jtL2hUIrm4JTRF16jPLL
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: source-alicloud-billing
  namespace: vanus
  labels:
    app: source-alicloud-billing
spec:
  selector:
    matchLabels:
      app: source-alicloud-billing
  replicas: 1
  template:
    metadata:
      labels:
        app: source-alicloud-billing
    spec:
      containers:
        - name: source-alicloud-billing
          image: public.ecr.aws/vanus/connector/source-alicloud-billing
          imagePullPolicy: Always
          volumeMounts:
            - name: source-alicloud-billing-config
              mountPath: /vanus-connect/config
      volumes:
        - name: source-alicloud-billing-config
          configMap:
            name: source-alicloud-billing
```

## Integrate with Vanus

This section shows how a source connector can send CloudEvents to a
running [Vanus cluster](https://github.com/vanus-labs/vanus).

### Prerequisites

- Have a running K8s cluster
- Have a running Vanus cluster
- Vsctl Installed

1. Export the VANUS_GATEWAY environment variable (the ip should be a host-accessible address of the vanus-gateway
   service)

```shell
export VANUS_GATEWAY=192.168.49.2:30001
```

2. Create an eventbus

```shell
vsctl eventbus create --name quick-start
```

3. Update the target config of the Alibaba Cloud Billing Source

```yaml
target: http://192.168.49.2:30001/gateway/quick-start
```

4. Run the Alibaba Cloud Billing Source

```shell
kubectl apply -f source-alicloud-billing.yaml
```

[vc]: https://docs.vanus.ai/introduction/concepts#vanus-connect
[alibill]: https://help.aliyun.com/document_detail/142608.html
[accessKey]: https://help.aliyun.com/document_detail/38738.html
