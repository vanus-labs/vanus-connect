---
title: MySql
---

# MySql Sink

## Introduction

The MySql Sink is a [Vance Connector][vc] which aims to handle incoming CloudEvents in a way that extracts the data part of the original event 
and deliver these extracted data to database mysql by use jdbc. You need create table before start the sink because table will not auto create .

For example, if the incoming CloudEvent looks like:

```json
{
  "id" : "vance.vance_test:binlog.000010:2515",
  "source" : "vance.debezium.mysql",
  "specversion" : "1.0",
  "type" : "debezium.mysql.vance.vance_test",
  "time" : "2022-07-08T03:17:03.139Z",
  "datacontenttype" : "application/json",
  "vancedebeziumop" : "r",
  "vancedebeziumversion" : "1.9.4.Final",
  "vancedebeziumconnector" : "mysql",
  "vancedebeziumname" : "vance",
  "vancedebeziumtsms" : "1657250223138",
  "vancedebeziumsnapshot" : "true",
  "vancedebeziumdb" : "vance",
  "vancedebeziumtable" : "vance_test",
  "vancedebeziumpos" : "2515",
  "vancedebeziumfile": "binlog.000010",
  "vancedebeziumrow": "0",
  "data" : {
    "id":18,
    "name":"xdl",
    "description":"Development Manager",
    "date": "2022-07-06"
  }
}
```

The MySql Sink will extract data field write to database table like:

```text
+----+---------+---------------------+------------+
| id | name    | description         | date       |
+----+---------+---------------------+------------+
| 18 | xdl     | Development Manager | 2022-07-06 |
+----+---------+---------------------+------------+
```

## MySql Sink Configs

Users can specify their configs by either setting environments variables or mount a config.json to
`/vance/config/config.json` when they run the connector. Find examples of setting configs [here][config].

### Config Fields of the MySql Sink

| name             | requirement | description                                                                |
|------------------|-------------|----------------------------------------------------------------------------|
| v_port           | optional    | v_port is used to specify the port MySql Sink is listening on,default 8080 |
| table_name       | required    | db table name                                                              |
| insert_mode      | optional    | insert mode: insert or upsert, default insert                              |
| commit_interval  | optional    | batch data commit to db interval, unit is millisecond default 1000         |

## MySql Sink Secrets

Users should set their sensitive data Base64 encoded in a secret file.
And mount your local secret file to `/vance/secret/secret.json` when you run the connector.

### Encode your sensitive data

```shell
$ echo -n ABCDEFG | base64
QUJDREVGRw==
```

Replace 'ABCDEFG' with your sensitive data.

### Secret Fields of the Mysql Sink

| name               | requirement | description        |
|--------------------|-------------|--------------------|
| host               | required    | db host            |
| port               | required    | db port            |
| username           | required    | db username        |
| password           | required    | db password        |
| dbName             | required    | db database name   |

## MySql Sink Image

> docker.io/vancehub/sink-mysql

## Local Development

You can run the sink codes of the MySql Sink locally as well.

### Building via Maven

```shell
$ cd sink-mysql 
$ mvn clean package
```

### Running via Maven

```shell
$ mvn exec:java -Dexec.mainClass="com.linkall.sink.mysql.Entrance"
```

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md
[config]: https://github.com/linkall-labs/vance-docs/blob/main/docs/connector.md
