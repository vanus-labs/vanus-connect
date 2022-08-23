# MongoDB Connector

## How to use

### quickstart

```bash
docker run -it --rm public.ecr.aws/vanus/connector/mongodb:latest /run/start.sh \
  --volume /xxx/secret.json /var/mongodb/secret.json \
  --env MONGODB_HOSTS=xxx \
  --env MONGODB_NAME=xxx \
  --env MONGODB_AUTHSOURCE=xxx
```

### vance

```bash
kubectl apply -f https://raw.githubusercontent.com/linkall-labs/vance/mongo-connector/connectors/mongodb/mongodb.yml
```

### k8s

```bash
kubectl apply -f https://raw.githubusercontent.com/linkall-labs/vance/mongo-connector/connectors/mongodb/mongodb-bare.yml
```

## Schema

The output events' schema is a [CloudEvent](https://github.com/cloudevents/spec) format, and each field are explained follows.

the original `ChangeEvent` can be found in [official document](https://www.mongodb.com/docs/manual/reference/change-events/)


| field                  | description                                                                                                       |
| ------------------------ | ------------------------------------------------------------------------------------------------------------------- |
| id                     | the bson`_id` will be set as the id                                                                               |
| source                 | mongodb.${relicaset_name}.${db_name}.${collection_name}                                                           |
| type                   | ${db_name}.${collection_name}                                                                                     |
| time                   | the time of this event generated with RFC3339 encoding                                                            |
| data                   | the body of `ChangeEvent`                                                                                           |
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
  "specversion":"1.0",
  "id": "unknown",
  "source": "unknown",
  "type": "unknown",
  "datacontenttype": "application/json",
  "time":"unknown",
  "data": {
    "rawKey": "xxxxx",
    "rawValue": "xxxx"
  },
  "vancemongodbrecognized": false,
  "vancemongodboperation": "unknown"
}
```

## Acknowledgement

The MongoDB Connector built on [debezium](https://github.com/debezium/debezium)
