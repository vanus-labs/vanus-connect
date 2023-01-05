---
title: MySQL
---

# MySQL Sink
This document provides a brief introduction to the MySQL Sink.
It is also designed to guide you through the process of running a
MySQL Sink Connector.

## Introduction
The MySQL Sink is a [Vance Connector][vc] that aims to handle incoming CloudEvents
in a way that extracts the data part of the original event and delivers these
extracted data to a MySQL database using JDBC. Before using this Sink, you will
need to create a database and a table.

## Handling incoming CloudEvent
For example, if the incoming CloudEvent looks like this:
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
###
The MySQL Sink will extract the data fields and write to the database table in the following way:
```text
+----+---------+---------------------+------------+
| id | name    | description         | date       |
+----+---------+---------------------+------------+
| 18 | xdl     | Development Manager | 2022-07-06 |
+----+---------+---------------------+------------+
```
---
## Quick Start
This quick start will guide you through the process of running an MySQL Sink Connector.

### Prerequisites
- Have container runtime (i.e., docker).
- Have running [MySQL][mysql] database.
- Have a database and table created.

### Set MySQL Sink Configurations
You can specify your configs by either setting environments
variables or mounting a config.json to `/vance/config/config.json`
when running the Connector.

Here is an example of a configuration file for the MySQL Sink.
```json
{
  "v_port": "8081",
  "table_name": "Costumer",
  "insert_mode": "insert",
  "commit_interval": "100" 
}
```

#### Config Fields of the MySQL Sink
| name             | requirement | description                                                                |
|------------------|-------------|----------------------------------------------------------------------------|
| v_port           | optional    | v_port is used to specify the port MySql Sink is listening on,default 8080 |
| table_name       | required    | db table name                                                              |
| insert_mode      | optional    | insert mode: insert or upsert, default insert                              |
| commit_interval  | optional    | batch data commit to db interval, unit is millisecond default 1000         |

### MySQL Sink Secrets
Users should set their sensitive data Base64 encoded in a secret file.
And mount your local secret file to `/vance/secret/secret.json` when you run the Connector.

#### Encode your Sensitive Data
Replace MY_SECRET with your sensitive data to get the Base64-based string.

```shell
$ echo -n MY_SECRET | base64
QUJDREVGRw==
```

Here is an example of a secret file for the MySQL Sink.
```json
{
  "host": "TVlfU0VDUkVUTVlfU0VDUkVU",
  "port": "OTA4Mw==",
  "username": "bG92ZWNob2NvbGF0ZQ==",
  "password": "MTIzNDU2Nzg5",
  "dbName": "SW1XYWxraW5PblN1blNoaW5l"
}
```

#### Secret Fields of the MySQL Sink

| name               | requirement | description        |
|--------------------|-------------|--------------------|
| host               | required    | db host            |
| port               | required    | db port            |
| username           | required    | db username        |
| password           | required    | db password        |
| dbName             | required    | db database name   |

### Run the MySQL Sink with Docker
Create your config.json and secret.json, and mount them to
specific paths to run the MySQL Sink using the following command.

> docker run -v $(pwd)/secret.json:/vance/secret/secret.json -v $(pwd)/config.json:/vance/config/config.json --rm vancehub/sink-mysql

### Verify the MySQL Sink
You can verify if the MySQL Sink works properly by
sending a CloudEvent with the POST terminal command, for example.

```shell
curl -X POST -d 
"{
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
}"
http://localhost:8081 
```
:::tip
Note that the last line contains the address and port targeted.
:::


[vc]: https://docs.linkall.com/concepts/connector
[mysql]: https://www.mysql.com