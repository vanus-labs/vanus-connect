---
title: Elasticsearch
---

# Elasticsearch Sink

## Introduction

The Elasticsearch Sink is a [Vanus Connector][vc] which aims to handle incoming CloudEvents in a way that extracts
the `data` part of the original event and deliver these extracted `data` to [Elasticsearch][es] cluster.

For example, an incoming CloudEvent looks like:

```json
{
  "specversion": "1.0",
  "id": "4395ffa3-f6de-443c-bf0e-bb9798d26a1d",
  "source": "vanus.source.test",
  "type": "vanus.type.test",
  "datacontenttype": "application/json",
  "time": "2022-06-14T07:05:55.777689Z",
  "data": {
    "id": 123,
    "date": "2022-06-13",
    "service": "test data"
  }
}
```

The Elasticsearch Sink will extract the `data` field and write it to the [Elasticsearch][es] cluster index as a
document:

```json
{
  "_index": "vanus_test",
  "_type": "_doc",
  "_id": "CqFnBIEBzJc0Oa5TERDD",
  "_version": 1,
  "_source": {
    "id": 123,
    "date": "2022-06-13",
    "service": "test data"
  }
}
```

## Quickstart

### Prerequisites

- Have a container runtime (i.e., docker).
- Have an Elasticsearch cluster.

### Create the config file

```shell
cat << EOF > config.yml
port: 8080
insert_mode: "upsert"
primary_key: "data.id"
es:
  address: "http://localhost:9200"
  index_name: "vanus_test"
  username: "elastic"
  password: "elastic"
EOF
```

| name         | requirement |     default     | description                                                   |
|:-------------|:-----------:|:---------------:|:--------------------------------------------------------------|
| port         |     NO      |      8080       | the port which Elasticsearch Sink listens on                  |
| address      |     YES     |                 | elasticsearch cluster address, multi split by ","             |
| index_name   |     YES     |                 | elasticsearch index name                                      |
| username     |     YES     |                 | elasticsearch cluster username                                |
| password     |     YES     |                 | elasticsearch cluster password                                |
| timeout      |     NO      |      10000      | elasticsearch index document timeout, unit millisecond        |
| insert_mode  |     NO      |     insert      | elasticsearch index document type: insert or upsert           |
| buffer_bytes |     NO      | 5 * 1024 * 1024 | elasticsearch each [bulk api][bulk api] request max body size |

The Elasticsearch Sink tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify
the position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm --network=host\
  -v ${PWD}:/vanus-connect/config \
  --name sink-elasticsearch public.ecr.aws/vanus/connector/sink-elasticsearch
```

### Test

Open a terminal and use the following command to send a CloudEvent to the Sink.

```shell
curl --location --request POST 'localhost:8080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
  "specversion": "1.0",
  "id": "4395ffa3-f6de-443c-bf0e-bb9798d26a1d",
  "source": "vanus.source.test",
  "type": "vanus.type.test",
  "datacontenttype": "application/json",
  "time": "2022-06-14T07:05:55.777689Z",
  "data": {
      "id": 123,
      "date": "2022-06-13",
      "service": "test data"
  }
}'
```

Use the following command to get an es document.

```shell
curl http://localhost:9200/vanus_test/_search?pretty
```

```json
{
  "_index": "vanus_test",
  "_type": "_doc",
  "_id": "123",
  "_version": 1,
  "_source": {
    "id": 123,
    "date": "2022-06-13",
    "service": "test data"
  }
}
```

### Clean resource

```shell
docker stop sink-elasticsearch
```

## Source details

### Extension Attributes

The Elasticsearch Sink have additional reactions if the incoming CloudEvent contains
following [Extension Attributes](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#extension-context-attributes)
.

| Attribute   | Required | Examples  | Description                                                                                                                                                                                |
|:------------|:--------:|:---------:|:-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| xvindexname |    NO    | "myindex" | the event want write to index name, if empty it will use config index_name                                                                                                                 |
| xvop        |    NO    |    "c"    | the event to document action:c, u, d; c convert to es index, u convert to es update with upsert, d convert to es delete, if config insert_mode is upsert the action c will be changed to u |
| xvid        |    NO    |   "123"   | the document id , it needs if action is delete, update                                                                                                                                     |

### Examples

#### Write to es with index

```shell
curl --location --request POST 'localhost:8080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
  "specversion": "1.0",
  "id": "4395ffa3-f6de-443c-bf0e-bb9798d26a1d",
  "source": "vanus.source.test",
  "type": "vanus.type.test",
  "datacontenttype": "application/json",
  "time": "2022-06-14T07:05:55.777689Z",
  "xvindexname": "myindex",
  "xvop": "c",
  "data": {
      "id": 123,
      "date": "2022-06-13",
      "service": "test data"
  }
}'
```

#### Write to es with document id

```shell
curl --location --request POST 'localhost:8080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
  "specversion": "1.0",
  "id": "4395ffa3-f6de-443c-bf0e-bb9798d26a1d",
  "source": "vanus.source.test",
  "type": "vanus.type.test",
  "datacontenttype": "application/json",
  "time": "2022-06-14T07:05:55.777689Z",
  "xvindexname": "myindex",
  "xvop": "c",
  "xvid": "123",
  "data": {
      "id": 123,
      "date": "2022-06-13",
      "service": "test data"
  }
}'
```

## Run in Kubernetes

```shell
kubectl apply -f sink-es.yaml
```

```yaml
apiVersion: v1
kind: Service
metadata:
  name: sink-es
  namespace: vanus
spec:
  selector:
    app: sink-es
  type: ClusterIP
  ports:
    - port: 8080
      name: sink-es
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: sink-es
  namespace: vanus
data:
  config.yml: |-
    port: 8080
    es:
      address: "http://localhost:9200"
      index_name: "vanus_test"
      username: "elastic"
      password: "elastic"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sink-es
  namespace: vanus
  labels:
    app: sink-es
spec:
  selector:
    matchLabels:
      app: sink-es
  replicas: 1
  template:
    metadata:
      labels:
        app: sink-es
    spec:
      containers:
        - name: sink-es
          image: public.ecr.aws/vanus/connector/sink-elasticsearch
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
            name: sink-es
```

## Integrate with Vanus

This section shows how a sink connector can receive CloudEvents from a
running [Vanus cluster](https://github.com/linkall-labs/vanus).

1. Run the sink-es.yaml

```shell
kubectl apply -f sink-es.yaml
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
  --sink 'http://sink-es:8080'
```

[vc]: https://www.vanus.ai/introduction/concepts#vanus-connect

[es]: https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html

[bulk api]: https://www.elastic.co/guide/en/elasticsearch/reference/master/docs-bulk.html
