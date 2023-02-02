---
title: HTTP
---

# HTTP Sink

## Introduction

The HTTP Sink is a [Vanus Connector][vc] which aims to handle incoming CloudEvents in a way that extracts the `data`
part of the original event and deliver to the target URL.

For example, if the incoming CloudEvent looks like:

```http
{
  "id" : "42d5b039-daef-4071-8584-e61df8fc1354",
  "source" : "quickstart",
  "specversion" : "1.0",
  "type" : "quickstart",
  "datacontenttype" : "application/json",
  "time" : "2023-01-26T10:38:29.345Z",
  "data" : {
    "headers":{
        "connect-name": "sink-http"
    },
    "query": "debug=true&type=curl",
    "body" : "simulation event data 1"
  }
}
```

The HTTP Sink will send an HTTP request looks like:

```text
POST /test?debug=true&type=curl

> Host: localhost:8081
> User-Agent: 	Go-http-client/1.1
> Content-Length: 23
> connect-name: sink-http

> simulation event data 1
```

## Quickstart

### Prerequisites

- Have a container runtime (i.e., docker).
- Have an HTTP server, you can go https://webhook.site to get a free URL

### Create the config file


```shell
cat << EOF > config.yml
target: <you server endpoint>
EOF
```

| Name            | Required | Default | Description                                                                      |
|:----------------|:--------:|:--------|----------------------------------------------------------------------------------|
| port            |    NO    | 8080    | the port which HTTP Sink listens on                                              |
| target          |   YES    |         | the target which HTTP Sink send http request, example: http://xxxxxx:8081/xxxxxx |
| method          |    NO    | POST    | the default http request method                                                  |
| headers         |    NO    |         | the default http request headers                                                 |
| auth.username   |    NO    |         | if your http server authentication by basic auth, username is needed             |
| auth.password   |    NO    |         | if your http server authentication by basic auth, password is needed             |

The HTTP Sink tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the
position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm \
  -p 31080:8080 \
  -v ${PWD}:/vanus-connect/config \
  --name sink-http public.ecr.aws/vanus/connector/sink-http
```

### Test

Open a terminal and use following command to send a CloudEvent to the Sink.

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
  "id" : "42d5b039-daef-4071-8584-e61df8fc1354",
  "source" : "quickstart",
  "specversion" : "1.0",
  "type" : "quickstart",
  "datacontenttype" : "application/json",
  "time" : "2023-01-26T10:38:29.345Z",
  "data" : {
    "headers":{
        "connect-name": "sink-http"
    },
    "query": "debug=true&type=curl",
    "body" : "simulation event data 1"
  }
}'
```

Then your HTTP server will receive the request.

### Clean resource

```shell
docker stop sink-http
```

## Sink details

### Data format

The HTTP Sink requires following data format in CloudEvent's `data` field.

```json
{
  "method": "POST",
  "path": "xxxxxx/xxxxxx",
  "headers": {
    "connect-name": "sink-http"
  },
  "query": "debug=true&type=curl",
  "body": "simulation event data 1"
}
```

## Run in Kubernetes

```shell
kubectl apply -f sink-http.yaml
```

```yaml
apiVersion: v1
kind: Service
metadata:
  name: sink-http
  namespace: vanus
spec:
  selector:
    app: sink-http
  type: ClusterIP
  ports:
    - port: 8080
      name: sink-http
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: sink-http
  namespace: vanus
data:
  config.yml: |-
    port: 8080
    target: http://vanus-gateway.vanus:8080/gateway/quick_start
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sink-http
  namespace: vanus
  labels:
    app: sink-http
spec:
  selector:
    matchLabels:
      app: sink-http
  replicas: 1
  template:
    metadata:
      labels:
        app: sink-http
    spec:
      containers:
        - name: sink-http
          image: public.ecr.aws/vanus/connector/sink-http:latest
          imagePullPolicy: Always
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "128Mi"
              cpu: "100m"
          ports:
            - name: http
              containerPort: 8080
          volumeMounts:
            - name: config
              mountPath: /vanus-connect/config
      volumes:
        - name: config
          configMap:
            name: sink-http
```

## Integrate with Vanus

This section shows how a sink connector can receive CloudEvents from a
running [Vanus cluster](https://github.com/linkall-labs/vanus).

1. Run the sink-http.yaml

```shell
kubectl apply -f sink-http.yaml
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
  --sink 'http://sink-http:8080'
```

[vc]: https://www.vanus.dev/introduction/concepts#vanus-connect
