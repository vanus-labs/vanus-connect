# Elasticsearch Sink

## Overview

A [Vance Connector][vc] which receives CloudEvents and deliver specific data to elasticsearch cluster

## User Guidelines

### Connector Introduction

The Elasticsearch Sink is a [Vance Connector][vc] which aims to handle incoming CloudEvents in a way that extracts the `data` part of the
original event and deliver these extracted `data` to elasticsearch cluster [index](index)


For example, if the incoming CloudEvent looks like:

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

The Elasticsearch Sink will write `data` to elasticsearch looks like:
```json
{
  "_index": "billing",
  "_type": "_doc",
  "_id": "CqFnBIEBzJc0Oa5TERDD",
  "_version": 1,
  "_source": {
    "vanceSource": "cloud.aws.billing",
    "vanceType": "aws.service.daily",
    "date": "2022-06-13",
    "service": "Amazon Elastic Compute Cloud - Compute",
    "amount": "12.294",
    "unit": "USD"
  }
}
```

## Sink Configs

Users can specify their configs by either setting environments variables or mount a config.json to
`/vance/config/config.json` when they run the connector. Find examples of setting configs [here][config].

### Config Fields

| name       | requirement | description                                                                        |
|------------|-------------|------------------------------------------------------------------------------------|
| v_port     | optional    | v_port is used to specify the port Elasticsearch Sink is listening on,default 8080 |
| address    | required    | elasticsearch cluster address, multi split by ","                                  | 
| index_name | required    | elasticsearch index name which to be write                                         | 
| username   | optional    | elasticsearch cluster username                                                     |
| password   | optional    | elasticsearch cluster password                                                     |


## Elasticsearch Sink Image

> public.ecr.aws/vanus/connector/essink

## Local Development

You can run the sink codes of the Sink Elasticsearch locally as well.

### Building

```shell
$ cd connectors/sink-elastisearch
$ go build -o bin/sink cmd/main.go
```

### Add and modify config

```json
{
  "address": "http://localhost:9200,http://uri:port",
  "index_name": "billing",
  "username":"elastic",
  "password": "elastic"
}
```

### Running

```shell
$ ./bin/sink
```

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md
[config]: https://github.com/linkall-labs/vance-docs/blob/main/docs/connector.md
[index]: https://www.elastic.co/guide/en/elasticsearch/reference/current/indices-create-index.html