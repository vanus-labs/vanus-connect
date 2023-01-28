---
title: Slack
---

# Slack Sink

## Introduction

The Slack Sink is a [Vanus Connector][vc] which aims to handle incoming CloudEvents in a way that extracts the `data` part of the
original event and deliver these extracted `data` to Slack channels.

For example, if the incoming CloudEvent looks like:

```json
{
  "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
  "source": "quick-start",
  "specversion": "1.0",
  "type": "quick-start",
  "datacontenttype": "application/json",
  "time": "2022-10-26T10:38:29.345Z",
  "data": {
    "subject": "Test",
    "message": "Hello Slack!:wave: This is Sink Slack!"
  }
}
```

then channels will receive a message like:
![message](https://github.com/linkall-labs/vance/blob/main/connectors/sink-slack/message.png?raw=true)

## Quick Start

in this section, we show how to use Slack Sink sends a text message to recipients.

### Prerequisites
- Have a container runtime (i.e., docker).
- Have a Slack App and should have at least `chat:write` and `chat:write.public` permission.

### Create Config file

```shell
cat << EOF > config.yml
default: "test_app"
slack:
  - app_name: "test_app"
    token: "<oauth token>"
    default_channel: "#general"
EOF
```

| Name                     | Required  | Default | Description                                                                                           |
|:-------------------------|:---------:|:--------|:------------------------------------------------------------------------------------------------------|
| port                     |    NO     | 8080    | the port which <name> Sink listens on                                                                 |
| default                  |    YES    |         | the default app name if event attribute doesn't have `xvslackapp`                                     |
| slack.[].app_name        |    YES    |         | custom slack app name as identifier                                                                   |
| slack.[].token           |    YES    |         | OAuth Token of this app, more visit: https://api.slack.com/legacy/oauth                               |
| slack.[].default_channel |    NO     |         | set default channel the messages send to if attribute was not be set, use `,` to separate multiples   |

The Slack Sink tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm \
  -p 31080:8080 \
  -v ${PWD}:/vanus-connect/config \
  --name sink-slack public.ecr.aws/vanus/connector/sink-slack
```

### Test

Open a terminal and use following command to send a CloudEvent to the Sink.

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "quick-start",
    "specversion": "1.0",
    "type": "quick-start",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "data": {
        "subject": "Test",
        "message": "Hello Slack!:wave: This is Sink Slack!"
    }
}'
```

now, you cloud see a new slack in your mailbox.
![message.png](https://github.com/linkall-labs/vance/blob/main/connectors/sink-slack/message.png?raw=true)

### Clean

```shell
docker stop sink-slack
```

## Sink details

### Extension Attributes

The <name> Sink have additional reactions if the incoming CloudEvent contains following[Extension Attributes](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#extension-context-attributes).

| Attribute       | Required  | Examples          | Description                                 |
|:----------------|:---------:|:------------------|:--------------------------------------------|
| xvslackapp      |    NO     | test_app          | Which slack app this event want to send to  |
| xvslackchannels |    NO     | #general,#random  | use `,` to separate multiples               |

### Data format

the event data must be `JSON` format, and only two key `subject` and `message` is valid for using, example:

```json
{
  "subject": "Test",
  "message": "Hello Slack!:wave: This is Sink Slack!"
}
```

### Examples

#### Sending message to the default app and default channel

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "quick-start",
    "specversion": "1.0",
    "type": "quick-start",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "data": {
        "subject": "Test",
        "message": "Hello Slack!:wave: This is Sink Slack!"
    }
}'
```

#### Sending message to the specified app and specified channel

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "quick-start",
    "specversion": "1.0",
    "type": "quick-start",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "xvslackapp": "test",
    "xvslackchannels": "#team-a",
    "data": {
        "subject": "Test",
        "message": "Hello Slack!:wave: This is Sink Slack!"
    }
}'
```

#### Sending message to the specified app and multiple channels

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "quick-start",
    "specversion": "1.0",
    "type": "quick-start",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "xvslackapp": "test",
    "xvslackchannels": "#team-a,#team-b,#team-c",
    "data": {
        "subject": "Test",
        "message": "Hello Slack!:wave: This is Sink Slack!"
    }
}'
```

## Run in Kubernetes

```yaml
apiVersion: v1
kind: Service
metadata:
  name: sink-slack
  namespace: vanus
spec:
  selector:
    app: sink-slack
  type: ClusterIP
  ports:
    - port: 8080
      name: sink-slack
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: sink-slack
  namespace: vanus
data:
  config.yml: |-
    default: "test-app1"
    slack:
      - app_name: "test-app1"
        token: "xoxb-xxxxxxxxxx"
        default_channel: "#general"
      - app_name: "test-app2"
        token: "xoxb-xxxxxxxxxx"
        default_channel: "#general"

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sink-slack
  namespace: vanus
  labels:
    app: sink-slack
spec:
  selector:
    matchLabels:
      app: sink-slack
  replicas: 1
  template:
    metadata:
      labels:
        app: sink-slack
    spec:
      containers:
        - name: sink-slack
          image: public.ecr.aws/vanus/connector/sink-slack
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "512Mi"
              cpu: "500m"
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
            name: sink-slack
```

## Integrate with Vanus

This section shows how a sink connector can receive CloudEvents from a running [Vanus cluster](https://github.com/linkall-labs/vanus).

1. Run the sink-slack.yaml
```shell
kubectl apply -f sink-slack.yaml
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
  --sink 'http://sink-name:8080'
```

[vc]: https://www.vanus.dev/introduction/concepts#vanus-connect
