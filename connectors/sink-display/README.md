# Display Sink 

## Overview

A [Vance Connector][vc] which prints received CloudEvents. This is commonly used as a logger to check incoming data.

## Introduction

The Display Sink is a single function [Connector][vc] which aims to print incoming CloudEvents in JSON format.

For example, it will print the incoming CloudEvent looks like:

```http
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

### Start Using Docker

mapping 8080 to 31080 in order to avoid port conflict.

```shell
docker run -d -p 31080:8080 --rm \
  -v ${PWD}:/vance/config \
  --name sink-display public.ecr.aws/vanus/connector/sink-display:latest
```

### Test
1. make a HTTP request
```shell
curl --location --request POST 'localhost:8080' \
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

2. view logs
```shell
docker logs sink-display
```

```shell
receive a new event, in total: 1
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

## How to use

### Run in Kubernetes
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
      initContainers:
      - name: touc-config
        image: busybox
        command:
          - touch
          - "${CONNECTOR_CONFIG}"
        volumeMounts:
          - name: config
            mountPath: /vance/config
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
              cpu: "100m"
          env:
            - name: LOG_LEVEL
              value: INFO
          volumeMounts:
            - name: config
              mountPath: /vance/config
      volumes:
        - name: config
          emptyDir: {}
```

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md