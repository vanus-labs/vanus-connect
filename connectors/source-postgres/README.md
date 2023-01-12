---
title: PostgreSQL CDC (Debezium)
---

# PostgreSQL CDC Source (Debezium)

## Introduction

The PostgreSQL Source is a [Vance Connector][vc] which use [Debezium][debezium] obtain a snapshot of the existing data
in a PostgreSQL schema and then monitor and record all subsequent row-level changes to that data.

For example, PostgreSQL database vance_test with schema public has table user Look:

```text
Column      |          Type          | Collation | Nullable | Default
------------+------------------------+-----------+----------+---------
 id         | character varying(100) |           | not null | 
 first_name | character varying(100) |           | not null | 
 last_name  | character varying(100) |           | not null | 
 email      | character varying(100) |           | not null | 
 
```

The row record will be transformed into a CloudEvent looks like:

```json
{
  "specversion": "1.0",
  "id": "e5f19d0a-8120-41a2-b4a3-ad3de6c66f6c",
  "source": "/vanus/debezium/postgresql/quick_start",
  "type": "vanus.debezium.postgresql.datachangeevent",
  "datacontenttype": "application/json",
  "time": "2023-01-11T03:23:20.973Z",
  "xvdebeziumname": "quick_start",
  "xvdebeziumop": "r",
  "xvop": "c",
  "xvdb": "vance_test",
  "xvschema": "public",
  "xvtable": "user",
  "data": {
    "id": "1",
    "first_name": "Anne",
    "last_name": "Kretchmar",
    "email": "annek@noanswer.org"
  }
}

```

## PostgreSQL Source Configs

Users can specify their configs by either setting environments variables or mount a config.json to
`/vance/config/config.json` when they run the connector. Find examples of setting configs [here][config].

### Config Fields of the PostgreSQL Source

| name             | requirement | description                                                                                            |
|------------------|-------------|--------------------------------------------------------------------------------------------------------|
| target           | required    | target URL will send CloudEvents to                                                                    |
| name             | required    | unique name for the connector                                                                          |
| db.host          | required    | IP address or host name of db                                                                          |
| db.port          | required    | integer port number of db                                                                              |
| db.username      | required    | username of db                                                                                         |
| db.password      | required    | password of db                                                                                         |
| db.database      | required    | database of db                                                                                         |
| schema_include   | optional    | schema name which want to capture changes, string array, can not set with schema_exclude               |
| schema_exclude   | optional    | schema name which don't want to capture changes,string array, can not set with schema_include          |
| table_include    | optional    | table name which want to capture changes, string array and format is schema.tableName                  |
| table_exclude    | optional    | table name which don't want to capture changes, string array and format is schema.tableName            |
| store.type       | required    | save offset type, support FILE, MEMORY                                                                 |
| store.pathname   | required    | it's needed when offset type is FIlE, save offset file name                                            |
| plugin_name      | optional    | The name of the [logical decoding plug-in] installed on the PostgreSQL server,default pgoutput         |
| slot_name        | optional    | The name of the logical decoding slot that was created for streaming changes from a particular plug-in |
| publication_name | optional    | The name of the publication created for streaming changes when using pgoutput                          |
| offset.lsn       | optional    | PostgreSQL Log Sequence Numbers which begin to capture the change,such as "0/17EFB50"                  |

### Config Example

```yaml
target: "http://localhost:8080"
name: "quick_start"
db:
  host: "localhost"
  port: 5432
  username: "vance_test"
  password: "123456"
  database: "vance_test"
schema_include: [ "public" ]
table_include: [ "public.user" ]

slot_name: vanus_slot
publication_name: vanus_publication
```

## PostgreSQL Source Image

> public.ecr.aws/vanus/connector/source-postgres

### Running with Docker

```shell
docker run --rm -v ${PWD}:/vance/config public.ecr.aws/vanus/connector/source-postgres
```

### K8S

```shell
  kubectl apply -f source-postgres.yaml
```

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md

[config]: https://github.com/linkall-labs/vance-docs/blob/main/docs/connector.md

[debezium]: https://debezium.io/documentation/reference/2.1/connectors/postgresql.html

[logical decoding plug-in]: https://debezium.io/documentation/reference/2.1/connectors/postgresql.html#postgresql-output-plugin
