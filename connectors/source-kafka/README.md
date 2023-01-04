---
title: Kafka
---

# Kafka Source

## Introduction

The Kafka Source is a [Vance Connector](https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md) which transforms Kafka messages from topics to CloudEvents and
deliver them to the target URL.

``` json
 { "name": "Jason", "age": "30" }
```

which is converted to a CloudEvent:

``` JSON
{
  "id" : "4ad0b59fc-3e1f-484d-8925-bd78aab15123",
  "source" : "kafka.localhost.topic2",
  "type" : "kafka.message",
  "datacontenttype" : "application/json or Plain/text",
  "time" : "2022-09-07T10:21:49.668Z",
  "data" : {
	 "name": "Jason",
	 "age": "30"
	 }
}
```

## Quick Start

in this section, we show how Kafka Source convert messages from topics to CloudEvents.

### Prerequisites
- Have a container runtime (i.e., [docker](https://www.docker.com)).
- Have a [kafka server](https://kafka.apache.org/quickstart) running.
- Have a or multiple topics created.


## Kafka Source Configs

```shell
cat << EOF > config.yml
v_target: 'http://<vanus_gateway_url><port>/gateway/<eventbus>'
KAFKA_SERVER_URL: YOUR_KAFKA_IP
KAFKA_SERVER_PORT: YOUR_KAFKA_PORT
CLIENT_ID: kafkaSource
TOPIC_LIST: YOUR_TOPIC
EOF
```
Config Fields of the kafka Source

| Configs   | Description                                                                     | Example                 |
|:----------|:--------------------------------------------------------------------------------|:------------------------|
| v_target  | v_target is used to specify the target URL HTTP Source will send CloudEvents to | "http://localhost:8081" |
| KAFKA_SERVER_URL    | The URL of the Kafka Cluster the Kafka Source is listening on                  | "8080"                  |
| KAFKA_SERVER_PORT    | v_port is used to specify the port Kafka Source is listening on                  | "8080"                  |
| CLIENT_ID    |  An optional identifier for multiple Kafka Sources that is passed to a Kafka broker with every request.                  | "kafkaSource"                  |
| TOPIC_LIST    | The source will listen to the topic or topics specified.                   | "topic1"  or "topic1, topic2, topic3"                 |

### start with Docker
```shell
docker run -d --rm \
  -v ${PWD}:/vance/config \
  --name source-kafka public.ecr.aws/vanus/connector/source-kafka:latest
```

### Test

From your kafka server run the following command to send messages:
> bin/kafka-console-producer.sh --topic YOUR_TOPIC --bootstrap-server YOUR_KAFKA_IP:PORT

now, you can use `vsctl event get <eventbus>` to view the event just sent. If you can't see the event you sent,
try to use `--offset` to get the event. (`vsctl` default retrieves event from earliest)

```
~> vsctl event get <eventbus>
+-----+----------------------------------------------------------------+
|     | Context Attributes,                                            |
|     |     {                                                          |
|     |       "id" : "ef26ed7b-9377-4bf5-b8d4-4fc6347e4fa2",           |
|     |       "source" : "kafka.host.docker.internal.topic1",          |
|     |       "specversion" : "V1",,                                   |
|     |       "type" : "kafka.message",                                |
|     |       "datacontenttype" : "plain/text",                        |
|     |       "time" : "2022-12-05T09:00:42.618Z",                     |
|     |       "data" : "Hello world!"                                  |
|     |     }                                                          |
|     |                                                                |
+-----+----------------------------------------------------------------+
```

### Clean

```shell
docker stop source-kafka
```

## How to use

### Configuration

The default path is `/vance/config/config.yml`. if you want to change the default path, you can set env `CONNECTOR_CONFIG` to
tell HTTP Source.


| Name                | Required |                       Default                       | Description                         |
|:--------------------|:--------:|:---------------------------------------------------:|-------------------------------------|
| v_target            | **YES**  | http://<vanus_gateway_url><port>/gateway/<eventbus> | the endpoint of CloudEvent sent to. |
| KAFKA_SERVER_URL    | **YES**  |                      localhost                      | The URL of the Kafka Cluster.       |
| KAFKA_SERVER_PORT   | **YES**  |                        9092                         | The PORT of the Kafka Cluster       |
| TOPIC_LIST          | **YES**  |                   topic1, topic2                    | The topic or topics to listen too.  |
| CLIENT_ID           |  **NO**  |                     KafkaSource                     | An optional identifier.             |

### Data

The ideal type of event for the Kafka source is a String in a JSON format. But it can handle any other type of data provided by Kafka.
> JSON Formatted String
> String = "{ "name": "Jason", "age": "30"}"
>

For example, if an original message looks like:
``` json
 { "name": "Jason", "age": "30" }
```

A Kafka message transformed into a CloudEvent looks like:

``` JSON
{
  "id" : "4ad0b59fc-3e1f-484d-8925-bd78aab15123",
  "source" : "kafka.localhost.topic2",
  "type" : "kafka.message",
  "datacontenttype" : "application/json or Plain/text",
  "time" : "2022-09-07T10:21:49.668Z",
  "data" : {
	 "name": "Jason",
	 "age": "30"
	 }
}
```

### Run in Kubernetes

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: source-kafka
  namespace: vanus
data:
  config.yml: |-
      KAFKA_SERVER_URL: "YOUR_KAFKA_IP"
      KAFKA_SERVER_PORT: "YOUR_KAFKA_PORT"
      CLIENT_ID: "kafkaSource"
      TOPIC_LIST: "YOUR_TOPIC"

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: source-kafka
  namespace: vanus
  labels:
    app: source-kafka
spec:
  selector:
    matchLabels:
      app: source-kafka
  replicas: 1
  template:
    metadata:
      labels:
        app: source-kafka
    spec:
      containers:
        - name: source-kafka
          image: public.ecr.aws/vanus/connector/source-kafka:latest
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
            name: source-kafka
```
