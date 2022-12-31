---
title: Slack
---

# Slack Sink

## Introduction

The Slack Sink is a [Vance Connector][vc] which aims to handle incoming CloudEvents in a way that extracts the `data` part of the
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

### Create Config file

replace `<app_id>`, `<custom app name>`, and `<oauth token>` to yours. If you haven't a Slack App, you could
create one by following [Create your Slack App](https://api.slack.com/apps/new), the App should have
at least `chat:write` and `chat:write.public` permission.

```shell
cat << EOF > config.yml
default: "<app_id>"
slack:
  - app_name: "<custom app name>"
    token: "<oauth token>"
    default_channel: "#general"
EOF
```

### Start Using Docker

mapping 8080 to 31080 in order to avoid port conflict.

```shell
docker run -d -p 31080:8080 --rm \
  -v ${PWD}:/vance/config \
  --name sink-slack public.ecr.aws/vanus/connector/sink-slack:latest
```

### Test

replace `<from_slack_address>`, `<from_slack_address>`, and `<smtp server address>` to yours.

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

## How to use

### Configuration

The default path is `/vance/config/config.yml`. if you want to change the default path, you can set env `CONNECTOR_CONFIG` to
tell Slack Sink.

| Name                     | Required | Default | Description                                                                                         |
| :----------------------- | :------: | :-----: | --------------------------------------------------------------------------------------------------- |
| default                  | **YES**  |    -    | Slack Sink supports multiple slack apps as target, you could set the default app by this field      |
| slack.[].app_name        | **YES**  |    -    | custom slack app name as identifier                                                                 |
| slack.[].token           | **YES**  |    -    | OAuth Token of this app, more visit: https://api.slack.com/legacy/oauth                             |
| slack.[].default_channel |    NO    |    -    | set default channel the messages send to if attribute was not be set, use `,` to separate multiples |

### Extension Attributes

Slack Sink has defined a few [CloudEvents Extension Attribute](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#extension-context-attributes)
to determine how to process event.

| Attribute       | Required | Examples         | Description                                |
| :-------------- | :------: | ---------------- | ------------------------------------------ |
| xvslackapp      |    NO    | test_app         | Which slack app this event want to send to |
| xvslackchannels |    NO    | #general,#random | use `,` to separate multiples              |

### Data

the event data must be `JSON` format, and only two key `subject` and `message` is valid for using, example:

```json
{
  "subject": "Test",
  "message": "Hello Slack!:wave: This is Sink Slack!"
}
```

## Examples

### Sending message to the default app and default channel

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

### Sending message to the specified app and specified channel

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

### Sending message to the specified app and multiple channels

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
    default: "test-app"
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
          #          For China mainland
          #          image: linkall.tencentcloudcr.com/vanus/connector/sink-slack:latest
          image: public.ecr.aws/vanus/connector/sink-slack:latest
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
              mountPath: /vance/config
      volumes:
        - name: config
          configMap:
            name: sink-slack
```

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md
