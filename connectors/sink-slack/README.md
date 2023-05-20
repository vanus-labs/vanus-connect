---
title: Slack
---

# Slack Sink

## Introduction

The Slack Sink is a [Vanus Connector][vc] that aims to handle incoming CloudEvents by extracting the `data` part of the
original event and delivering it to a Slack channel.

For example, if an incoming CloudEvent looks like:

```json
{
  "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
  "source": "quick-start",
  "specversion": "1.0",
  "type": "quick-start",
  "datacontenttype": "application/json",
  "time": "2022-10-26T10:38:29.345Z",
  "data": {
    "blocks": [
      {
        "type": "section",
        "text": {
          "text": "A message italicized text_.",
          "type": "mrkdwn"
        }
      },
      {
        "type": "section",
        "text": {
          "type": "plain_text",
          "text": "This is a plain text section block.",
          "emoji": true
        }
      }
    ]
  }
}
```

The Slack channel will receive a message like:
![message](https://github.com/vanus-labs/vanus-connect/blob/main/connectors/sink-slack/message.png?raw=true)

## Quick Start

In this section we will show you how to use Slack Sink to send a message to a Slack Channel.

### Prerequisites

- Have a container runtime (i.e., docker).
- Have a [Slack account](https://api.slack.com/apps).

### Create an App in Slack

1. Create an app on slack.
   ![message.png](https://github.com/vanus-labs/vanus-connect/blob/main/connectors/sink-slack/createApp.png?raw=true)
2. Select `From scratch`.
   ![message.png](https://github.com/vanus-labs/vanus-connect/blob/main/connectors/sink-slack/selectFromScratch.png?raw=true)
3. Set the bot name and Workspace.
4. Click on permissions in the central menu.
   ![message.png](https://github.com/vanus-labs/vanus-connect/blob/main/connectors/sink-slack/clickPerm.png?raw=true)
5. Scopes 'Add OAuth Scope' `chat:write` and `chat:write.public`.
   ![message.png](https://github.com/vanus-labs/vanus-connect/blob/main/connectors/sink-slack/setPerm.png?raw=true)
6. Install to workspace.
   ![message.png](https://github.com/vanus-labs/vanus-connect/blob/main/connectors/sink-slack/installWorkspace.png?raw=true)
7. Set your configurations with the `Bot User OAuth Token` in OAuth & Permissions.
   ![message.png](https://github.com/vanus-labs/vanus-connect/blob/main/connectors/sink-slack/oath.png?raw=true)

### Create the config file

```shell
cat << EOF > config.yml
token: "xoxp-422301774731343243235Example"
default_channel: "#general"
```

| Name            | Required | Default | Description                                                             |
|:----------------|:--------:|:--------|:------------------------------------------------------------------------|
| port            |    NO    | 8080    | the port which Slack Sink listens on                                    |
| token           |   YES    |         | OAuth Token of this app, more visit: https://api.slack.com/legacy/oauth |
| default_channel |   YES    |         | set default channel the messages send to if attribute was not be set    |

The Slack Sink tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the
position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm \
  -p 31080:8080 \
  -v ${PWD}:/vanus-connect/config \
  --name sink-slack public.ecr.aws/vanus/connector/sink-slack
```

### Test

We have designed for you a sandbox environment, removing the need to use your local machine.
You can run Connectors directly and safely on the Playground.

Open a terminal and use the following command to send a CloudEvent to the Sink.

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
       "blocks": [
         {
           "type": "section",
           "text": {
             "text": "A message italicized text_.",
             "type": "mrkdwn"
           }
         },
         {
           "type": "section",
           "text": {
             "type": "plain_text",
             "text": "This is a plain text section block.",
             "emoji": true
           }
         }
       ]
  }
}'
```

Now, you should see in your slack channel your message.
![message.png](https://github.com/vanus-labs/vanus-connect/blob/main/connectors/sink-slack/message.png?raw=true)

### Clean

```shell
docker stop sink-slack
```

## Sink details

### Extension Attributes

The Slack Sink has additional options if the incoming CloudEvent contains the
following[Extension Attributes](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#extension-context-attributes)
.

| Attribute | Required  | Examples        | Description |
|:----------|:---------:|:----------------|:------------|
| xvchannel |    NO     | #general        | chanel      |

### Data format

the event data must be `JSON` format, and only with key `blocks` , refer [doc](https://api.slack.com/reference/block-kit/blocks) 

### Examples

#### Sending a message to the default channel.

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
       "blocks": [
         {
           "type": "section",
           "text": {
             "text": "A message italicized text_.",
             "type": "mrkdwn"
           }
         },
         {
           "type": "section",
           "text": {
             "type": "plain_text",
             "text": "This is a plain text section block.",
             "emoji": true
           }
         }
       ]
    }
}'
```

#### Sending a message to a specific channel.

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
    "xvchannel": "#team-a",
    "data": {
       "blocks": [
         {
           "type": "section",
           "text": {
             "text": "A message italicized text_.",
             "type": "mrkdwn"
           }
         },
         {
           "type": "section",
           "text": {
             "type": "plain_text",
             "text": "This is a plain text section block.",
             "emoji": true
           }
         }
       ]
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
    token: "xoxp-422301774731343243235Example"
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

[vc]: https://docs.vanus.ai/introduction/concepts#vanus-connect
