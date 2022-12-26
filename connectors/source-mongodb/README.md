---
title: MongoDB
---

# MongoDB Source

## Introduction

This connector capturing mongodb [ChangeEvent](https://www.mongodb.com/docs/manual/reference/change-events/)

## Quickstart

### create config file

change configurations to yours.

```shell
cat << EOF > config.yml
target: "<url of cloud event receiver>"
name: "quick-start"
hosts: [<your_mongodb_hosts >]
EOF
```

For full configuration, you can see [config](#config) section.

### Start Using Docker

it assumes that the mongodb instance doesn't need authentication. For how to use authentication please see
[secret](#secret) section.

```shell
docker run -d --rm \
  --network host \
  -v ${PWD}:/vance/config \
  -e CONNECTOR_CONFIG=/vance/config/config.yml \
  --name source-mongodb public.ecr.aws/vanus/connector/source-mongodb:latest
```

### capture a insert event

if you insert a new document to your mongodb instance

```shell
db.mongo_source.insert({"test":"demo"})
```

and you will receive an event like this:

```json
{
  "specversion":"1.0",
  "id":"4cbc7a65-5338-41aa-8f16-8fd164146975",
  "source":"/debezium/mongodb/test",
  "type":"io.debezium.mongodb.datachangeevent",
  "datacontenttype":"application/json",
  "time":"2022-12-21T21:23:40Z",
  "data":{
    "after":{
      "_id":"63a3795c8835b568e786e26a",
      "test":"demo"
    }
  },
  "iodebeziumversion":"2.0.1.Final",
  "xvanuslogoffset":"AAAAAAAAAHk=",
  "iodebeziumord":"1",
  "iodebeziumdb":"test",
  "iodebeziumrs":"replicaset-01",
  "iodebeziumname":"test",
  "iodebeziumcollection":"mongo_source",
  "xvanusstime":"2022-12-21T21:23:41.109Z",
  "iodebeziumsnapshot":"false",
  "iodebeziumtsms":"1671657820000",
  "xvanuseventbus":"test",
  "iodebeziumconnector":"mongodb",
  "iodebeziumop":"c"
}
```

please see [Event Structure](#Event Structure) to understanding it.

### clean resource

```shell
docker stop source-mongodb
```

## How to use

### Configuration

```yaml
target: "http://localhost:8080",
store:
  type: "FILE"
  pathname: "/tmp/vance/source-mongodb/offset.data",
name: "test",
hosts: "127.0.0.1:27017",
credential:
  username: "vanus"
  password: "abc123"
  auth_source: "admin"
database_include: ["test"]
collection_include: ["test.demo"]
```

| Name                   | Required |   Default   | Description                                                                                                                                                                                                                                                                                                                                                         |
|:-----------------------|:--------:|:-----------:|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| target                 | **YES**  |      -      | Target URL will send CloudEvents to                                                                                                                                                                                                                                                                                                                                 |
| store.type             |    NO    |   MEMORY    | KV store type for metadata, one of MEMORY, FILE                                                                                                                                                                                                                                                                                                                     |
| store.pathname         |    NO    |      -      | file pathname if `store.type=FILE                                                                                                                                                                                                                                                                                                                                   |
| name                   | **YES**  |      -      | Unique name for the connector. Attempting to register again with the same name will fail.                                                                                                                                                                                                                                                                           |
| connection_url         |    NO    |      -      | Specifies a connection string that the connector uses during the initial discovery of a MongoDB replica set. To use this option, you must set the value of mongodb.members.auto.discover to true. Do not set this property and the mongodb.hosts property at the same time.                                                                                         |
| hosts                  |    NO    | empty array | The host addresses to use to connect to the MongoDB replica set                                                                                                                                                                                                                                                                                                     |
| credential.username    |    NO    |      -      | username of mongodb                                                                                                                                                                                                                                                                                                                                                 |
| credential.password    |    NO    |      -      | password of mongodb                                                                                                                                                                                                                                                                                                                                                 |
| credential.auth_source |    NO    |      -      | authSource of mongodb                                                                                                                                                                                                                                                                                                                                               |
| database_include       |    NO    | empty array | Database names to be monitored; any database name not included in database.include is excluded from monitoring. By default all databases are monitored. Must not be used with database.exclude                                                                                                                                                                      |
| database_exclude       |    NO    | empty array | Database names to be excluded from monitoring; any database name not included in database.exclude is monitored. Must not be used with database.include                                                                                                                                                                                                              |
| collection_include     |    NO    | empty array | Match fully-qualified namespaces for MongoDB collections to be monitored; any collection not included in collection_include is excluded from monitoring. Each identifier is of the form databaseName.collectionName. By default the connector will monitor all collections except those in the local and admin databases. Must not be used with collection_exclude. |
| collection_exclude     |    NO    | empty array | Match fully-qualified namespaces for MongoDB collections to be excluded from monitoring; any collection not included in collection_exclude is monitored. Each identifier is of the form databaseName.collectionName. Must not be used with collection_include                                                                                                       | 

Note: the `name` can't be modified once it has been started.

For more explanation, you can view https://debezium.io/documentation/reference/stable/connectors/mongodb.html#mongodb-example-configuration
## Example
(TODO here can refer another page)
Use `mongo-source` and `mysql-sink` to build a data pipeline in minutes.

## Run in Kubernetes
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
      pathname: "/tmp/vance/source-mongodb/offset.data",
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
              mountPath: /vance/config
      volumes:
        - name: config
          configMap:
            name: source-mongodb
```