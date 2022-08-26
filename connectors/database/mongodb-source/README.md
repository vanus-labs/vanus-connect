# MongoDB Connector

## Introduction

## How to use

### quickstart

```bash
docker run -it --rm public.ecr.aws/vanus/connector/mongodb:latest /etc/vance/mongodb/start.sh \
  --volume /xxx/config.json /etc/vance/mongodb/config.json \
  --volume /xxx/secret.json /etc/vance/mongodb/secret.json \
  --env MONGODB_CONNECTOR_HOME=/etc/vance/mongodb
```

### vance

```bash
kubectl apply -f https://raw.githubusercontent.com/linkall-labs/vance/mongo-connector/connectors/mongodb/mongodb.yml
```

### k8s

Coming soon, it depends on Vance Operator, the experience of it will be like follow:

```bash
kubectl apply -f https://raw.githubusercontent.com/linkall-labs/vance/main/connectors/xxx/xxx.yml
# or
vsctl connectors create mongodb --source --config /xxx/config.josn --secret /xxx/secret.json
```

## Configuration

### config.json

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

| Name               | Required |   default   | description                                                                                                                                                                                                                                                                                                                                                         |
| :------------------- | :-----------: | :-----------: | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| v_target           |  required  |      -      | Target URL will send CloudEvents to                                                                                                                                                                                                                                                                                                                                 |
| v_store_file       |  required  |      -      | KV store file name                                                                                                                                                                                                                                                                                                                                                  |
| name               |  required  |      -      | Unique name for the connector. Attempting to register again with the same name will fail.                                                                                                                                                                                                                                                                           |
| db_hosts           |  required  |      -      | The host addresses to use to connect to the MongoDB replica set                                                                                                                                                                                                                                                                                                     |
| db_name            |  required  |      -      | A unique name that identifies the connector and/or MongoDB replica set or sharded cluster that this connector monitors.                                                                                                                                                                                                                                             |
| database.include   |  optional  | empty array | Database names to be monitored; any database name not included in database.include is excluded from monitoring. By default all databases are monitored. Must not be used with database.exclude                                                                                                                                                                      |
| database.exclude   |  optional  | empty array | Database names to be excluded from monitoring; any database name not included in database.exclude is monitored. Must not be used with database.include                                                                                                                                                                                                              |
| collection.include |  optional  | empty array | Match fully-qualified namespaces for MongoDB collections to be monitored; any collection not included in collection.include is excluded from monitoring. Each identifier is of the form databaseName.collectionName. By default the connector will monitor all collections except those in the local and admin databases. Must not be used with collection.exclude. |
| collection.exclude |  optional  | empty array | Match fully-qualified namespaces for MongoDB collections to be excluded from monitoring; any collection not included in collection.exclude is monitored. Each identifier is of the form databaseName.collectionName. Must not be used with collection.include                                                                                                       |

Note: the `name` property can't be modified once it has been started.

Excepting `v_target` and `v_store_file`, all item are listed are mapping
to [debezium-mongo](https://debezium.io/documentation/reference/stable/connectors/mongodb.html#mongodb-example-configuration)
, the properties are not listed are not supported now.

### secret.json

```json
{
  "user": "admin",
  "password": "admin",
  "authsoure": "admin"
}
```

| name      | requirement | description                                                     |
| ----------- | ------------- | ----------------------------------------------------------------- |
| user      | optional    | Name of the database user to be used when connecting to MongoDB |
| password  | optional    | Password to be used when connecting to MongoDB                  |
| authsoure | optional    | Database (authentication source) containing MongoDB credentials |

The `user` and `password` are required only when MongoDB is configured to use authentication. This `authsoure` required
only when MongoDB is configured to use authentication with another authentication database than admin.

## Schema

The output events' schema is a [CloudEvent](https://github.com/cloudevents/spec) format, and each field are explained
follows.

the original `ChangeEvent` can be found
in [official document](https://www.mongodb.com/docs/manual/reference/change-events/)

| field                  | description                                                                                                       |
| ------------------------ | ------------------------------------------------------------------------------------------------------------------- |
| id                     | the bson`_id` will be set as the id                                                                               |
| source                 | `mongodb.{relicaset_name}.{db_name}.{collection_name}`                                                            |
| type                   | `{db_name}.{collection_name}`                                                                                     |
| time                   | the time of this event generated with RFC3339 encoding                                                            |
| data                   | the body of`ChangeEvent`                                                                                          |
| data.full              | the full document of each bson, not empty when operation is`insert` and `update`, mapping to`insert.fullDocument` |
| data.changed           | the data changed when updating, mapping to`update.updateDescription`                                              |
| data.changed.updated   | mapping to`update.updateDescription.updatedFields`                                                                |
| data.changed.deleted   | mapping to`update.updateDescription.removedFields`                                                                |
| data.changed.truncated | mapping to`update.updateDescription.trucatedArrays`                                                               |
| vancemongodbrecognized | if this event was recognized with well-schema, there is a detail explanation in follow section                    |
| vancemongodboperation  | the operation type of this event, it's enum in`insert, update, delete`                                            |
| vancemongodb*          | other metadata from debezium may helpful                                                                          |

### Create Event

```json
{
  "specversion": "1.0",
  "id": "6304855bccaea8fcf8a159f2",
  "source": "mongodb.replicaset-01.test.source",
  "type": "test.source",
  "datacontenttype": "application/json",
  "time": "2022-08-23T07:44:27Z",
  "data": {
    "full": {
      "download": "1234",
      "connector": "mongodb",
      "_id": "6304855bccaea8fcf8a159f2",
      "version": "v0.3.0"
    }
  },
  "vancemongodbrecognized": true,
  "vancemongodbversion": "1.9.4.Final",
  "vancemongodbsnapshot": "false",
  "vancemongodbname": "test",
  "vancemongodbord": "1",
  "vancemongodboperation": "insert"
}
```

### Update Event

```json
{
  "specversion": "1.0",
  "id": "6304855bccaea8fcf8a159f2",
  "source": "mongodb.replicaset-01.test.source",
  "type": "test.source",
  "datacontenttype": "application/json",
  "time": "2022-08-23T08:08:05Z",
  "data": {
    "full": {
      "download": "1240",
      "connector": "mongodb",
      "_id": "6304855bccaea8fcf8a159f2",
      "version": "v0.3.0"
    },
    "changed": {
      "updated": {
        "download": 1240
      }
    }
  },
  "vancemongodbrecognized": true,
  "vancemongodbversion": "1.9.4.Final",
  "vancemongodbsnapshot": "false",
  "vancemongodbname": "test",
  "vancemongodbord": "1",
  "vancemongodboperation": "update"
}
```

### Delete Event

```json
{
  "specversion": "1.0",
  "id": "6304855bccaea8fcf8a159f2",
  "source": "mongodb.replicaset-01.test.source",
  "type": "test.source",
  "datacontenttype": "application/json",
  "time": "2022-08-23T08:09:24Z",
  "data": {},
  "vancemongodbord": "1",
  "vancemongodbrecognized": true,
  "vancemongodbversion": "1.9.4.Final",
  "vancemongodbsnapshot": "false",
  "vancemongodbname": "test",
  "vancemongodboperation": "delete"
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
    "rawKey": "xxxxx",
    "rawValue": "xxxx"
  },
  "vancemongodbrecognized": false,
  "vancemongodboperation": "unknown"
}
```

## example

Use mongo-source to build a data pipeline to MySQL in minutes.

## Acknowledgement

The MongoDB Connector built on [debezium](https://github.com/debezium/debezium)
