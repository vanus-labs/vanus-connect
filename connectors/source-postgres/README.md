---
title: PostgreSQL CDC (Debezium)
---

# PostgreSQL CDC Source (Debezium)

## Introduction

The PostgreSQL Source is a [Vance Connector][vc] which use [Debezium][debezium] obtain a snapshot of the existing data
in a PostgreSQL schema and then monitor and record all subsequent row-level changes to that data.

For example,PostgreSQL dababase vance_test with schema public has table customers Look:

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
  "id": "name:quick_start;lsn:25090824;txId:826",
  "source": "/debezium/postgresql/quick_start",
  "specversion": "1.0",
  "type": "io.debezium.postgresql.datachangeevent",
  "time": "2023-01-04T13:00:58.360Z",
  "datacontenttype": "application/json",
  "iodebeziumop": "r",
  "iodebeziumversion": "2.1.1.Final",
  "iodebeziumconnector": "postgresql",
  "iodebeziumname": "quick_start",
  "iodebeziumtsms": "1672837258360",
  "iodebeziumsnapshot": "last",
  "iodebeziumdb": "vance_test",
  "iodebeziumsequence": "[null,\"25090824\"]",
  "iodebeziumschema": "public",
  "iodebeziumtable": "user",
  "iodebeziumtxid": "826",
  "iodebeziumlsn": "25090824",
  "iodebeziumxmin": null,
  "iodebeziumtxid": null,
  "iodebeziumtxtotalorder": null,
  "iodebeziumtxdatacollectionorder": null,
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

| name                | requirement | description                                                                                                               |
|---------------------|-------------|---------------------------------------------------------------------------------------------------------------------------|
| host                | required    | db host                                                                                                                   |
| port                | required    | db port                                                                                                                   |
| username            | required    | db username                                                                                                               |
| password            | required    | db password                                                                                                               |
| db_name             | required    | db database name                                                                                                          |
| schema_name         | required    | the name of schema want to capture changes,default public                                                                 |
| include_table       | optional    | the name of table want to capture changes, many split by comma                                                            |
| plugin_name         | optional    | The name of the [logical decoding plug-in] installed on the PostgreSQL server,default pgoutput                            |
| slot_name           | optional    | The name of the logical decoding slot that was created for streaming changes from a particular plug-in,default vance_slot |
| publication_name    | optional    | The name of the publication created for streaming changes when using pgoutput,default vance_publication                   |
| v_target            | required    | target URL will send CloudEvents to                                                                                       |
| v_store_file        | required    | save offset file name                                                                                                     |
| store_offset_key    | optional    | offset store use key, default is vance_debezium_offset                                                                    |
| offset_lsn          | optional    | PostgreSQL Log Sequence Numbers which begin to capture the change                                                         |

### Config Example

```json
{
  "host": "localhost",
  "port": "5432",
  "username": "postgres",
  "password": "123456",
  "db_name": "dbname",
  "include_table": "user",
  "v_store_file": "/vance/data/offset.dat",
  "v_target": "http://localhost:8080"
}
```

## PostgreSQL Source Image

> docker.io/vancehub/source-postgres

### Running with Docker

```shell
docker run -v $(pwd)/config.json:/vance/config/config.json -v $(pwd)/data:/vance/data --rm vancehub/source-postgres
```

## Local Development

You can run the source codes of the PostgreSQL Source locally as well.

### Building via Maven

```shell
cd source-postgres 
mvn clean package
```

### Running via Maven

```shell
mvn exec:java -Dexec.mainClass="com.linkall.source.postgresql.Entrance"
```

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md
[config]: https://github.com/linkall-labs/vance-docs/blob/main/docs/connector.md
[debezium]: https://debezium.io/documentation/reference/1.9/connectors/postgresql.html
[logical decoding plug-in]: https://debezium.io/documentation/reference/1.9/connectors/postgresql.html#postgresql-output-plugin
