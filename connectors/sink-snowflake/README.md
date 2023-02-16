---
title: Snowflake
---

# Snowflake Sink

## Introduction

The Snowflake Sink is a [Vanus Connector][vc] that aims to handle incoming CloudEvents in a way that extracts the data
part of the original event and delivers these extracted data to a Snowflake database using [bulk loading](loadfile).

For example, if the incoming CloudEvent looks like this:

```json
{
  "id": "88767821-92c2-477d-9a6f-bfdfbed19c6a",
  "source": "quickstart",
  "specversion": "1.0",
  "type": "quickstart",
  "time": "2022-07-08T03:17:03.139Z",
  "datacontenttype": "application/json",
  "data": {
    "id": 18,
    "name": "xdl",
    "email": "Development Manager",
    "date": "2022-07-06"
  }
}
```

The Snowflake Sink will extract the data fields and write to the database table in the following way:

```text
+----+---------+---------------------+------------+
| id | name    | description         | date       |
+----+---------+---------------------+------------+
| 18 | xdl     | Development Manager | 2022-07-06 |
+----+---------+---------------------+------------+
```

## Quick Start

This quick start will guide you through the process of running a Snowflake Sink Connector.

### Prerequisites

- Have a container runtime (i.e., docker).
- Have a running [snowflake][snowflake] database.

### Create the config file

```shell
cat << EOF > config.yml
snowflake:
  host: "myaccount.ap-northeast-1.aws.snowflakecomputing.com"
  username: "vanus_user" 
  password: "snowflake"
  role: "ACCOUNTADMIN"
  warehouse: "xxxxxx"
  database: "VANUS_DB"
  schema: "public"
  table: "vanus_test"

EOF
```

| Name                 | Required |    Default    | Description                                                                          |
|:---------------------|:--------:|:-------------:|--------------------------------------------------------------------------------------|
| port                 |    NO    |     8080      | the port which Snowflake Sink listens on                                             |
| snowflake.host       |   YES    |               | [account] of snowflake, example: myaccount.ap-northeast-1.aws.snowflakecomputing.com |
| snowflake.username   |   YES    |               | username of snowflake                                                                |
| snowflake.password   |   YES    |               | password of snowflake                                                                |
| snowflake.role       |   YES    |               | [role] of snowflake                                                                  |
| snowflake.warehouse  |   YES    |               | [warehouse] of snowflake                                                             |
| snowflake.database   |   YES    |               | [database] of snowflake                                                              |
| snowflake.schema     |   YES    |               | [schema](database) of snowflake                                                      |
| snowflake.table      |   YES    |               | table name of snowflake, the table no need exist                                     |
| snowflake.properties |    NO    |               | the other properties for jdbc [jdbc parameters](jdbc-parameter) of snowflake         |
| flush_time           |    NO    |      10       | the time of second for make a file and flush to snowflake                            |
| flush_size_bytes     |    NO    | 100*1024*1024 | the size of bytes for make a file and flush to snowflake                             |

The Snowflake Sink tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the
position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm \
  -p 31080:8080 \
  -v ${PWD}:/vanus-connect/config \
  --name sink-snowflake public.ecr.aws/vanus/connector/sink-snowflake
```

### Test

Open a terminal and use the following command to send a CloudEvent to the Sink. The data field must be according to your
database structure.

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
  "id" : "88767821-92c2-477d-9a6f-bfdfbed19c6a",
  "source" : "quickstart",
  "specversion" : "1.0",
  "type" : "quickstart",
  "time" : "2022-07-08T03:17:03.139Z",
  "datacontenttype" : "application/json",
  "data" : {
    "id":18,
    "name":"xdl",
    "description":"Development Manager",
    "date": "2022-07-06"
  }
}'
```

Open the snowflake console and use the following command to make sure Snowflake has the data

```sql
select * from public.vanus_test;
```

### Clean resource

```shell
docker stop sink-snowflake
```

## Run in Kubernetes

```shell
kubectl apply -f sink-snowflake.yaml
```

```yaml
apiVersion: v1
kind: Service
metadata:
  name: sink-snowflake
  namespace: vanus
spec:
  selector:
    app: sink-snowflake
  type: ClusterIP
  ports:
    - port: 8080
      name: sink-snowflake
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: sink-snowflake
  namespace: vanus
data:
  config.yml: |-
    port: 8080
    snowflake:
      host: "myaccount.ap-northeast-1.aws.snowflakecomputing.com"
      username: "vanus_user"
      password: "snowflake"
      role: "ACCOUNTADMIN"
      warehouse: "xxxxxx"
      database: "VANUS_DB"
      schema: "public"
      table: "vanus_test"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sink-snowflake
  namespace: vanus
  labels:
    app: sink-snowflake
spec:
  selector:
    matchLabels:
      app: sink-snowflake
  replicas: 1
  template:
    metadata:
      labels:
        app: sink-snowflake
    spec:
      containers:
        - name: sink-snowflake
          image: public.ecr.aws/vanus/connector/sink-snowflake
          imagePullPolicy: Always
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "128Mi"
              cpu: "100m"
          ports:
            - name: http
              containerPort: 8080
          volumeMounts:
            - name: config
              mountPath: /vanus-connect/config
      volumes:
        - name: config
          configMap:
            name: sink-snowflake
```

## Integrate with Vanus

This section shows how a sink connector can receive CloudEvents from a
running [Vanus cluster](https://github.com/linkall-labs/vanus).

1. Run the sink-snowflake.yaml

```shell
kubectl apply -f sink-snowflake.yaml
```

2. Create an eventbus

```shell
vsctl eventbus create --name quick-start
```

3. Create a subscription (the sink should be specified as the sink service address or the host name with its port)

```shell
vsctl subscription create \
  --name quick-start \
  --eventbus quick-start \
  --sink 'http://sink-snowflake:8080'
```

[vc]: https://www.vanus.dev/introduction/concepts#vanus-connect
[snowflake]: https://www.snowflake.com
[account]: https://docs.snowflake.com/en/user-guide/admin-account-identifier
[role]: https://docs.snowflake.com/en/user-guide/security-access-control-overview#roles
[warehouse]: https://docs.snowflake.com/en/user-guide/warehouses-overview#overview-of-warehouses
[database]: https://docs.snowflake.com/en/sql-reference/ddl-database
[loadfile]: https://docs.snowflake.com/en/user-guide/data-load-local-file-system
[jdbc-parameter]: https://docs.snowflake.com/en/user-guide/jdbc-parameters
