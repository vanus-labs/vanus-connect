---
title: MongoDB
---

#  MongoDB Sink

## Introduction

The Sink MongoDB is a [Vance Connector][vc] which aims to handle incoming CloudEvents in a way that extracts the `data` part of the
original event and insert/update/delete this data to mongodb.

For examples, If incoming event looks like:

```json
{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "quick-start",
    "specversion": "1.0",
    "type": "sink-mongodb",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "xvdatabasedb": "test",
    "xvdatabasecoll": "demo",
    "data": {
        "inserts": [
            {
                "scenario":"quick-start"
            }
        ]
    }
}
```

which equals to

```shell
use test;
db.demo.insertMany([{"scenario":"quick-start"}])
```

## Quickstart

### create config file

use your mongodb's hosts, username and password.

```shell
cat << EOF > config.yml
connection_uri: "mongodb+srv://<hosts>/?retryWrites=true&w=majority"
credential:
  username: "<username>"
  password: "<password>"
EOF
```

### start with docker

```shell
docker run -d --rm \
  -p 31080:8080 \
  -v ${PWD}:/vance/config \
  --name sink-mongodb public.ecr.aws/vanus/connector/sink-mongodb:latest
```

### insert document to mongodb

For more details on how to understand, please see [Examples](#examples) section.

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "quick-start",
    "specversion": "1.0",
    "type": "sink-mongodb",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "xvdatabasedb": "test",
    "xvdatabasecoll": "demo",
    "data": {
        "inserts": [
            {
                "scenario": "quick-start"
            }
        ]
    }
}'
```

find in mongodb

```shell
shard-0 [primary] test> db.demo.find()
[
  {
    _id: ObjectId("63a56b176dcdb253ae4924f0"),
    scenario: 'quick-start'
  }
]
shard-0 [primary] test>
```

### clean resource

```shell
docker stop sink-mongodb  
```

## How to use

### Configuration

The default path is `/vance/config/config.yml`. if you want to change the default path, you can set env `CONNECTOR_CONFIG` to
tell Sink MongoDB.


| Name                                 | Required | Default | Description                                                                                                                                       |
|:-------------------------------------|:--------:|:-------:|---------------------------------------------------------------------------------------------------------------------------------------------------|
| port                                 |    NO    |  8080   | the pot Sink MongoDB receives incoming events                                                                                                     |
| connection_uri                       | **YES**  |    -    | the URI to connect MongoDB, view[Connection String URI Format](https://www.mongodb.com/docs/manual/reference/connection-string/) for more details |
| credential.username                  |    NO    |    -    | https://www.mongodb.com/docs/drivers/go/current/fundamentals/auth/                                                                                |
| credential.password                  |    NO    |    -    | https://www.mongodb.com/docs/drivers/go/current/fundamentals/auth/                                                                                |
| credential.auth_source               |    NO    |    -    | https://www.mongodb.com/docs/drivers/go/current/fundamentals/auth/                                                                                |
| credential.auth_mechanism            |    NO    |    -    | https://www.mongodb.com/docs/drivers/go/current/fundamentals/auth/                                                                                |
| credential.auth_mechanism_properties |    NO    |    -    | https://www.mongodb.com/docs/drivers/go/current/fundamentals/auth/                                                                                |

```yaml
connection_uri: "mongodb+srv://<host1>,<host2>/?retryWrites=true&w=majority"
credential:
  username: "vanus"
  password: "demo"
  auth_source: "admin"
```

### Extension Attributes

Sink Source has defined a few [CloudEvents Extension Attribute](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#extension-context-attributes)
to determine how to process event.


| Attribute      | Required | Examples | Description                          |
|:---------------|:--------:|----------|--------------------------------------|
| xvdatabasedb   | **YES**  | test     | which database this event write to   |
| xvdatabasecoll | **YES**  | demo     | which collection this event write to |

### Data


| Item                  | Required |   Type   | Default | Description                                                   |
|:----------------------|:--------:|:--------:|:-------:|---------------------------------------------------------------|
| inserts               |    NO    | []Object |  null   | insert data                                                   |
| updates               |    NO    | []Object |  null   | https://www.mongodb.com/docs/manual/tutorial/update-documents |
| updates[].filter      |    NO    |  Object  |  null   |                                                               |
| updates[].update      |    NO    |  Object  |  null   |                                                               |
| updates[].update_many |    NO    | boolean  |  false  | update many records when filter matches more than one         |
| deletes               |    NO    | []Object |  null   | delete data                                                   |
| deletes[].filter      |    NO    |  Object  |  null   | delete data                                                   |
| deletes[].delete_many |    NO    |  Object  |  false  | delete many records when filter matches more than one         |

```json
{
  "inserts":[
    {
      "_id": "63a56aed6dcdb253ae4924ee",
      "key1": "value1"
    },
    {
      "key2": "value2"
    }
  ],
  "updates":[
    {
      "filter":{
        "_id": "63a56aed6dcdb253ae4924ee"
      },
      "update": {
        "$set": {
          "key1": "value2_updated"
        }
      },
      "update_many": true
    }
  ],
  "deletes":[
    {
      "filter": {
        "key2": "value2"
      },
      "delete_many": true
    }
  ]
}
```

### Examples

#### insert multiple documents to mongodb

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "quick-start",
    "specversion": "1.0",
    "type": "sink-mongodb",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "xvdatabasedb": "test",
    "xvdatabasecoll": "demo",
    "data": {
        "inserts": [
            {
                "scenario": "quick-start-1"
            },
            {
                "scenario": "quick-start-2"
            }
        ]
    }
}'
```

#### update multiple documents in mongodb

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "quick-start",
    "specversion": "1.0",
    "type": "sink-mongodb",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "xvdatabasedb": "test",
    "xvdatabasecoll": "demo",
    "data": {
        "updates": [
            {
                "filter":{
                  "scenario": "quick-start-1"
                },
                "update": {
                    "$set": {
                      "scenario": "quick-start-1-updated"
                    }
                },
                "update_many": false
            }
        ]
    }
}'
```

#### delete document

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "quick-start",
    "specversion": "1.0",
    "type": "sink-mongodb",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "xvdatabasedb": "test",
    "xvdatabasecoll": "demo",
    "data": {
        "deletes": [
            {
                "filter":{
                  "scenario": "quick-start-1-updated"
                },                
                "delete_many": false
            }
        ]
    }
}'
```

### Run in kubernetes
```yaml
apiVersion: v1
kind: Service
metadata:
  name: sink-mongodb
  namespace: vanus
spec:
  selector:
    app: sink-mongodb
  type: ClusterIP
  ports:
    - port: 8080
      name: sink-mongodb
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: sink-mongodb
  namespace: vanus
data:
  config.yml: |-
    connection_uri: "mongodb+srv://<hosts>/?retryWrites=true&w=majority"
    credential:
      username: "<username>"
      password: "<password>"

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sink-mongodb
  namespace: vanus
  labels:
    app: sink-mongodb
spec:
  selector:
    matchLabels:
      app: sink-mongodb
  replicas: 1
  template:
    metadata:
      labels:
        app: sink-mongodb
    spec:
      containers:
        - name: sink-mongodb
          image: public.ecr.aws/vanus/connector/sink-mongodb:latest
          imagePullPolicy: Always
          ports:
            - name: http
              containerPort: 8080
          volumeMounts:
            - name: config
              mountPath: /vance/config
          # env: see README for more about how to set env
      volumes:
        - name: config
          configMap:
            name: sink-mongodb
```

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md