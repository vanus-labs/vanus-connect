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
```

```

The event schema that the mongodb source output looks like follows.

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
  }
}
```

## Acknowledgement
The MongoDB Connector built on top of by [debezium](https://github.com/debezium/debezium)