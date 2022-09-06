# MongoDB Source Connector

## Introduction

This connector capturing mongodb [ChangeEvent](https://www.mongodb.com/docs/manual/reference/change-events/)

## Quickstart

### create config file

```shell
cat << EOF > config.yml
# change this hosts to your mongodb's address
{
  "v_target": "http://localhost:8080",
  "v_store_file": "/vance/tmp/offset.data",
  "db_hosts":[
    "127.0.0.1:27017"
  ],
  "port": 8080
}
EOF
```

For full configuration, you can see [config](#config) section.

### run mongodb-source

it assumes that the mongodb instance doesn't need authentication. For how to use authentication please see
[secret](#secret) section.

```shell
docker run -d \
  -p 8080:8080 \
  -v ${PWD}:/vance/config \
  -v /tmp:/vance/tmp \
  --name mongodb-source \
  --rm public.ecr.aws/vanus/connector/mongodb-source:dev
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
    "id":"630e32fa020bac5f5f1dcb74",
    "source":"mongodb.replicaset-01.test.mongo_source",
    "type":"test.mongo_source",
    "datacontenttype":"application/json",
    "time":"2022-08-30T15:55:38Z",
    "data":{
        "metadata":{
            "id":"630e32fa020bac5f5f1dcb74",
            "recognized":true,
            "extension":{
                "ord":1,
                "rs":"replicaset-01",
                "collection":"mongo_source",
                "version":"1.9.4.Final",
                "connector":"mongodb",
                "name":"replica-set01",
                "ts_ms":1661874938000,
                "snapshot":"false",
                "db":"test"
            }
        },
        "op":"INSERT",
        "insert":{
            "document":{
                "_id":{
                    "$oid":"630e32fa020bac5f5f1dcb74"
                },
                "test":"demo"
            }
        }
    },
    "vancemongodbversion":"1.9.4.Final",
    "vancemongodboperation":"INSERT",
    "vancemongodbrecognized":true,
    "vancemongodbsnapshot":"false",
    "vancemongodbname":"replica-set01",
    "vancemongodbord":""
}
```

please see [Event Structure](#Event Structure) to understanding it.

### clean resource

```shell
docker stop mongodb-source
```

## Configuration

### config

```json
{
  "v_target": "http://localhost:8080",
  "v_store_file": "/tmp/vance/source-mongodb/offset.data",
  "name": "test",
  "db_name": "replica-set01",
  "db_hosts": "127.0.0.1:27017",
  "database": {
    "include": [],
    "exclude": []
  },
  "collection": {
    "include": [],
    "exclude": []
  }
}
```

| Name               | Required |   Default   | Description                                                                                                                                                                                                                                                                                                                                                         |
|:-------------------|:--------:|:-----------:|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| v_target           | **YES**  |      -      | Target URL will send CloudEvents to                                                                                                                                                                                                                                                                                                                                 |
| v_store_file       | **YES**  |      -      | KV store file name                                                                                                                                                                                                                                                                                                                                                  |
| name               | **YES**  |      -      | Unique name for the connector. Attempting to register again with the same name will fail.                                                                                                                                                                                                                                                                           |
| db_hosts           | **YES**  |      -      | The host addresses to use to connect to the MongoDB replica set                                                                                                                                                                                                                                                                                                     |
| db_name            | **YES**  |      -      | A unique name that identifies the connector and/or MongoDB replica set or sharded cluster that this connector monitors.                                                                                                                                                                                                                                             |
| database.include   | optional | empty array | Database names to be monitored; any database name not included in database.include is excluded from monitoring. By default all databases are monitored. Must not be used with database.exclude                                                                                                                                                                      |
| database.exclude   | optional | empty array | Database names to be excluded from monitoring; any database name not included in database.exclude is monitored. Must not be used with database.include                                                                                                                                                                                                              |
| collection.include | optional | empty array | Match fully-qualified namespaces for MongoDB collections to be monitored; any collection not included in collection.include is excluded from monitoring. Each identifier is of the form databaseName.collectionName. By default the connector will monitor all collections except those in the local and admin databases. Must not be used with collection.exclude. |
| collection.exclude | optional | empty array | Match fully-qualified namespaces for MongoDB collections to be excluded from monitoring; any collection not included in collection.exclude is monitored. Each identifier is of the form databaseName.collectionName. Must not be used with collection.include                                                                                                       |

Note: the `name` property can't be modified once it has been started.

Excepting `v_target` and `v_store_file`, all item are listed here mapping
to [debezium-mongo](https://debezium.io/documentation/reference/stable/connectors/mongodb.html#mongodb-example-configuration), the properties are not listed are not supported now.

### secret

| Name       | Required | Default | Description                      |
|:-----------|:--------:|:-------:|----------------------------------|
| username   | **YES**  |    -    | the username to connect mongodb  |
| password   | **YES**  |    -    | the password to connect mongodb  |
| authSource |    NO    |  admin  | the authSource to authentication |

The `user` and `password` are required only when MongoDB is configured to use authentication. This `authSource` required
only when MongoDB is configured to use authentication with another authentication database than admin.

- example: create a `secert.json` that its content like follow, and mount it to container inside.

```json
{
  "username": "test",
  "password": "123456",
  "authSource": "admin"
}
```

and

```shell
docker run -d \
  -p 8080:8080 \
  -v ${PWD}:/vance/config \
  --env CONNECTOR_SECRET_ENABLE=true 
  --name mongodb-source \
  --rm public.ecr.aws/vanus/connector/mongodb-source:dev
```

## Deploy

### using k8s(recommended)

```shell
kubectl apply -f https://raw.githubusercontent.com/linkall-labs/vance/main/connectors/database/mongodb-source/mongodb-source.yml
```

### using vance Operator

Coming soon, it depends on Vance Operator, the experience of it will be like follow:

```shell
kubectl apply -f https://raw.githubusercontent.com/linkall-labs/vance/main/connectors/database/mongodb-source/crd.yml
```

or

```shell
vsctl connectors create mongodb --source --config /xxx/config.josn --secret /xxx/secret.json
```

## Event Structure

The output events' structure is a [CloudEvent](https://github.com/cloudevents/spec) format, and each field are explained
follows.

the original `ChangeEvent` can be found in [official document](https://www.mongodb.com/docs/manual/reference/change-events/)

| Field                  | Required | Description                                                                                                                                 |
|------------------------|:--------:|---------------------------------------------------------------------------------------------------------------------------------------------|
| id                     | **YES**  | the bson`_id` will be set as the id                                                                                                         |
| source                 | **YES**  | `mongodb.{relicaset_name}.{db_name}.{collection_name}`                                                                                      |
| type                   | **YES**  | `{db_name}.{collection_name}`                                                                                                               |
| time                   | **YES**  | the time of this event generated with RFC3339 encoding                                                                                      |
| data                   | **YES**  | the body of`ChangeEvent`, it's defined as `Event` in [mongodb.proto](../../proto/database/mongodb.proto)                                    |
| data.metadata          | **YES**  | the metadata of this event, it's defined as`Metadata` in [base.proto](../../proto/base/base.proto) , in the most cases users can be ignored |
| data.op                | **YES**  | the event operation of this event, it's defined as`Operation` in  [database.proto](../../proto/database/database.proto)                     |
| data.raw               |    NO    | the raw data of this event, it's defined as "Raw" in[database.proto](../../proto/database/database.proto)                                   |
| data.insert            |    NO    | it's defined as`InsertEvent` in [mongodb.proto](../../proto/database/mongodb.proto)                                                         |
| data.update            |    NO    | it's defined as`UpdateEvent` in [mongodb.proto](../../proto/database/mongodb.proto)                                                         |
| vancemongodbrecognized | **YES**  | if this event was recognized with well-Event Structure, the further explanation in [Unrecognized Event](#unrecognized-event)                |
| vancemongodboperation  | **YES**  | the operation type of this event, it's enum in`insert, update, delete`                                                                      |
| vancemongodb*          | **YES**  | other metadata from debezium may helpful                                                                                                    |

`Required=YES` means it must appear in event, `NO` means it only appears in some conditional cases.

### Create Event

```json
{
  "specversion":"1.0",
  "id":"630e32fa020bac5f5f1dcb74",
  "source":"mongodb.replicaset-01.test.mongo_source",
  "type":"test.mongo_source",
  "datacontenttype":"application/json",
  "time":"2022-08-30T15:55:38Z",
  "data":{
    "metadata":{
      "id":"630e32fa020bac5f5f1dcb74",
      "recognized":true,
      "extension":{
        "ord":1,
        "rs":"replicaset-01",
        "collection":"mongo_source",
        "version":"1.9.4.Final",
        "connector":"mongodb",
        "name":"replica-set01",
        "ts_ms":1661874938000,
        "snapshot":"false",
        "db":"test"
      }
    },
    "op":"INSERT",
    "insert":{
      "document":{
        "_id":{
          "$oid":"630e32fa020bac5f5f1dcb74"
        },
        "test":"demo"
      }
    }
  },
  "vancemongodbversion":"1.9.4.Final",
  "vancemongodboperation":"INSERT",
  "vancemongodbrecognized":true,
  "vancemongodbsnapshot":"false",
  "vancemongodbname":"replica-set01",
  "vancemongodbord":""
}
```

### Update Event

```json
{
  "specversion":"1.0",
  "id":"630e3293020bac5f5f1dcb73",
  "source":"mongodb.replicaset-01.test.mongo_source",
  "type":"test.mongo_source",
  "datacontenttype":"application/json",
  "time":"2022-08-30T16:03:30Z",
  "data":{
    "metadata":{
      "id":"630e3293020bac5f5f1dcb73",
      "recognized":true,
      "extension":{
        "ord":1,
        "rs":"replicaset-01",
        "collection":"mongo_source",
        "version":"1.9.4.Final",
        "connector":"mongodb",
        "name":"replica-set01",
        "ts_ms":1661875410000,
        "snapshot":"false",
        "db":"test"
      }
    },
    "op":"UPDATE",
    "insert":{
      "document":{
        "_id":{
          "$oid":"630e3293020bac5f5f1dcb73"
        },
        "test":"update"
      }
    },
    "update":{
      "updateDescription":{
        "removedFields":[],
        "truncatedArrays":[],
        "updatedFields":{
          "test":"update"
        }
      }
    }
  },
  "vancemongodbrecognized":true,
  "vancemongodbsnapshot":"false",
  "vancemongodbname":"replica-set01",
  "vancemongodbord":"",
  "vancemongodbversion":"1.9.4.Final",
  "vancemongodboperation":"UPDATE"
}
```

### Delete Event

```json
{
  "specversion":"1.0",
  "id":"630e3293020bac5f5f1dcb73",
  "source":"mongodb.replicaset-01.test.mongo_source",
  "type":"test.mongo_source",
  "datacontenttype":"application/json",
  "time":"2022-08-30T16:04:35Z",
  "data":{
    "metadata":{
      "id":"630e3293020bac5f5f1dcb73",
      "recognized":true,
      "extension":{
        "ord":1,
        "rs":"replicaset-01",
        "collection":"mongo_source",
        "version":"1.9.4.Final",
        "connector":"mongodb",
        "name":"replica-set01",
        "ts_ms":1661875475000,
        "snapshot":"false",
        "db":"test"
      }
    },
    "op":"DELETE"
  },
  "vancemongodbord":"",
  "vancemongodbversion":"1.9.4.Final",
  "vancemongodboperation":"DELETE",
  "vancemongodbrecognized":true,
  "vancemongodbsnapshot":"false",
  "vancemongodbname":"replica-set01"
}
```

### Unrecognized Event

Although we do our best to deal with different events, but it's not easy to make sure that all raw data are parsed to a
structured format because we just only can deal the raw we knew. So, if there is error happened when we try to parse raw
data, we will see the event is an unrecognized event and set `vancemongodbrecognized` to false instead of discard in
order to guaranty no data loss.

User should put event with `vancemongodbrecognized=false` to the `deadLetter` in the further step. we're appreciated
that you can create an issue to feedback us about the unrecognized event, we will fix it as soon as possible.

```json
{
  "specversion": "1.0",
  "id": "unknown",
  "source": "unknown",
  "type": "unknown",
  "datacontenttype": "application/json",
  "time": "unknown",
  "data": {
    "raw": {
      "key":"xxx",
      "value": "xxxx"
    }
  },
  "vancemongodbrecognized": false,
  "vancemongodboperation": "unknown"
}
```

## Example
(TODO here can refer another page)
Use `mongo-source` and `mysql-sink` to build a data pipeline in minutes.
