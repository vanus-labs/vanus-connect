---
title: AliCloud Billing
---

# AliCloud Billing Source

## Introduction

The AliCloud Billing Source is a [Vance Connector][vc] which use [AliCloud billing][alibill] api pull yesterday billing
data by fix time.The data group by product

For example,billing data output a CloudEvent looks like:

```json
{
  "specversion": "1.0",
  "id": "bd64e9e0-cd46-43f1-95fa-2008b6b49e85",
  "source": "cloud.alicloud.billing",
  "type": "alicloud.account_billing.daily",
  "datacontenttype": "application/json",
  "time": "2022-06-14T07:01:55.277687Z",
  "data": {
    "VanceSource": "cloud.alicloud.billing",
    "VanceType": "alicloud.account_billing.daily",
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

## AliCloud Billing Source Configs

### Config

| name      | requirement | default               | description                                               |
|:----------|:------------|:----------------------|:----------------------------------------------------------|
| target    | required    |                       | target URL will send CloudEvents to                       |
| endpoint  | optional    | business.aliyuncs.com | the AliCloud business api endpoint                        |
| pull_hour | optional    | 2                     | AliCloud billing source pull billing data time(unit hour) |

### Secret

| name              | requirement | default  | description                                       |
|-------------------|-------------|----------|---------------------------------------------------|
| access_key_id     | required    |          | the AliCloud account [accessKeyID][accessKey]     |
| secret_access_Key | required    |          | the AliCloud account [secretAccessKey][accessKey] |

## AliCloud Billing Source Image

> public.ecr.aws/vanus/connector/source-alicloud-billing

## Deploy

### Docker

#### create config file

refer [config](#Config) to create `config.yml`. for example:

```yaml
"target": "http://localhost:8080"
```

#### create secret file

refer [secret](#Secret) to create `secret.yml`. for example:

```yaml
"access_key_id": "xxxxxx"
"secret_access_key": "xxxxxx"
```

#### run

```shell
 docker run --rm -v ${PWD}:/vance/config -v ${PWD}:/vance/secret public.ecr.aws/vanus/connector/source-alicloud-billing
```

### K8S

```shell
  kubectl apply -f source-alicloud-billing.yaml
```

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md
[alibill]: https://help.aliyun.com/document_detail/142608.html
[accessKey]: https://help.aliyun.com/document_detail/38738.html
