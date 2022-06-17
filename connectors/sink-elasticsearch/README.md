# Elasticsearch Sink

## Introduction

The Elasticsearch Sink is a [Vance Connector][vc] which aims to handle incoming CloudEvents in a way that extracts the `data` part of the original event and deliver these extracted `data` to [Elasticsearch][es] cluster

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

Users can specify their configs by either setting environments variables or mount a config.json to
`/vance/config/config.json` when they run the connector. Find examples of setting configs [here][config].

### Config Fields of Elasticsearch Sink

| name       | requirement | description                                                                        |
|------------|-------------|------------------------------------------------------------------------------------|
| v_port     | optional    | v_port is used to specify the port Elasticsearch Sink is listening on,default 8080 |
| address    | required    | elasticsearch cluster address, multi split by ","                                  |
| index_name | required    | elasticsearch index name                                                           |
| username   | optional    | elasticsearch cluster username                                                     |
| password   | optional    | elasticsearch cluster password                                                     |

## Elasticsearch Sink Image

> docker.io/vancehub/sink-elasticsearch

## Local Development

You can run the sink codes of the Elasticsearch Sink locally as well.

### Building

```shell
cd connectors/sink-elastisearch
go build -o bin/sink cmd/main.go
```

### Running

```shell
bin/sink
```

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md
[config]: https://github.com/linkall-labs/vance-docs/blob/main/docs/connector.md
[es]: https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html
