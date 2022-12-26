---
title: Elasticsearch
---

# Elasticsearch Sink

## Introduction

The Elasticsearch Sink is a [Vance Connector][vc] which aims to handle incoming CloudEvents in a way that extracts
the `data` part of the original event and deliver these extracted `data` to [Elasticsearch][es] cluster

For example, if the incoming CloudEvent looks like:

```json
{
  "specversion": "1.0",
  "id": "4395ffa3-f6de-443c-bf0e-bb9798d26a1d",
  "source": "vance.source.test",
  "type": "vance.type.test",
  "datacontenttype": "application/json",
  "time": "2022-06-14T07:05:55.777689Z",
  "data": {
    "date": "2022-06-13",
    "service": "test data",
    "amount": "12.294",
    "unit": "USD"
  }
}
```

The Elasticsearch Sink will extract `data` field write to [Elasticsearch][es] cluster index as a document looks like:

```json
{
  "_index": "vance_test",
  "_type": "_doc",
  "_id": "CqFnBIEBzJc0Oa5TERDD",
  "_version": 1,
  "_source": {
    "date": "2022-06-13",
    "service": "test data",
    "amount": "12.294",
    "unit": "USD"
  }
}
```

## Elasticsearch Sink Configs

### Config

| name        | requirement  | default  | description                                                          |
|:------------|:-------------|:---------|:---------------------------------------------------------------------|
| port        | optional     | 8080     | the port Elasticsearch Sink is listening on                          |
| address     | required     |          | elasticsearch cluster address, multi split by ","                    |
| index_name  | required     |          | elasticsearch index name                                             |
| timeout     | optional     | 10000    | elasticsearch index document timeout, unit millisecond               |
| insert_mode | optional     | insert   | elasticsearch index document type: insert or upsert                  |
| primary_key | optional     |          | elasticsearch index document primary key in event, example: data.id  |

### Secret

| name        | requirement | default  | description                     |
|-------------|-------------|----------|---------------------------------|
| username    | optional    |          | elasticsearch cluster username  |
| password    | optional    |          | elasticsearch cluster password  |

## Image

> public.ecr.aws/vanus/connector/sink-elasticsearch

## Deploy

### Docker

#### create config file

refer [config](#Config) to create `config.yml`. for example:

```yaml
"port": 8080
"address": "http://localhost:9200"
"index_name": "vance_test"
"primary_key": "data.id"
"insert_mode": "upsert"
```

#### create secret file

refer [secret](#Secret) to create `secret.yml`. for example:

```yaml
"username": "elastic"
"password": "elastic"
```

#### run

```shell
 docker run --rm -v ${PWD}:/vance/config -v ${PWD}:/vance/secret public.ecr.aws/vanus/connector/sink-elasticsearch
```

### K8S

```shell
  kubectl apply -f sink-es.yaml
```

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md
[es]: https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html
