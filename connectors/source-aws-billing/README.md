# AWS Billing Source

## Overview

A [Vance Connector][vc] which transforms aws billing to CloudEvents and deliver them to the target URL.

## User Guidelines

### Connector Introduction

The AWS Billing Source is a [Vance Connector][vc] which use aws cost explorer api pull yesterday billing data by fix time. 
The data group by aws service

For example, output a CloudEvent looks like:

``` json
{
  "specversion": "1.0",
  "id": "4395ffa3-f6de-443c-bf0e-bb9798d26a1d",
  "source": "cloud.billing.aws",
  "type": "aws.service.daily",
  "datacontenttype": "application/json",
  "time": "2022-06-14T07:05:55.777689Z",
  "data": {
    "vanceSource": "cloud.billing.aws",
    "vanceType": "aws.service.daily",
    "date": "2022-06-13",
    "service": "Amazon Elastic Compute Cloud - Compute",
    "amount": "12.294",
    "unit": "USD"
  }
}
```

## Vance Connector Configs

The AWS Billing source use Vance cdk develop, so users can specify config by either setting environments variables or mount a config.json to
`/vance/config/config.json` when run the connector. 

Configuration

| name              | requirement | description                                                                   |
|-------------------|-------------|-------------------------------------------------------------------------------|
| vance_sink        | required    | target URL will send CloudEvents to                                           |
| access_key_id     | required    | the aws account accessKeyID                                                   |
| secret_access_Key | required    | the aws account secretAccessKey                                               | 
| endpoint          | optional    | the aws cost explorer api endpoint,default https://ce.us-east-1.amazonaws.com |
| pull_hour         | optional    | aws billing source pull billing data time(unit hour),default 2                |


## HTTP Source Image

> public.ecr.aws/vanus/connector/awsbill

## Local Development

You can run the source codes of the AWS Billing Source locally as well.

### Building

```shell
$ cd connectors/source-aws-billing
$ go build -o bin/source cmd/main.go
```

### Add and modify config

```json
{
  "access_key_id": "xxxxxx",
  "secret_access_key":"xxxxxx",
  "vance_sink": "http://localhost:8080"
}
```

### Running

```shell
$ ./bin/source VANCE_CONFIG_PATH=./config.json
```

[vc]: https://github.com/JieDing/vance-docs/blob/main/docs/concept.md