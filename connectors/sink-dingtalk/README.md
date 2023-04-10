---
title: Dingtalk
---

# Dingtalk Sink

## Introduction

The Dingtalk Sink is a [Vanus Connector][vc] which aims to handle incoming CloudEvents in a way that extracts the `data`
part of the
original event and deliver these extracted `data` to the Dingtalk APIs. Now the Sink support Dingtalk Bot: pushing a
message to Group Chat with message of text, link, markdown, actionCard, feedCard.

## Quick Start

In this section, we show how to use Dingtalk Sink push a text message to your group chat.

### Prerequisites

#### Add a bot to your group chat

Go to your target group, click Group Settings > Bots > Add RoBot, and select Custom to add the bot to the group chat.
You can refer [docs][robot].
You will get the webhook address of the bot in the following format:

```
https://oapi.dingtalk.com/robot/send?access_token=xxxxxx
```

> ⚠️ You must set your signature verification to make sure push messages work.

### Create the config file

Replace `chat_group`, `signature`, and `url` to yours. `chat_group` can be fill in any value as you want.

```shell
cat << EOF > config.yml
bot:
  webhooks:
    - chat_group: test
      url: https://oapi.dingtalk.com/robot/send?access_token=xxxxxx
      signature: SECxxxxxx
EOF
```

| Name                       | Required | Default | Description                                                                       |
|:---------------------------|:--------:|:-------:|-----------------------------------------------------------------------------------|
| bot.default                |  **NO**  |    -    | default chat group, if not set it will use webhooks first element chat_group      |
| bot.webhooks.[].chat_group | **YES**  |    -    | the chat_group name, you can set any value to it                                  |
| bot.webhooks.[].signature  | **YES**  |    -    | the signature to sign request, you can get it when you create Chat Bot            |
| bot.webhooks.[].url        | **YES**  |    -    | the webhook address that message sent to, you can get it when you create Chat Bot |

The Dingtalk Sink tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the
position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm \
  -p 31080:8080 \
  -v ${PWD}:/vanus-connect/config \
  --name sink-Dingtalk public.ecr.aws/vanus/connector/sink-dingtalk
```

### Test

Open a terminal and use the following command to send a CloudEvent to the Sink.

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "sink-Dingtalk-quickstart",
    "specversion": "1.0",
    "type": "quickstart",
    "datacontenttype": "text/plain",
    "time": "2022-10-26T10:38:29.345Z",
    "data": "Hello dingtalk"
}'
```

now, you can see a notification from your bot in your group chat.

### Clean

```shell
docker stop sink-dingtalk
```

## Sink details

### Extension Attributes

Dingtalk Sink has defined a
few [CloudEvents Extension Attribute](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#extension-context-attributes)
to determine how to process event.

| Attribute   | Required | Examples | Description                                                                                                   |
|:------------|:--------:|----------|---------------------------------------------------------------------------------------------------------------|
| xvchatgroup |    NO    | text     | which Dingtalk chat-group the event sent for, the default value is config default                             |
| xvmsgtype   |    NO    | text     | which Message Type the event convert to, default is text, support: text, link, markdown, actionCard, feedCard |


## Examples

### Dingtalk Bot

you could find official docs of Dingtalk [robot][robot]

#### Text Message

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "sink-Dingtalk-quickstart",
    "specversion": "1.0",
    "type": "quickstart",
    "datacontenttype": "text/plain",
    "time": "2022-10-26T10:38:29.345Z",
    "xvchatgroup": "bot1",
    "xvmsgtype": "text",
    "data": "Hello dingtalk"
}'
```

#### Link Message

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "sink-dingtalk-quickstart",
    "specversion": "1.0",
    "type": "quickstart",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "xvmsgtype": "link",
    "data": {
        "text": "text", 
        "title": "title", 
        "picUrl": "https://www.vanus.ai/images/ng/vanus-black.svg", 
        "messageUrl": "https://github.com/vanus-labs"
    }
}'
```

#### Markdown Message

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "sink-dingtalk-quickstart",
    "specversion": "1.0",
    "type": "quickstart",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "xvmsgtype": "markdown",
    "data": {
        "text": "text", 
        "title": "title"
    }
}'
```

#### ActionCard Message

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "sink-dingtalk-quickstart",
    "specversion": "1.0",
    "type": "quickstart",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "xvmsgtype": "actionCard",
    "data": {
        "text": "text", 
        "title": "title",
        "btnOrientation": "0", 
        "btns": [
            {
                "title": "title1", 
                "actionURL": "https://github.com/vanus-labs"
            }, 
            {
                "title": "title2", 
                "actionURL": "https://github.com/vanus-labs"
            }
        ]
    }
}'
```


#### FeedCard Message

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "sink-dingtalk-quickstart",
    "specversion": "1.0",
    "type": "quickstart",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "xvmsgtype": "feedCard",
    "data": {
        "links": [
            {
                "title": "title1", 
                "messageURL": "https://github.com/vanus-labs", 
                "picURL": "https://www.vanus.ai/images/ng/vanus-black.svg"
            },
            {
                "title": "title2", 
                "messageURL": "https://github.com/vanus-labs", 
                "picURL": "https://www.vanus.ai/images/ng/vanus-black.svg"
            }
        ]
    }
}'
```


## Run in Kubernetes

```shell
kubectl apply -f sink-dingtalk.yaml
```

```yaml
apiVersion: v1
kind: Service
metadata:
  name: sink-dingtalk
  namespace: vanus
spec:
  selector:
    app: sink-dingtalk
  type: ClusterIP
  ports:
    - port: 8080
      name: sink-dingtalk
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: sink-dingtalk
  namespace: vanus
data:
  config.yml: |-
    bot:
      defalult: "bot1"
      webhooks:
        - chat_group: "bot1"
          url: "https://oapi.dingtalk.com/robot/send?access_token=XXXXXX"
          signature: "xxxxxxxxxx"
        - chat_group: "bot2"
          url: "https://oapi.dingtalk.com/robot/send?access_token=XXXXXX"
          signature: "xxxxxxxxxx"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sink-dingtalk
  namespace: vanus
  labels:
    app: sink-dingtalk
spec:
  selector:
    matchLabels:
      app: sink-dingtalk
  replicas: 1
  template:
    metadata:
      labels:
        app: sink-dingtalk
    spec:
      containers:
        - name: sink-dingtalk
          image: public.ecr.aws/vanus/connector/sink-dingtalk:latest
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
            name: sink-dingtalk
```

[vc]: https://docs.vanus.ai/introduction/concepts#vanus-connect
[robot]: https://open.dingtalk.com/document/orgapp/custom-robot-access
