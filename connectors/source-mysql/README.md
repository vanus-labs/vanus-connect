# MySql Source

## Introduction

The MySql Source is a [Vance Connector][vc] which use [Debezium][debezium] obtain a snapshot of the existing data in a
MySql database and then monitor and record all subsequent row-level changes to that data.

For example,MySql database vance has table vance_test Look:

```text
+-------------+--------------+------+-----+---------+----------------+
| Field       | Type         | Null | Key | Default | Extra          |
+-------------+--------------+------+-----+---------+----------------+
| id          | int          | NO   | PRI | NULL    | auto_increment |
| name        | varchar(100) | NO   |     | NULL    |                |
| description | varchar(100) | NO   |     | NULL    |                |
| date        | date         | YES  |     | NULL    |                |
+-------------+--------------+------+-----+---------+----------------+
```

The row record will be transformed into a CloudEvent looks like:

```json
{
  "id": "vance.vance_test:binlog.000010:2515",
  "source": "vance.debezium.mysql",
  "specversion": "1.0",
  "type": "debezium.mysql.vance.vance_test",
  "time": "2022-07-08T03:17:03.139Z",
  "datacontenttype": "application/json",
  "vancedebeziumop": "r",
  "vancedebeziumversion": "1.9.4.Final",
  "vancedebeziumconnector": "mysql",
  "vancedebeziumname": "vance",
  "vancedebeziumtsms": "1657250223138",
  "vancedebeziumsnapshot": "true",
  "vancedebeziumdb": "vance",
  "vancedebeziumtable": "vance_test",
  "vancedebeziumpos": "2515",
  "vancedebeziumfile": "binlog.000010",
  "vancedebeziumrow": "0",
  "data": {
    "id": 18,
    "name": "xdl",
    "description": "Development Manager",
    "date": "2022-07-06"
  }
}
```

## MySql Source Configs

Users can specify their configs by either setting environments variables or mount a config.json to
`/vance/config/config.json` when they run the connector. Find examples of setting configs [here][config].

### Config Fields of the Mysql Source

| name                 | requirement | description                                                                    |
|----------------------|-------------|--------------------------------------------------------------------------------|
| host                 | required    | db host                                                                        |
| port                 | required    | db port                                                                        |
| username             | optional    | db username                                                                    |
| password             | optional    | db password                                                                    |
| db_name              | required    | db database name                                                               |
| v_target             | required    | target URL will send CloudEvents to                                            |
| v_store_file         | required    | save offset file name                                                          |
| db_history_file      | required    | save db schema history file name                                               |
| store_offset_key     | optional    | offset store use key, default is vance_debezium_offset                         |
| offset_binlog_file   | optional    | binlog filename, increment sync start binlog file name if not set is full sync |
| offset_binlog_pos    | optional    | binlog position, use with config offset_binlog_file                            |
| offset_binlog_gtids  | optional    | binlog grids                                                                   |

### Config Example

```json
{
  "host": "localhost",
  "port": "3306",
  "username": "root",
  "password": "123456",
  "db_name": "dbname",
  "include_table": "user",

  "v_store_file": "/vance/data/offset.dat",
  "db_history_file": "/vance/data/history.dat",

  "v_target": "http://localhost:8080"
}
```

## MySql Source Image

> docker.io/vancehub/source-mysql

### Running with Docker

> docker run -v $(pwd)/config.json:/vance/config/config.json -v $(pwd)/data:/vance/data --rm vancehub/source-mysql

## Local Development

You can run the source codes of the MySql Source locally as well.

### Building via Maven

```shell
cd source-mysql 
mvn clean package
```

### Running via Maven

```shell
mvn exec:java -Dexec.mainClass="com.linkall.source.mysql.Entrance"
```

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md

[config]: https://github.com/linkall-labs/vance-docs/blob/main/docs/connector.md

[debezium]: https://debezium.io/documentation/reference/1.9/connectors/mysql.html
