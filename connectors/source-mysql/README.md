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
  "id": "ad3eaae9-2d6d-46e7-88d3-4b657155f183",
  "source": "/debezium/mysql/dbname",
  "type": "io.debezium.mysql.datachangeevent",
  "datacontenttype": "application/json",
  "time": "2022-12-23T05:54:49Z",
  "iodebeziumconnector": "mysql",
  "iodebeziumserverid": "0",
  "iodebeziumsnapshot": "last",
  "iodebeziumdb": "dbname",
  "iodebeziumfile": "binlog.000009",
  "iodebeziumpos": "197",
  "iodebeziumname": "dbname",
  "iodebeziumtsms": "1671774889000",
  "iodebeziumtable": "user",
  "iodebeziumop": "r",
  "iodebeziumversion": "2.0.1.Final",
  "iodebeziumrow": "0",
  "data": {
    "after": {
      "id": 1,
      "name": "user_name",
      "email": "user_email"
    }
  }
}
```

## MySql Source Configs

### Config

| name                    | requirement | description                                                                    |
|-------------------------|-------------|--------------------------------------------------------------------------------|
| target                  | required    | target URL will send CloudEvents to                                            |
| db_config.host          | required    | db host                                                                        |
| db_config.port          | required    | db port                                                                        |
| db_config.username      | optional    | db username                                                                    |
| db_config.password      | optional    | db password                                                                    |
| db_config.database      | required    | db database name                                                               |
| include_tables          | required    | include table                                                                  |
| exclude_tables          | required    | exclude table                                                                  |
| store_config.type       | required    | save offset type, support FILE,MEMORY                                          |
| store_config.store_file | required    | it's needed when offset type is FIlE, save offset file name                    |
| db_history_file         | required    | save db schema history file name                                               |
| binlog_offset.file      | optional    | binlog filename, increment sync start binlog file name if not set is full sync |
| binlog_offset.pos       | optional    | binlog position, use with config offset_binlog_file                            |
| binlog_offset.gtids     | optional    | binlog grids                                                                   |

### Config Example

```yaml
target: "http://localhost:8080"
db_config:
  host: "localhost"
  port: 3306
  username: "root"
  password: "vanus123456"
  database: "dbname"
include_tables:
  - user

store_config:
  type: FILE
  store_file: "/tmp/mysql/offset.dat"

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
