---
title: Doris
---

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

### Config

| name          | requirement  | default       | description                               |
|:--------------|:-------------|:--------------|:------------------------------------------|
| port          | optional     | 8080          | the port Doris Sink is listening on       |
| fenodes       | required     |               | doris fenodes, example: "17.0.0.1:8003"   |
| db_name       | required     |               | doris database name                       |
| table_name    | required     |               | doris table name                          |
| stream_load   | optional     |               | doris stream load properties, map struct  |
| load_interval | optional     | 5             | doris stream load interval, unit second   |
| load_size     | optional     | 10*1024*1024  | doris stream load max body size           |
| timeout       | optional     | 30            | doris stream load timeout, unit second    |

### Secret

| name          | requirement | default  | description    |
|---------------|-------------|----------|----------------|
| username      | required    |          | doris username |
| password      | required    |          | doris password |

## Doris Sink Image

> public.ecr.aws/vanus/connector/sink-doris

## Deploy

### Docker

#### create config file

refer [config](#Config) to create `config.yml`. for example:

```yaml
"port": 8080
"fenodes": "172.31.57.192:8030"
"db_name": "vance_test"
"table_name": "user"
```

#### create secret file

refer [secret](#Secret) to create `secret.yml`. for example:

```yaml
"username": "vance_test"
"password": "123456"
```

#### run

```shell
 docker run --rm -v ${PWD}:/vance/config -v ${PWD}:/vance/secret public.ecr.aws/vanus/connector/sink-doris
```

### K8S

```shell
  kubectl apply -f sink-doris.yaml
```

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md
[doris]: https://doris.apache.org/docs/summary/basic-summary
[stream load]: https://doris.apache.org/docs/dev/data-operate/import/import-way/stream-load-manual/
