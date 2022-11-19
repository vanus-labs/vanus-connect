# Doris Sink

## Introduction

The Doris Sink is a [Vance Connector][vc] which aims to handle incoming CloudEvents in a way that extracts the `data`
part of the original event and deliver these extracted `data` to [Doris][doris]. The Sink use [Stream Load][stream load]
way to import data. 

For example, if the incoming CloudEvent looks like:

```json
{
  "specversion": "1.0",
  "id": "4395ffa3-f6de-443c-bf0e-bb9798d26a1d",
  "source": "vance.source.test",
  "type": "vance.type.test",
  "datacontenttype": "application/json",
  "time": "2022-11-20T07:05:55.777689Z",
  "data": {
    "id": 1,
    "username": "name",
    "birthday": "2022-11-20"
  }
}
```

The Doris Sink will extract `data` field write to [Doris][doris] table like:

```text
+------+----------+------------+
| id   | username | birthday   |
+------+----------+------------+
|    1 | name     | 2022-11-20 |
+------+----------+------------+
```

## Doris Sink Configs

Users can specify their configs by either setting environments variables or mount a config.yaml to
`/vance/config/config.yaml` when they run the connector. Find examples of setting configs [here][config].

### Config Fields of Doris Sink

| name          | requirement | description                                                                |
|---------------|-------------|----------------------------------------------------------------------------|
| v_port        | optional    | v_port is used to specify the port Doris Sink is listening on,default 8080 |
| fenodes       | required    | doris fenodes, example: "17.0.0.1:8003"                                    |
| username      | required    | doris username                                                             |
| password      | optional    | doris password                                                             |
| db_name       | required    | doris database name                                                        |
| table_name    | required    | doris table name                                                           |
| stream_load   | optional    | doris stream load properties, map struct                                   |
| load_interval | optional    | doris stream load interval, unit second, default 5                         |
| load_size     | optional    | doris stream load max body size, default 10 * 1024 * 1024                  |
| timeout       | optional    | doris stream load timeout, unit second, default 30                         |

## Doris Sink Image

> docker.io/vancehub/sink-doris

## Local Development

You can run the sink codes of the Doris Sink locally as well.

### Building

```shell
cd connectors/sink-doris
go build -o bin/sink cmd/main.go
```

### Running

```shell
bin/sink
```

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md
[config]: https://github.com/linkall-labs/vance-docs/blob/main/docs/connector.md
[doris]: https://doris.apache.org/docs/summary/basic-summary
[stream load]: https://doris.apache.org/docs/dev/data-operate/import/import-way/stream-load-manual/
