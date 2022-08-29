# MongoDB Connector

## Introduction

This is the template README.md of the connector, a real example can be found in [mongo](../mongodb/README.md)

## How to use

Copy this template to your connector directory, and finish all section.

### quickstart

```bash
cat << EOF > config.yml
# change this hosts to your mongodb's address
db_hosts:
  - xxx.xxx.xxx.xx:27017
port: 8080
EOF

docker run -d \
  -p 8080:8080 \
  -v config.yml:/vance/config/config.yml \
  --rm image.linkall.com/connector/mongo-sink:v0.2.1
```

### vance

Coming soon, it depends on Vance Operator, the experience of it will be like follow:

```bash
kubectl apply -f https://raw.githubusercontent.com/linkall-labs/vance/main/connectors/xxx/xxx.yml
# or
vsctl connectors create mongodb --source --config /xxx/config.josn --secret /xxx/secret.json
```

### k8s

```bash
kubectl apply -f https://raw.githubusercontent.com/linkall-labs/vance/main/connectors/xxx/xxx-bare.yml
```

## Configuration

### config.json

```json
{
  "v_target": "http://localhost:8080",
  "xxx": "xxx",
  "yyy": "yyy"
}
```

| Name     | Required | Default | Description                         |
| :--------- | :--------: | :-------: | ------------------------------------- |
| v_target | **YES** |    -    | Target URL will send CloudEvents to |
| xxx      | **YES** |    -    | xxxxx                               |
| yyy      |    NO    |  empty  |                                     |

xxxxxx

### secret.json

```json
{
  "xxx": "xxx",
  "yyy": "xxx",
  "zzz": "xxx"
}
```


| Name | Required | Default | Description |
| :----- | :--------: | :-------: | ------------- |
| xxx  | **YES** | babala | xxx         |
| yyy  |    NO    |  empty  | yyy         |
| zzz  |    NO    |  empty  | zzz         |

The `user` and `password` are required only when MongoDB is configured to use authentication. This `authsoure` required
only when MongoDB is configured to use authentication with another authentication database than admin.

## Schema

The output events' schema is a [CloudEvent](https://github.com/cloudevents/spec) format, and each field are explained
follows.


| field                    | description                         |
| -------------------------- | ------------------------------------- |
| id                       | the bson`_id` will be set as the id |
| source                   | xxx                                 |
| type                     | xxx                                 |
| time                     | xxx                                 |
| data                     | the body of`ChangeEvent`            |
| data.xxx                 | ...                                 |
| data.yyy                 | ...                                 |
| data.aaa.bbb             | ...                                 |
| vance{conenctorname}xxx} | ...                                 |
| vance{conenctorname}yyy  | ...                                 |

### Different Data Body explanation1(if it's needed)

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

### Different Data Body explanation2(if it's needed)

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

...

### Unrecognized Event(if it's needed)

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

xxxx

## Acknowledgement

The MongoDB Connector built on [debezium](https://github.com/debezium/debezium)
