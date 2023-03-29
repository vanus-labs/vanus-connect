---
title: ChatGPT
---

# ChatGPT Source

## Introduction

The ChatGPT Source is a [Vanus Connector][vc] which aims to read an incoming request body, then call openai api and
convert the response content to a CloudEvent.

For example, the incoming Request looks like:

```bash
curl --location --request POST 'localhost:8080' \
--header '' \
--data-raw 'what is vanus'
```

which is converted to:

```json
{
  "specversion": "1.0",
  "id": "0effe4cc-06c7-4fe9-9180-aa7c3b30777e",
  "source": "vanus-chatGPT-source",
  "type": "vanus-chatGPT-type",
  "datacontenttype": "application/json",
  "time": "2023-03-28T09:15:10.70413Z",
  "data": {
    "content": "vanus is a message queue"
  }
}
```

## Quick Start

This section will show you how to use ChatGPT Source to read request body and call openai api to obtains response then
convert to a CloudEvent.

### Create Config file

```shell
cat << EOF > config.yml
target: http://localhost:31081
port: 8082
token: xxxxxx
EOF
```

| Name           | Required | Default | Description                                        |
|:---------------|:--------:|:--------|:---------------------------------------------------|
| target         |   YES    |         | the target URL to send CloudEvents                 |
| port           |    NO    | 8080    | the port to receive HTTP request                   |
| token          |   YES    |         | the ChatGPT auth token                             |
| everyday_limit |    NO    | 100     | the ChatGPT Source call openapi api count everyday |

The ChatGPT Source tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the
position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-chatgpt public.ecr.aws/vanus/connector/source-chatgpt
```

### Test

Open a terminal and use the following command to run a Display sink, which receives and prints CloudEvents.

```shell
docker run -it --rm \
  -p 31081:8080 \
  --name sink-display public.ecr.aws/vanus/connector/sink-display
```

Make sure the `target` value in your config file is `http://localhost:31081` so that the Source can send the CloudEvents
to the Display Sink.

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-chatgpt public.ecr.aws/vanus/connector/source-chatgpt
```

Open a terminal and use the following command to send a request to ChatGPT Source

```shell
curl --location --request POST 'localhost:8082' \
--header 'Content-Type: text/plain' \
--data-raw 'what is vanus'
```

Here is the sort of CloudEvent you should expect to receive in the Display Sink:

```json
{
  "specversion": "1.0",
  "id": "0effe4cc-06c7-4fe9-9180-aa7c3b30777e",
  "source": "vanus-chatGPT-source",
  "type": "vanus-chatGPT-type",
  "datacontenttype": "application/json",
  "time": "2023-03-28T09:15:10.70413Z",
  "data": {
    "content": "vanus is a message queue"
  }
}
```

### Clean

```shell
docker stop source-chatgpt sink-display
```

## Source details

### Attributes

#### Changing Default Required Attributes

If you want to change the default attributes of `source`, `type` (defined
by [CloudEvents](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#required-attributes)) to your own,
you could set the request header to set them.

| Header        | Description       |
|:--------------|:------------------|
| vanus-source  | cloudevent source |
| vanus-type    | coudevent type    |

## Run in Kubernetes

```shell
kubectl apply -f source-chatgpt.yaml
```

```yaml
apiVersion: v1
kind: Service
metadata:
  name: source-chatgpt
  namespace: vanus
spec:
  selector:
    app: source-chatgpt
  type: ClusterIP
  ports:
    - port: 8080
      name: source-chatgpt
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: source-chatgpt
  namespace: vanus
data:
  config.yml: |-
    target: "http://localhost:18080"
    token: "sk-k7UNuxZiZZVOYEU8xxxxxxxxxxxxxxxxx"

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: source-chatgpt
  namespace: vanus
  labels:
    app: source-chatgpt
spec:
  selector:
    matchLabels:
      app: source-chatgpt
  replicas: 1
  template:
    metadata:
      labels:
        app: source-chatgpt
    spec:
      containers:
        - name: source-chatgpt
          image: public.ecr.aws/vanus/connector/source-chatgpt:latest
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
            name: source-chatgpt
```

## Integrate with Vanus

This section shows how a source connector can send CloudEvents to a
running [Vanus cluster](https://github.com/vanus-labs/vanus).

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

3. Update the target config of the ChatGPT Source

```yaml
target: http://192.168.49.2:30001/gateway/quick-start
```

4. Run the ChatGPT Source

```shell
kubectl apply -f source-chatgpt.yaml
```

[vc]: https://docs.vanus.ai/introduction/concepts#vanus-connect
