---
title: Doris
---

# Doris Sink

## Introduction

The Doris Sink is a [Vanus Connector][vc] which aims to handle incoming CloudEvents in a way that extracts the `data`
part of the original event and deliver these extracted `data` to [Doris][doris]. The Doris Sink use [Stream Load][stream load]
way to import data.

For example, the incoming CloudEvent looks like:

```json
{
  "specversion": "1.0",
  "id": "4395ffa3-f6de-443c-bf0e-bb9798d26a1d",
  "source": "vanus.source.test",
  "type": "vanus.type.test",
  "datacontenttype": "application/json",
  "time": "2022-11-20T07:05:55.777689Z",
  "data": {
    "id": 1,
    "username": "name",
    "birthday": "2022-11-20"
  }
}
```

The Doris Sink will extract `data` field write to [Doris][doris] table like:

```text
+------+----------+------------+
| id   | username | birthday   |
+------+----------+------------+
|    1 | name     | 2022-11-20 |
+------+----------+------------+
```

## Quickstart

### Prerequisites
- Have a container runtime (i.e., docker).
- Have a [Doris cluster](https://doris.apache.org/docs/dev/get-starting/)

### Create the config file

```shell
cat << EOF > config.yml
port: 8080
secret:
  # doris info
  fenodes: "localhost:8030"
  db_name: "vanus_test"
  table_name: "vanus_test"
  username: "vanus_test"
  password: "123456"
EOF
```

| Name            | Required | Default      | Description                                |
|:----------------|:--------:|:-------------|--------------------------------------------|
| port            |    NO    | 8080         | the port which Doris Sink listens on       |
| fenodes         |   YES    |              | doris fenodes, example: "17.0.0.1:8003"    |
| db_name         |   YES    |              | doris database name                        |
| table_name      |   YES    |              | doris table name                           |
| username        |   YES    |              | doris username                             |
| password        |   YES    |              | doris password                             |
| stream_load     |    NO    |              | doris stream load properties, map struct   |
| load_interval   |    NO    | 5            | doris stream load interval, unit second    |
| load_size       |    NO    | 10*1024*1024 | doris stream load max body size            |
| timeout         |    NO    | 30           | doris stream load timeout, unit second     |

The Doris Sink tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.


### Start with Docker

```shell
docker run -it --rm --network=host\
  -v ${PWD}:/vanus-connect/config \
  --name sink-doris public.ecr.aws/vanus/connector/sink-doris
```

### Test

Connect to Doris and use command to create database and table.

```shell
create database vanus_test;
use vanus_test;
CREATE TABLE IF NOT EXISTS vanus_test.vanus_test
(
    `id` LARGEINT NOT NULL COMMENT "id",
    `username` VARCHAR(64) COMMENT "username",
    `birthday` DATE NOT NULL COMMENT "birthday"
)
AGGREGATE KEY(`id`, `username`, `birthday`)
DISTRIBUTED BY HASH(`id`) BUCKETS 10
PROPERTIES (
    "replication_allocation" = "tag.location.default: 1"
);
```

Open a terminal and use following command to send a CloudEvent to the Sink.

```shell
curl --location --request POST 'localhost:8080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
  "specversion": "1.0",
  "id": "4395ffa3-f6de-443c-bf0e-bb9798d26a1d",
  "source": "vanus.source.test",
  "type": "vanus.type.test",
  "datacontenttype": "application/json",
  "time": "2022-11-20T07:05:55.777689Z",
  "data": {
    "id": 1,
    "username": "name",
    "birthday": "2022-11-20"
  }
}'
```

you will see data in doris table vanus_test

```text
+------+----------+------------+
| id   | username | birthday   |
+------+----------+------------+
|    1 | name     | 2022-11-20 |
+------+----------+------------+
```

### Clean resource

```shell
docker stop sink-doris
```

## Run in Kubernetes

```shell
kubectl apply -f sink-doris.yaml
```

```yaml
apiVersion: v1
kind: Service
metadata:
  name: sink-doris
  namespace: vanus
spec:
  selector:
    app: sink-doris
  type: ClusterIP
  ports:
    - port: 8080
      name: sink-doris
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: sink-doris
  namespace: vanus
data:
  config.yml: |-
    port: 8080
    secret:
      # doris info
      fenodes: "localhost:8030"
      db_name: "vanus_test"
      table_name: "user"
      username: "vanus_test"
      password: "123456"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sink-doris
  namespace: vanus
  labels:
    app: sink-doris
spec:
  selector:
    matchLabels:
      app: sink-doris
  replicas: 1
  template:
    metadata:
      labels:
        app: sink-doris
    spec:
      containers:
        - name: sink-doris
          image: public.ecr.aws/vanus/connector/sink-doris
          imagePullPolicy: Always
          volumeMounts:
            - name: config
              mountPath: /vanus-connect/config
      volumes:
        - name: config
          configMap:
            name: sink-doris

```

## Integrate with Vanus

This section shows how a sink connector can receive CloudEvents from a running [Vanus cluster](https://github.com/linkall-labs/vanus).

1. Run the sink-doris.yaml
```shell
kubectl apply -f sink-doris.yaml
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
  --sink 'http://sink-doris:8080'
```

[vc]: https://www.vanus.dev/introduction/concepts#vanus-connect
[doris]: https://doris.apache.org/docs/summary/basic-summary
[stream load]: https://doris.apache.org/docs/dev/data-operate/import/import-way/stream-load-manual/
