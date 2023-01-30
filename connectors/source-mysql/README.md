---
title: MySQL CDC（Debezium)
---

# MySQL CDC Source (Debezium)

## Introduction

The MySQL Source is a [Vanus Connector][vc] which use [Debezium][debezium] obtain a snapshot of the existing data in a
MySql database and then monitor and record all subsequent row-level changes to that data.

For example, MySQL database vanus_test has table user Look:

```text
+----------+-----------+-----------------+
| id       | name      | email           |
+----------+-----------+-----------------+
| 100      | vanus     | dev@example.com |
+----------+-----------+-----------------+
```

The row record will be transformed into a CloudEvent looks like:

```json
{
  "specversion": "1.0",
  "id": "88767821-92c2-477d-9a6f-bfdfbed19c6a",
  "source": "/debezium/mysql/quick_start",
  "type": "debezium.mysql.datachangeevent",
  "datacontenttype": "application/json",
  "time": "2023-01-11T07:25:39.557Z",
  "xvdebeziumname": "quick_start",
  "xvdebeziumop": "r",
  "xvop": "c",
  "xvdb": "vanus_test",
  "xvtable": "user",
  "data": {
    "id": 100,
    "name": "vanus",
    "email": "dev@example.com"
  }
}
```

## Quick Start

This section shows how MySQL Source convert db record to a CloudEvent.

### Prerequisites
- Have a container runtime (i.e., docker).
- Have a MySQL Server `8.0`, `5.7`, or `5.6`.

#### Setting up MySQL

1. Enable binary logging.
         
   You must enable binary logging for MySQL replication. The binary logs record transaction updates for replication tools to propagate changes.
You can configure your MySQL server configuration file with the following properties, which are described in below:

    ```text
    server-id         = 223344
    log_bin           = mysql-bin
    binlog_format     = ROW
    binlog_row_image  = FULL
    expire_logs_days  = 10
    ```
    
    See the [MySQL doc](https://dev.mysql.com/doc/refman/8.0/en/replication-options-binary-log.html) for more details;

2. Enable GTIDs (Optional).

   GTIDs are available in MySQL 5.6.5 and later. See the [MySQL doc](https://dev.mysql.com/doc/refman/8.0/en/replication-options-gtids.html#option_mysqld_gtid-mode) for more details.
   1. Enable gtid_mode
      ```sql
       mysql> gtid_mode=ON;
      ```
   2. Enable enforce_gtid_consistency
      ```sql
       mysql> enforce_gtid_consistency=ON;
      ```

#### Prepare data

1. Create database and table 
   ```sql
   create database vanus_test;
   CREATE TABLE IF NOT EXISTS vanus_test.user
   (
     `id` int NOT NULL,
     `name` varchar(100) NOT NULL,
     `email` varchar(100) NOT NULL,
     PRIMARY KEY (`id`)
   ) ENGINE=InnoDB;
   ```
2. Insert data
   ```sql
   insert into vanus_test.`user` values(100,"vanus","dev@example.com");
   ```
3. Create user and grant role
   ```sql
   CREATE USER 'vanus_test'@'%' IDENTIFIED WITH mysql_native_password BY '123456';
   GRANT SELECT, RELOAD, SHOW DATABASES, REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO 'vanus_test'@'%';
   ```

### Create the config file

```shell
cat << EOF > config.yml
target: http://localhost:31081
name: "quick_start"
db:
  host: "localhost"
  port: 3306
  username: "vanus_test"
  password: "123456"
database_include: [ "vanus_test" ]
# format is vanus_test.tableName
table_include: [ "vanus_test.user" ]

store:
  type: FILE
  pathname: "/vanus-connect/data/offset.dat"

db_history_file: "/vanus-connect/data/history.dat"
EOF
```

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

The MySQL Source tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  -v ${PWD}:/vanus-connect/data \
  --name source-mysql public.ecr.aws/vanus/connector/source-mysql
```

### Test

Open a terminal and use the following command to run a Display sink, which receives and prints CloudEvents.

```shell
docker run -it --rm \
  -p 31081:8080 \
  --name sink-display public.ecr.aws/vanus/connector/sink-display
```

Make sure the `target` value in your config file is `http://localhost:31081` so that the Source can send CloudEvents to our Display Sink.
 

Here is the sort of CloudEvent you should expect to receive in the Display Sink:

```json
{
  "specversion": "1.0",
  "id": "88767821-92c2-477d-9a6f-bfdfbed19c6a",
  "source": "/debezium/mysql/quick_start",
  "type": "debezium.mysql.datachangeevent",
  "datacontenttype": "application/json",
  "time": "2023-01-11T07:25:39.557Z",
  "xvdebeziumname": "quick_start",
  "xvdebeziumop": "r",
  "xvop": "c",
  "xvdb": "vanus_test",
  "xvtable": "user",
  "data": {
    "id": 100,
    "name": "vanus",
    "email": "dev@example.com"
  }
}
```

### Clean

```shell
docker stop source-mysql sink-display
```

## Run in Kubernetes

```shell
  kubectl apply -f source-mysql.yaml
```

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: source-mysql
  namespace: vanus
data:
  config.yml: |-
    target: "http://vanus-gateway.vanus:8080/gateway/quick_start"
    name: "quick_start"
    db:
      host: "localhost"
      port: 3306
      username: "root"
      password: "123456"
    database_include: [ "vanus_test" ]
    # format is vanus_test.tableName
    table_include: [ "vanus_test.user" ]

    store:
      type: FILE
      pathname: "/vanus-connect/data/offset.dat"

    db_history_file: "/vanus-connect/data/history.dat"

---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: source-mysql
  namespace: vanus
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: source-mysql
  namespace: vanus
  labels:
    app: source-mysql
spec:
  selector:
    matchLabels:
      app: source-mysql
  replicas: 1
  template:
    metadata:
      labels:
        app: source-mysql
    spec:
      containers:
        - name: source-mysql
          image: public.ecr.aws/vanus/connector/source-mysql
          imagePullPolicy: Always
          volumeMounts:
            - name: config
              mountPath: /vanus-connect/config
            - name: data
              mountPath: /vanus-connect/data
      volumes:
        - name: config
          configMap:
            name: source-mysql
        - name: data
          persistentVolumeClaim:
            claimName: source-mysql
```

## Integrate with Vanus

This section shows how a source connector can send CloudEvents to a running [Vanus cluster](https://github.com/linkall-labs/vanus).

### Prerequisites
- Have a running K8s cluster
- Have a running Vanus cluster
- Vsctl Installed

1. Export the VANUS_GATEWAY environment variable (the ip should be a host-accessible address of the vanus-gateway service)
```shell
export VANUS_GATEWAY=192.168.49.2:30001
```

2. Create an eventbus
```shell
vsctl eventbus create --name quick-start
```

3. Update the target config of the MySQL Source
```yaml
target: http://192.168.49.2:30001/gateway/quick-start
```

4. Run the MySQL Source
```shell
  kubectl apply -f source-mysql.yaml
```

[vc]: https://www.vanus.dev/introduction/concepts#vanus-connect
[debezium]: https://debezium.io/documentation/reference/2.0/connectors/mysql.html
