---
title: HTTP
---

# HTTP Source

## Introduction

The HTTP Source is a [Vanus Connector][vc] which aims to convert an incoming HTTP Request to a CloudEvent.

For example, the incoming HTTP Request looks like:

```bash
curl --location --request POST 'localhost:8080/webhook?source=123&id=abc&type=456&subject=def&test=demo' \
--header 'Content-Type: text/plain' \
--data-raw '{
    "test":"demo"
}'
```

which is converted to:

```json
{
  "specversion": "1.0",
  "id": "abc",
  "source": "123",
  "type": "456",
  "subject": "def",
  "datacontenttype": "application/json",
  "time": "2023-01-29T03:25:26.229114Z",
  "data": {
    "body": {
      "test": "demo"
    },
    "headers": {
      "Accept": "*/*",
      "Content-Length": "21",
      "Content-Type": "text/plain",
      "Host": "localhost:8080",
      "User-Agent": "curl/7.85.0"
    },
    "method": "POST",
    "path": "/webhook",
    "query_args": {
      "id": "abc",
      "source": "123",
      "subject": "def",
      "test": "demo",
      "type": "456"
    }
  },
  "xvhttpremoteip": "::1",
  "xvhttpremoteaddr": "[::1]:57822",
  "xvhttpbodyisjson": true
}
```

## Quick Start

This section will show you how to use HTTP Source to convert an HTTP request(made by cURL) to a CloudEvent.

### Create Config file

```shell
cat << EOF > config.yml
target: http://localhost:31081
port: 8082
EOF
```

| Name   | Required | Default | Description                        |
| :----- | :------: | :-----: | :--------------------------------- |
| target |   YES    |         | the target URL to send CloudEvents |
| port   |    NO    |  8080   | the port to receive HTTP request   |

The HTTP Source tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-http public.ecr.aws/vanus/connector/source-http
```

### Test

Open a terminal and use the following command to run a Display sink, which receives and prints CloudEvents.

```shell
docker run -it --rm \
  -p 31081:8080 \
  --name sink-display public.ecr.aws/vanus/connector/sink-display
```

Make sure the `target` value in your config file is `http://localhost:31081` so that the Source can send the CloudEvents to the Display Sink.

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-http public.ecr.aws/vanus/connector/source-http
```

Open a terminal and use the following command to send an http request to HTTP Source

```shell
curl --location --request POST 'localhost:8082/webhook?source=123&id=abc&type=456&subject=def' \
--header 'Content-Type: text/plain' \
--data-raw '{
    "test":"demo"
}'
```

Here is the sort of CloudEvent you should expect to receive in the Display Sink:

```json
{
  "specversion": "1.0",
  "id": "abc",
  "source": "123",
  "type": "456",
  "subject": "def",
  "datacontenttype": "application/json",
  "time": "2023-01-29T03:25:26.229114Z",
  "data": {
    "body": {
      "test": "demo"
    },
    "headers": {
      "Accept": "*/*",
      "Content-Length": "21",
      "Content-Type": "text/plain",
      "Host": "localhost:8080",
      "User-Agent": "curl/7.85.0"
    },
    "method": "POST",
    "path": "/webhook",
    "query_args": {
      "id": "abc",
      "source": "123",
      "subject": "def",
      "type": "456"
    }
  },
  "xvhttpremoteip": "::1",
  "xvhttpremoteaddr": "[::1]:57822",
  "xvhttpbodyisjson": true
}
```

### Clean

```shell
docker stop source-http sink-display
```

## Source details

### Attributes

#### Changing Default Required Attributes

If you want to change the default attributes of `id`, `source`, `type`, and `subject`(defined by [CloudEvents](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#required-attributes)) to your own, you could use the `Query Parameter` to set them.

| Attribute  |      Default       | Query Parameter | Example                                 |
| :--------: | :----------------: | :-------------- | :-------------------------------------- |
|     id     |        UUID        | ?id=xxx         | http://url:port/webhook?id=xxxx         |
|   source   | vanus-http-source  | ?source=xxx     | http://url:port/webhook?source=xxxx     |
|    type    | naive-http-request | ?type=xxx       | http://url:port/webhook?type=xxxx       |
|  subject   |       empty        | ?subject=xxx    | http://url:port/webhook?subject=xxxx    |
| dataschema |       empty        | ?dataschema=xxx | http://url:port/webhook?dataschema=xxxx |

`datacontenttype` will be automatically inferred based on the request body. If the body can be converted to `JSON`, the `application/json` will be set. Otherwise, `text/plain` will be set.

#### Extension Attributes

The HTTP Source defines following [CloudEvents Extension Attributes](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#extension-context-attributes)

|    Attribute     |  Type   | Description                                                                                                                      |
| :--------------: | :-----: | :------------------------------------------------------------------------------------------------------------------------------- |
| xvhttpbodyisjson | boolean | HTTP Sink will validate if request body is JSON format data, if it is, this attribute is `true`, otherwise `false`               |
|  xvhttpremoteip  | string  | The IP of the request from where, if the request was through reverse-proxy like Nginx, the value may be not the original IP      |
| xvhttpremoteaddr | string  | The address of the request from where, if the request was through reverse-proxy like Nginx, the value may be not the original IP |

## Run in Kubernetes

```shell
kubectl apply -f source-http.yaml
```

```yaml
apiVersion: v1
kind: Service
metadata:
  name: source-http
  namespace: vanus
spec:
  selector:
    app: source-http
  type: ClusterIP
  ports:
    - port: 8080
      name: source-http
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: source-http
  namespace: vanus
data:
  config.yml: |-
    target: http://<url>:<port>/gateway/<eventbus>

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: source-http
  namespace: vanus
  labels:
    app: source-http
spec:
  selector:
    matchLabels:
      app: source-http
  replicas: 1
  template:
    metadata:
      labels:
        app: source-http
    spec:
      containers:
        - name: source-http
          image: public.ecr.aws/vanus/connector/source-http:latest
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "512Mi"
              cpu: "500m"
          imagePullPolicy: Always
          volumeMounts:
            - name: config
              mountPath: /vanus-connect/config
      volumes:
        - name: config
          configMap:
            name: source-http
```

## Integrate with Vanus

This section shows how a source connector can send CloudEvents to a running [Vanus cluster](https://github.com/linkall-labs/vanus).

### Prerequisites

- Have a running K8s cluster
- Have a running Vanus cluster
- Vsctl Installed

1. Export the VANUS_GATEWAY environment variable (the ip should be a host-accessible address of the vanus-gateway service)

```shell
export VANUS_GATEWAY=192.168.49.2:30001
```

2. Create an eventbus

```shell
vsctl eventbus create --name quick-start
```

3. Update the target config of the HTTP Source

```yaml
target: http://192.168.49.2:30001/gateway/quick-start
```

4. Run the HTTP Source

```shell
kubectl apply -f source-http.yaml
```

[vc]: https://www.vanus.ai/introduction/concepts#vanus-connect
