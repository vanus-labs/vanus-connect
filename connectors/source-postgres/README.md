---
title: PostgreSQL CDC (Debezium)
---

# PostgreSQL CDC Source (Debezium)

## Introduction

The PostgreSQL Source is a [Vanus Connector][vc] which use [Debezium][debezium] obtain a snapshot of the existing data
in a PostgreSQL schema and then monitor and record all subsequent row-level changes to that data.

For example, a PostgreSQL database schema look like this:

```text
Column      |          Type          | Collation | Nullable | Default
------------+------------------------+-----------+----------+---------
 id         | character varying(100) |           | not null | 
 first_name | character varying(100) |           | not null | 
 last_name  | character varying(100) |           | not null | 
 email      | character varying(100) |           | not null | 
 
```

The row record will be transformed into a CloudEvent in the following way:

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
  "xvdb": "vanus_test",
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

## Quick Start

This section will teach you how to use PostgreSQL Source to convert db record to a CloudEvent.

### Prerequisites
- Have a container runtime (i.e., docker).
- Have a PostgreSQL database。

#### Setting up Postgres

1. Enable logical replication.
  configure the following parameters in the [postgresql.conf](https://www.postgresql.org/docs/current/config-setting.html) file
  ```text
    wal_level = logical             # type of coding used within the Postgres write-ahead log.minimal, archive, hot_standby, or logical (change requires restart)
    max_wal_senders = 1             # the maximum number of processes used for handling WAL changes (change requires restart)
    max_replication_slots = 1       # the maximum number of replication slots that are allowed to stream WAL changes (change requires restart)
  ```
2. Select a replication plugin.
   We recommend using a [pgoutput](https://www.postgresql.org/docs/9.6/logicaldecoding-output-plugin.html) plugin (the standard logical decoding plugin in Postgres).
   The PostgreSQl Source support logical decoding plugins from Debezium:
   - [protobuf](https://github.com/debezium/postgres-decoderbufs/blob/main/README.md) : To encode changes in Protobuf format
   - [wal2json](https://github.com/eulerto/wal2json/blob/master/README.md) : To encode changes in JSON format

#### Prepare data

1. Connect to PostgreSQL and create database `vanus_test`.
2. Create table use following command:
    ```shell
    CREATE TABLE IF NOT EXISTS public."user"
    (
        id character varying(100) NOT NULL,
        first_name character varying(100) NOT NULL,
        last_name character varying(100) NOT NULL,
        email character varying(100) NOT NULL,
        CONSTRAINT user_pkey PRIMARY KEY (id)
    );
    ```
3. Insert data
    ```sql
    insert into public."user"(id,first_name,last_name,email) values(1,'Anne','Kretchmar','annek@noanswer.org');
    ```
4. Create user and grant role
    ```sql
    CREATE USER vanus_test WITH PASSWORD '123456' REPLICATION LOGIN;
    GRANT SELECT ON TABLE "user" to vanus_test;
    ```
5. Create replication slot using pgoutput, run:
  ```sql
    SELECT pg_create_logical_replication_slot('vanus_slot', 'pgoutput');
  ```
6. Create publications and replication identities for tables
    ```sql
    CREATE PUBLICATION vanus_publication FOR TABLE "user";
    ```

Refer to the [Postgres docs](https://www.postgresql.org/docs/10/sql-alterpublication.html)

### Create the config file

```shell
cat << EOF > config.yml
target: http://localhost:31081
name: "quick_start"
db:
  host: "localhost"
  port: 5432
  username: "vanus_test"
  password: "123456"
  database: "vanus_test"
schema_include: [ "public" ]
table_include: [ "public.user" ]
plugin_name: pgoutput
slot_name: vanus_slot
publication_name: vanus_publication

EOF
```

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

The PostgreSQL Source tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-postgres public.ecr.aws/vanus/connector/source-postgres
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
  "id": "e5f19d0a-8120-41a2-b4a3-ad3de6c66f6c",
  "source": "/vanus/debezium/postgresql/quick_start",
  "type": "vanus.debezium.postgresql.datachangeevent",
  "datacontenttype": "application/json",
  "time": "2023-01-11T03:23:20.973Z",
  "xvdebeziumname": "quick_start",
  "xvdebeziumop": "r",
  "xvop": "c",
  "xvdb": "vanus_test",
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

### Clean

```shell
docker stop source-postgres sink-display
```

### K8S

```shell
  kubectl apply -f source-postgres.yaml
```

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: source-postgres
  namespace: vanus
data:
  config.yaml: |-
    target: "http://vanus-gateway.vanus:8080/gateway/quick_start"
    name: "quick_start"
    db:
      host: "localhost"
      port: 5432
      username: "vanus_test"
      password: "123456"
      database: "vanus_test"
    schema_include: [ "public" ]
    table_include: [ "public.user" ]

    plugin_name: pgoutput
    slot_name: vanus_slot
    publication_name: vanus_publication
    
    store:
        type: FILE
        pathname: "/vanus-connect/data/offset.dat"
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: source-postgres
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
  name: source-postgres
  namespace: vanus
  labels:
    app: source-postgres
spec:
  selector:
    matchLabels:
      app: source-postgres
  replicas: 1
  template:
    metadata:
      labels:
        app: source-postgres
    spec:
      containers:
        - name: source-postgres
          image: public.ecr.aws/vanus/connector/source-postgres
          imagePullPolicy: Always
          volumeMounts:
            - name: config
              mountPath: /vanus-connect/config
            - name: data
              mountPath: /vanus-connect/data
      volumes:
        - name: config
          configMap:
            name: source-postgres
        - name: data
          persistentVolumeClaim:
            claimName: source-postgres
```

## Integrate with Vanus

This section shows how a source connector can send CloudEvents to a running [Vanus cluster](https://github.com/vanus-labs/vanus).

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

3. Update the target config of the PostgreSQL Source
```yaml
target: http://192.168.49.2:30001/gateway/quick-start
```

4. Run the PostgreSQL Source
```shell
  kubectl apply -f source-postgres.yaml
```

[vc]: https://docs.vanus.ai/introduction/concepts#vanus-connect
[debezium]: https://debezium.io/documentation/reference/2.1/connectors/postgresql.html
[logical decoding plug-in]: https://debezium.io/documentation/reference/2.1/connectors/postgresql.html#postgresql-output-plugin
