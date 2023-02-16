---
title: Kafka
---

# Kafka Source

## Overview

The Kafka Source is a [Vanus Connector][vc] which aims to consume Kafka messages from topics, converts them into CloudEvents, and then deliver them to the target URL.

For example, if an original message looks like:

```text
> { "name": "Jason", "age": "30" }
```

It will be converted to CloudEvent this way:

```JSON
{
  "id" : "4ad0b59fc-3e1f-484d-8925-bd78aab15123",
  "source" : "kafka_bootstrap_servers.mytopic",
  "type" : "kafka.message",
  "datacontenttype" : "application/json",
  "time" : "2022-09-07T10:21:49.668Z",
  "data" : {
	 "name": "Jason",
	 "age": "30"
  }
}
```

## Quick Start

This section will show you how to use Kafka Source to convert Kafka messages to CloudEvents.

### Prerequisites

- Have a container runtime (i.e., docker).
- Have a [Kafka cluster](https://kafka.apache.org)

### Create the config file

```shell
cat << EOF > config.yml
target: http://localhost:31081
bootstrap_servers: "localhost:9092"
group_id: "vanus-source-kafka"
topics: [ "mytopic" ]
EOF
```

| Name              | Required | Default | Description                                                |
| :---------------- | :------- | :-----: | :--------------------------------------------------------- |
| target            | YES      |         | the target URL which Kafka Source will send CloudEvents to |
| bootstrap_servers | YES      |         | the kafka cluster bootstrap servers                        |
| group_id          | YES      |         | the kafka cluster consumer group id                        |
| topics            | YES      |         | the kafka topics listened by kafka source                  |

The Kafka Source tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the
position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-kafka public.ecr.aws/vanus/connector/source-kafka
```

### Test

Open a terminal and use the following command to run a Display sink, which receives and prints CloudEvents.

```shell
docker run -it --rm \
  -p 31081:8080 \
  --name sink-display public.ecr.aws/vanus/connector/sink-display
```

Make sure the `target` value in your config file is `http://localhost:31081` so that the Source can send CloudEvents to our Display Sink.

After running Display Sink, run the Kafka Source

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-kafka public.ecr.aws/vanus/connector/source-kafka
```

Send kafka message use the following command:

```shell
bin/kafka-console-producer.sh --topic mytopic --bootstrap-server localhost:9092
```

Here is the sort of CloudEvent you should expect to receive in the Display Sink:

```json
{
  "id": "4ad0b59fc-3e1f-484d-8925-bd78aab15123",
  "source": "kafka_bootstrap_servers.mytopic",
  "type": "kafka.message",
  "datacontenttype": "application/json",
  "time": "2022-09-07T10:21:49.668Z",
  "data": {
    "name": "Jason",
    "age": "30"
  }
}
```

### Clean

```shell
docker stop source-kafka sink-display
```

## Run in Kubernetes

```shell
kubectl apply -f source-kafka.yaml
```

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: source-kafka
  namespace: vanus
data:
  config.yaml: |-
    target: "http://vanus-gateway.vanus:8080/gateway/quick_start"
    bootstrap_servers: "localhost:9092"
    group_id: "vanus-source-kafka"
    topics: [ "mytopic" ]
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
          image: public.ecr.aws/vanus/connector/source-kafka
          imagePullPolicy: Always
          volumeMounts:
            - name: config
              mountPath: /vanus-connect/config
      volumes:
        - name: config
          configMap:
            name: source-kafka
```

## Integrate with Vanus

This section shows how a source connector can send CloudEvents to a
running [Vanus cluster](https://github.com/linkall-labs/vanus).

### Prerequisites

- Have a running K8s cluster
- Have a running Vanus cluster
- Vsctl Installed

1. Export the VANUS_GATEWAY environment variable (the ip should be a host-accessible address of the vanus-gateway
   service)

```shell
export VANUS_GATEWAY=192.168.49.2:30001
```

2. Create an eventbus

```shell
vsctl eventbus create --name quick-start
```

3. Update the target config of the Kafka Source

```yaml
target: http://192.168.49.2:30001/gateway/quick-start
```

4. Run the Kafka Source

```shell
kubectl apply -f source-kafka.yaml
```

[vc]: https://www.vanus.ai/introduction/concepts#vanus-connect
