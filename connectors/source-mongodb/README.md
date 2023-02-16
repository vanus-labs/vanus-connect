---
title: MongoDB
---

# MongoDB Source

## Introduction

The MongoDB Source is a [Vanus Connector][vc] which aims to
capturing mongodb [ChangeEvent](https://www.mongodb.com/docs/manual/reference/change-events/) use [Debezium][debezium]
and convert to a CloudEvent.

## Quickstart

This section shows how MongoDB Source convert MongoDB data to a CloudEvent.

### Prerequisites

- Have a container runtime (i.e., docker).
- Have a MongoDB server.

### create config file

change configurations to yours.

```shell
cat << EOF > config.yml
target: http://localhost:31081

name: "quick-start"
hosts: "127.0.0.1:27017",
credential:
  username: "vanus"
  password: "abc123"
  auth_source: "admin"
database_include: [ "test" ]
collection_include: [ "test.demo" ]
store:
  type: "FILE"
  pathname: "/tmp/vanus-connect/source-mongodb/offset.data",

EOF
```

| Name                    | Required  |   Default   | Description                                                                                                                                                                                                                                                                                                                                                         |
|:------------------------|:---------:|:-----------:|:--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| target                  |    YES    |     ""      | the target URL which MongoDB Source will send CloudEvents to                                                                                                                                                                                                                                                                                                        |
| store.type              |    NO     |   MEMORY    | KV store type for metadata, one of MEMORY, FILE                                                                                                                                                                                                                                                                                                                     |
| store.pathname          |    NO     |      -      | file pathname if `store.type=FILE                                                                                                                                                                                                                                                                                                                                   |
| name                    |    YES    |      -      | Unique name for the connector. Attempting to register again with the same name will fail.                                                                                                                                                                                                                                                                           |
| connection_url          |    NO     |      -      | Specifies a connection string that the connector uses during the initial discovery of a MongoDB replica set. To use this option, you must set the value of mongodb.members.auto.discover to true. Do not set this property and the mongodb.hosts property at the same time.                                                                                         |
| hosts                   |    NO     | empty array | The host addresses to use to connect to the MongoDB replica set                                                                                                                                                                                                                                                                                                     |
| credential.username     |    NO     |      -      | username of mongodb                                                                                                                                                                                                                                                                                                                                                 |
| credential.password     |    NO     |      -      | password of mongodb                                                                                                                                                                                                                                                                                                                                                 |
| credential.auth_source  |    NO     |      -      | authSource of mongodb                                                                                                                                                                                                                                                                                                                                               |
| database_include        |    NO     | empty array | Database names to be monitored; any database name not included in database.include is excluded from monitoring. By default all databases are monitored. Must not be used with database.exclude                                                                                                                                                                      |
| database_exclude        |    NO     | empty array | Database names to be excluded from monitoring; any database name not included in database.exclude is monitored. Must not be used with database.include                                                                                                                                                                                                              |
| collection_include      |    NO     | empty array | Match fully-qualified namespaces for MongoDB collections to be monitored; any collection not included in collection_include is excluded from monitoring. Each identifier is of the form databaseName.collectionName. By default the connector will monitor all collections except those in the local and admin databases. Must not be used with collection_exclude. |
| collection_exclude      |    NO     | empty array | Match fully-qualified namespaces for MongoDB collections to be excluded from monitoring; any collection not included in collection_exclude is monitored. Each identifier is of the form databaseName.collectionName. Must not be used with collection_include                                                                                                       |

The MongoDB Source tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the
position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

it assumes that the mongodb instance doesn't need authentication. For how to use authentication please see
[secret](#secret) section.

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-mongodb public.ecr.aws/vanus/connector/source-mongodb
```

### Test

Open a terminal and use the following command to run a Display sink, which receives and prints CloudEvents.

```shell
docker run -it --rm \
  -p 31081:8080 \
  --name sink-display public.ecr.aws/vanus/connector/sink-display
```

Make sure the `target` value in your config file is `http://localhost:31081` so that the Source can send CloudEvents to
our Display Sink.

### capture a insert event

if you insert a new document to your mongodb instance

```shell
db.mongo_source.insert({"test":"demo"})
```

and you will receive an event like this:

```json
{
  "specversion": "1.0",
  "id": "4cbc7a65-5338-41aa-8f16-8fd164146975",
  "source": "/debezium/mongodb/test",
  "type": "io.debezium.mongodb.datachangeevent",
  "datacontenttype": "application/json",
  "time": "2022-12-21T21:23:40Z",
  "data": {
    "after": {
      "_id": "63a3795c8835b568e786e26a",
      "test": "demo"
    }
  },
  "iodebeziumversion": "2.0.1.Final",
  "xvanuslogoffset": "AAAAAAAAAHk=",
  "iodebeziumord": "1",
  "iodebeziumdb": "test",
  "iodebeziumrs": "replicaset-01",
  "iodebeziumname": "test",
  "iodebeziumcollection": "mongo_source",
  "xvanusstime": "2022-12-21T21:23:41.109Z",
  "iodebeziumsnapshot": "false",
  "iodebeziumtsms": "1671657820000",
  "xvanuseventbus": "test",
  "iodebeziumconnector": "mongodb",
  "iodebeziumop": "c"
}
```

### clean resource

```shell
docker stop source-mongodb sink-display
```

## Run in Kubernetes

```shell
kubectl apply -f sink-mongodb.yaml
```

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: source-mongodb
  namespace: vanus
data:
  config.yml: |-
    target: "http://localhost:8080",
    store:
      type: "FILE"
      pathname: "/tmp/vanus-connect/source-mongodb/offset.data",
    name: "test",
    hosts: "127.0.0.1:27017",
    credential:
      username: "vanus"
      password: "abc123"
      auth_source: "admin"
    database_include: ["test"]
    collection_include: ["test.demo"]

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: source-mongodb
  namespace: vanus
  labels:
    app: source-mongodb
spec:
  selector:
    matchLabels:
      app: source-mongodb
  replicas: 1
  template:
    metadata:
      labels:
        app: source-mongodb
    spec:
      containers:
        - name: source-mongodb
          image: public.ecr.aws/vanus/connector/source-mongodb:latest
          imagePullPolicy: Always
          ports:
            - name: http
              containerPort: 8080
          volumeMounts:
            - name: config
              mountPath: /vanus-connect/config
      volumes:
        - name: config
          configMap:
            name: source-mongodb
```

## Integrate with Vanus

This section shows how a sink connector can receive CloudEvents from a
running [Vanus cluster](https://github.com/linkall-labs/vanus).

1. Run the sink-mongodb.yaml

```shell
kubectl apply -f sink-mongodb.yaml
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
  --sink 'http://sink-mongodb:8080'
```

[vc]: https://www.vanus.ai/introduction/concepts#vanus-connect
[debezium]: https://debezium.io/documentation/reference/2.1/connectors/mongodb.html