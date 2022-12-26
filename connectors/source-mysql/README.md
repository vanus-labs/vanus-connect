---
title: MySql
---

# MySql Source

## Introduction

The MySql Source is a [Vance Connector][vc] which use [Debezium][debezium] obtain a snapshot of the existing data in a
MySql database and then monitor and record all subsequent row-level changes to that data.

For example,MySql database dbname has table user Look:

```text
+-------------+--------------+------+-----+---------+----------------+
| Field       | Type         | Null | Key | Default | Extra          |
+-------------+--------------+------+-----+---------+----------------+
| id          | int          | NO   | PRI | NULL    | auto_increment |
| name        | varchar(100) | NO   |     | NULL    |                |
| email       | varchar(100) | NO   |     | NULL    |                |
+-------------+--------------+------+-----+---------+----------------+
```

The row record will be transformed into a CloudEvent looks like:

```json
{
  "specversion": "1.0",
  "id": "a67f31d6-a0c2-4124-b794-4139a9525ea8",
  "source": "/debezium/mysql/quick_start",
  "type": "io.debezium.mysql.datachangeevent",
  "datacontenttype": "application/json",
  "time": "2022-12-23T16:06:14Z",
  "iodebeziumconnector": "mysql",
  "iodebeziumserverid": "0",
  "iodebeziumsnapshot": "last",
  "iodebeziumdb": "dbname",
  "iodebeziumfile": "binlog.000009",
  "iodebeziumpos": "197",
  "iodebeziumname": "quick_start",
  "iodebeziumtsms": "1671811574000",
  "iodebeziumtable": "user",
  "iodebeziumop": "r",
  "iodebeziumversion": "2.0.1.Final",
  "iodebeziumrow": "0",
  "data": {
      "id": 100,
      "name": "user_name",
      "email": "user_email"
  }
}
```

## MySql Source Configs

### Config

| name                | requirement | description                                                                                       |
|---------------------|-------------|---------------------------------------------------------------------------------------------------|
| target              | required    | target URL will send CloudEvents to                                                               |
| name                | required    | unique name for the connector                                                                     |
| db.host             | required    | IP address or host name of db                                                                     |
| db.port             | required    | integer port number of db                                                                         |
| db.username         | required    | username of db                                                                                    |
| db.password         | required    | password of db                                                                                    |
| database_include    | optional    | database name which want to capture changes, string array, can not set with exclude_database      |
| database_exclude    | optional    | database name which don't want to capture changes,string array, can not set with include_database |
| table_include       | optional    | table name which want to capture changes, string array and format is databaseName.tableName       |
| table_exclude       | optional    | table name which don't want to capture changes, string array and format is databaseName.tableName |
| store.type          | required    | save offset type, support FILE, MEMORY                                                            |
| store.pathname      | required    | it's needed when offset type is FIlE, save offset file name                                       |
| db_history_file     | required    | save db schema history file name                                                                  |
| binlog_offset.file  | optional    | binlog filename, increment sync start binlog file name if not set is full sync                    |
| binlog_offset.pos   | optional    | binlog position, use with config offset_binlog_file                                               |
| binlog_offset.gtids | optional    | binlog grids                                                                                      |

### Config Example

```yaml
target: "http://localhost:8080"
name: "quick_start"
db:
  host: "localhost"
  port: 3306
  username: "root"
  password: "vanus123456"
database_include:
  - dbname
table_include:
  - dbname.user

store:
  type: FILE
  pathname: "/tmp/mysql/offset.dat"

db_history_file: "/tmp/mysql/history.dat"
```

## MySql Source Image

> public.ecr.aws/vanus/connector/source-mysql

### Running with Docker

```shell
docker run --rm -v ${PWD}:/vance/config public.ecr.aws/vanus/connector/source-mysql
```

### K8S

```shell
  kubectl apply -f source-mysql.yaml
```

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md
[debezium]: https://debezium.io/documentation/reference/2.0/connectors/mysql.html
