---
title: Display
---

# Display Sink 

## Introduction

The Display Sink is a [Vanus Connector](https://www.vanus.dev/introduction/concepts#vanus-connect) which aims to print incoming CloudEvents in JSON format.

For example, it will print the incoming CloudEvent looks like:

```json
{
  "specversion": "1.0",
  "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
  "source": "quickstart",
  "type": "quickstart",
  "datacontenttype": "application/json",
  "time": "2022-10-26T10:38:29.345Z",
  "data": {
    "myData": "simulation event data"
  }
}
```

## Quick Start

### Start with Docker

```shell
docker run -it --rm \
-p 31080:8080 \
--name sink-display public.ecr.aws/vanus/connector/sink-display
```

### Test

Open a terminal and use following command to send a CloudEvent to the Sink.

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "quickstart",
    "specversion": "1.0",
    "type": "quickstart",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "data": {
        "myData": "simulation event data"
    }
}'
```

The Display Sink will print:

```shell
INFO[2022-10-26T06:22:41.754221044Z] logger level is set  log_level=INFO
INFO[2022-10-26T06:22:41.849166961Z] the connector started  connector-name="Display Sink" listening=8080
INFO[2022-10-26T03:25:26.262083591Z] receive a new event  in_total=1
{
  "specversion": "1.0",
  "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
  "source": "quickstart",
  "type": "quickstart",
  "datacontenttype": "application/json",
  "time": "2022-10-26T10:38:29.345Z",
  "data": {
    "myData": "simulation event data"
  }
}
```
### Clean resource

```shell
docker stop sink-display
```

## Run in Kubernetes

```shell
kubectl apply -f sink-display.yaml
```

```yaml
apiVersion: v1
kind: Service
metadata:
  name: sink-display
  namespace: vanus
spec:
  selector:
    app: sink-display
  type: ClusterIP
  ports:
    - port: 8080
      name: sink-display
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sink-display
  namespace: vanus
  labels:
    app: sink-display
spec:
  selector:
    matchLabels:
      app: sink-display
  replicas: 1
  template:
    metadata:
      labels:
        app: sink-display
    spec:
      containers:
        - name: sink-display
          image: public.ecr.aws/vanus/connector/sink-display:latest
          imagePullPolicy: Always
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "128Mi"
```

## Integrate with Vanus

This section shows how a sink connector can receive CloudEvents from a running [Vanus cluster](https://github.com/linkall-labs/vanus).

1. Run the sink-display.yaml
```shell
kubectl apply -f sink-display.yaml
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
  --sink 'http://sink-display:8080'
```

[vc]: https://www.vanus.dev/introduction/concepts#vanus-connect
