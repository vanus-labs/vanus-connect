# MongoDB Sink Connector

## Introduction

The MongoDB Sink is a [Vance Connector](https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md), 
which now supports insert/update/delete operations.

## Quickstart

### create config file

```shell
cat << EOF > config.yml
# change this hosts to your mongodb's address
db_hosts:
  - 44.242.140.28:27017
port: 8080
EOF
```

### run mongodb-sink

it assumes that the mongodb instance doesn't need authentication. For how to use authentication please see
[secret](#secret) section.

```shell
docker run -d --rm \
  --network host \
  -p 8080:8080 \
  -v ${PWD}:/vance/config \
  --name mongodb-sink public.ecr.aws/vanus/connector/mongodb-sink:dev
```

### insert document to mongodb
For more details on how to understand, please see [Structure](#structure) and [Examples](#examples) section.

```shell
curl --location --request POST 'http://127.0.0.1:8080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "specversion": "1.0",
    "id": "62ff305f779a73966deb3877",
    "source": "mongodb.replicaset-01.test.source",
    "type": "test.source",
    "datacontenttype": "application/json",
    "time": "2022-08-26T18:42:16Z",
    "data": {
        "op": "INSERT",
        "insert": {
            "document": {
                "a": 1234
            }
        }
    },
    "vancemongosinkdatabase":"test",
    "vancemongosinkcollection": "sink",  
}'
```

### clean resource

```shell
docker stop mongodb-sink  
```

## Configuration

the configuration of mongodb-sink based on [Connection String URI Format](https://www.mongodb.com/docs/v6.0/reference/connection-string/)

### config

| Name     | Required | Default | Description                                     |
|:---------|:--------:|:-------:|-------------------------------------------------|
| db_hosts | **YES**  |    -    | the mongodb cluster hosts                       |
| port     | **YES**  |    -    | the port the mongodb-sink for listening request |

- example

create a `config.yml` with its content as below, and mount it to container inside.

```yaml
db_hosts:
  - 127.0.0.1:27017
port: 8080
```

```shell
docker run -d \
  -p 8080:8080 \
  -v ${PWD}:/vance/config \
  --name mongodb-sink \
  --rm public.ecr.aws/vanus/connector/mongodb-sink:v0.2.0-alpha
```

### secret

| Name       | Required | Default | Description                      |
|:-----------|:--------:|:-------:|----------------------------------|
| username   | **YES**  |    -    | the username to connect mongodb  |
| password   | **YES**  |    -    | the password to connect mongodb  |
| authSource |    NO    |  admin  | the authSource to authentication |

The `user` and `password` are required only when MongoDB is configured to use authentication. This `authSource` required
only when MongoDB is configured to use authentication with another authentication database than admin.

- example: create a `secert.yml` that its content like follow, and mount it to container inside.

```yaml
username: "test"
password: "123456"
authSource: "admin"
```

```shell
docker run -d \
  -p 8080:8080 \
  -v ${PWD}:/vance/config \
  --env CONNECTOR_SECRET_ENABLE=true 
  --name mongodb-sink \
  --rm public.ecr.aws/vanus/connector/mongodb-sink:v0.2.0-alpha
```

## Deploy

### using k8s(recommended)

```shell
kubectl apply -f https://raw.githubusercontent.com/linkall-labs/vance/main/connectors/database/mongodb-sink/mongodb-sink.yml
```

### using vance Operator

Coming soon, it depends on Vance Operator, the experience of it will be like follow:

```shell
kubectl apply -f https://raw.githubusercontent.com/linkall-labs/vance/main/connectors/database/mongodb-sink/crd.yml
```

or

```shell
vsctl connectors create mongodb --source --config /xxx/config.josn --secret /xxx/secret.json
```

## Structure

The input events' structure is a [CloudEvent](https://github.com/cloudevents/spec) format, and each field are explained
follows.

the original `ChangeEvent` can be found in [official document](https://www.mongodb.com/docs/manual/reference/change-events/)

| Field                    | Required | Description                                                                                                                                 |
|--------------------------|:--------:|---------------------------------------------------------------------------------------------------------------------------------------------|
| id                       | **YES**  | the bson`_id` will be set as the id                                                                                                         |
| source                   | **YES**  | where the event come from                                                                                                                   |
| type                     | **YES**  | what's the event's type                                                                                                                     |
| time                     |    NO    | the time of this event generated with RFC3339 encoding                                                                                      |
| data                     | **YES**  | the body of`ChangeEvent`, it's defined as `Event` in [mongodb.proto](../../proto/database/mongodb.proto)                                    |
| data.metadata            |    NO    | the metadata of this event, it's defined as`Metadata` in [base.proto](../../proto/base/base.proto) , in the most cases users can be ignored |
| data.op                  | **YES**  | the event operation of this event, it's defined as`Operation` in [database.proto](../../proto/database/database.proto)                      |
| data.raw                 |    NO    | the raw data of this event, it's defined as "Raw" in [database.proto](../../proto/database/database.proto)                                  |
| data.insert              |    NO    | it's defined as`InsertEvent` in [mongodb.proto](../../proto/database/mongodb.proto)                                                         |
| data.update              |    NO    | it's defined as`UpdateEvent` in [mongodb.proto](../../proto/database/mongodb.proto)                                                         |
| vancemongosinkdatabase   | **YES**  | which `database` the event into                                                                                                             |
| vancemongosinkcollection | **YES**  | which `collection` the event into                                                                                                           |

## Examples

### insert document

```shell
curl --location --request POST 'http://127.0.0.1:8080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "specversion": "1.0",
    "id": "62ff305f779a73966deb3877",
    "source": "mongodb.replicaset-01.test.source",
    "type": "test.source",
    "datacontenttype": "application/json",
    "time": "2022-08-26T18:42:16Z",
    "data": {
        "op": "INSERT",
        "insert": {
            "document": {
                "a": 1234
            }
        }
    },
    "vancemongosinkdatabase":"test",
    "vancemongosinkcollection": "sink",  
}'
```

### update document

```shell
curl --location --request POST 'http://127.0.0.1:8080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "specversion": "1.0",
    "id": "62ff305f779a73966deb3877",
    "source": "mongodb.replicaset-01.test.source",
    "type": "test.source",
    "datacontenttype": "application/json",
    "time": "2022-08-26T18:42:16Z",
    "data": {
        "op": "UPDATE",
        "update": {
            "updateDescription": {
                "removedFields": [],
                "truncatedArrays": [],
                "updatedFields": {
                    "a": 12314
                }
            }
        }
    },
    "vancemongosinkdatabase":"test",
    "vancemongosinkcollection": "sink",  
}'
```

### delete document

```shell
curl --location --request POST 'http://127.0.0.1:8080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "specversion": "1.0",
    "id": "62ff305f779a73966deb3877",
    "source": "mongodb.replicaset-01.test.source",
    "type": "test.source",
    "datacontenttype": "application/json",
    "time": "2022-08-26T18:42:16Z",
    "data": {
        "op": "DELETE"    
    },
    "vancemongosinkdatabase":"test",
    "vancemongosinkcollection": "sink",  
}'
```
