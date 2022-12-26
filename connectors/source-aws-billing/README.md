---
title: AWS Billing
---

# AWS Billing Source

## Introduction

The AWS Billing Source is a [Vance Connector][vc] which use [AWS Cost Explorer][awsbill] api pull yesterday billing data
by fix time.The data group by aws service For example,billing data output a CloudEvent looks like:

```json
{
  "specversion": "1.0",
  "id": "4395ffa3-f6de-443c-bf0e-bb9798d26a1d",
  "source": "cloud.aws.billing",
  "type": "aws.service.daily",
  "datacontenttype": "application/json",
  "time": "2022-06-14T07:05:55.777689Z",
  "data": {
    "vanceSource": "cloud.aws.billing",
    "vanceType": "aws.service.daily",
    "date": "2022-06-13",
    "service": "Amazon Elastic Compute Cloud - Compute",
    "amount": "12.294",
    "unit": "USD"
  }
}
```

## AWS Billing Source Configs

### Config

| name      | requirement | default                             | description                                           |
|:----------|:------------|:------------------------------------|:------------------------------------------------------|
| target    | required    |                                     | target URL will send CloudEvents to                   |
| endpoint  | optional    | https://ce.us-east-1.amazonaws.com  | the aws cost explorer api endpoint                    |
| pull_hour | optional    | 2                                   | aws billing source pull billing data time(unit hour)  |

### Secret

| name              | requirement | default  | description                                  |
|-------------------|-------------|----------|----------------------------------------------|
| access_key_id     | required    |          | the aws account [accessKeyID][accessKey]     |
| secret_access_Key | required    |          | the aws account [secretAccessKey][accessKey] |

## AWS Billing Source Image

> public.ecr.aws/vanus/connector/source-aws-billing

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
 docker run --rm -v ${PWD}:/vance/config -v ${PWD}:/vance/secret public.ecr.aws/vanus/connector/source-aws-billing
```

### K8S

```shell
  kubectl apply -f source-aws-billing.yaml
```

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md
[awsbill]: https://docs.aws.amazon.com/aws-cost-management/latest/APIReference/API_GetCostAndUsage.html
[accessKey]: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html
