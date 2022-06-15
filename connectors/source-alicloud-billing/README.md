# AliCloud Billing Source

## Overview

A [Vance Connector][vc] which transforms AliCloud billing to CloudEvents and deliver them to the target URL.

## User Guidelines

### Connector Introduction

The AliCloud Billing Source is a [Vance Connector][vc] which use [AliCloud billing][alibill] api pull yesterday billing data by fix time. 
The data group by product

For example, output a CloudEvent looks like:

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

## Source Configs

Users can specify their configs by either setting environments variables or mount a config.json to
`/vance/config/config.json` when they run the connector. Find examples of setting configs [here][config].

### Config Fields

| name              | requirement | description                                                           |
|-------------------|-------------|-----------------------------------------------------------------------|
| v_target          | required    | target URL will send CloudEvents to                                   |
| access_key_id     | required    | the AliCloud account accessKeyID                                      |
| secret_access_Key | required    | the AliCloud account secretAccessKey                                  | 
| endpoint          | optional    | the AliCloud business api endpoint,default business.aliyuncs.com      |
| pull_hour         | optional    | AliCloud billing source pull billing data time(unit hour),default 2   |


## HTTP Source Image

> public.ecr.aws/vanus/connector/alicloudbill

## Local Development

You can run the source codes of the AliCloud Billing Source locally as well.

### Building

```shell
$ cd connectors/source-alicloud-billing
$ go build -o bin/source cmd/main.go
```

### Add and modify config

```json
{
  "access_key_id": "xxxxxx",
  "secret_access_key":"xxxxxx",
  "v_target": "http://localhost:8080"
}
```

### Running

```shell
$ ./bin/source
```

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md
[config]: https://github.com/linkall-labs/vance-docs/blob/main/docs/connector.md
[alibill]: https://help.aliyun.com/document_detail/142608.html