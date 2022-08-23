# MongoDB Connector

## Quickstart

### Docker

```bash
docker run -it --rm public.ecr.aws/vanus/connector/mongodb:latest /run/start.sh \
  --volume /xxx/secret.json /var/mongodb/secret.json \
  --env MONGODB_HOSTS=xxx \
  --env MONGODB_NAME=xxx \
  --env MONGODB_AUTHSOURCE=xxx
```

## Schema

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
    "id": "6304855bccaea8fcf8a159f2",
    "full": {
      "download": "1234",
      "connector": "mongodb",
      "_id": "6304855bccaea8fcf8a159f2",
      "version": "v0.3.0"
    }
  },
  "vancemongodbformatted": true,
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
    "id": "6304855bccaea8fcf8a159f2",
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
  "vancemongodbformatted": true,
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
  "data": {
    "id": "6304855bccaea8fcf8a159f2"
  },
  "vancemongodbord": "1",
  "vancemongodbformatted": true,
  "vancemongodbversion": "1.9.4.Final",
  "vancemongodbsnapshot": "false",
  "vancemongodbname": "test",
  "vancemongodboperation": "delete"
}
```

### Unrecognized Event

```json
    "specversion":"1.0",
"id": "unknown",
"source": "unknown",
"type": "unknown",
"datacontenttype": "application/json",
"time":"unknown",
"data": {
"rawKey": "xxxxx",
"rawValue": "xxxx",
},
```

## Acknowledgement

### k8s

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mongodb-connector
  labels:
    app: mongodb-connector
spec:
  selector:
    matchLabels:
      app: mongodb-connector
  replicas: 1
  template:
    metadata:
      labels:
        app: mongodb-connector
    spec:
      containers:
        - name: mongodb-connector
          image: public.ecr.aws/vanus/connector/mongodb:latest
          imagePullPolicy: Always
          command: [ "sh", "-c", "/var/mongodb/start.sh" ]
          resources:
            requests:
              cpu: 100m
              memory: 1000Mi
          env:
            - name: MONGODB_HOSTS
              value: "localhost:27017"
            - name: MONGODB_NAME
              value: "admin"
            - name: MONGODB_AUTHSOURCE
              value: "admin"
            - name: DB_INCLUDE_LIST
              value: "test"
          volumeMounts:
            - name: secret
              mountPath: "/var/mongodb/secret.json"
              readOnly: true
        volumes:
          - name: secret
            secret:
              secretName: mongodb-secret
---
apiVersion: v1
kind: Secret
metadata:
  name: mongodb-secret
type: Opaque
data:
  user: admin
  password: admin
```