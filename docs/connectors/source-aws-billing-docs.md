# AWS Billing Source

## Introduction

The AWS Billing Source is a [Vance Connector][vc] which use [AWS Cost Explorer][awsbill] api pull yesterday billing data by fix time.The data group by aws service

For example,billing data output a CloudEvent looks like:

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

Users can specify their configs by either setting environments variables or mount a config.json to
`/vance/config/config.json` when they run the connector. Find examples of setting configs [here][config].

### Config Fields of the AWS Billing Source

| name              | requirement | description                                                                     |
|-------------------|-------------|---------------------------------------------------------------------------------|
| v_target          | required    | target URL will send CloudEvents to                                             |
| access_key_id     | required    | the aws account [accessKeyID][accessKey]                                        |
| secret_access_Key | required    | the aws account [secretAccessKey][accessKey]                                    |
| endpoint          | optional    | the aws cost explorer api endpoint,default <https://ce.us-east-1.amazonaws.com> |
| pull_hour         | optional    | aws billing source pull billing data time(unit hour),default 2                  |

## AWS Billing Source Image

> docker.io/vancehub/source-aws-billing

## Local Development

You can run the source codes of the AWS Billing Source locally as well.

### Building

```shell
cd connectors/source-aws-billing
go build -o bin/source cmd/main.go
```

### Running

```shell
bin/source
```

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md
[config]: https://github.com/linkall-labs/vance-docs/blob/main/docs/connector.md
[awsbill]: https://docs.aws.amazon.com/aws-cost-management/latest/APIReference/API_GetCostAndUsage.html
[accessKey]: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html
