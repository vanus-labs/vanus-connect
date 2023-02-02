---
title: Slack
---

# Slack Sink

## Introduction

The Slack Sink is a [Vanus Connector][vc] that aims to handle incoming CloudEvents by extracting the `data` part of the original event and delivering it to a Slack channel.

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
    "subject": "Test",
    "message": "Hello Slack!:wave: This is Sink Slack!"
  }
}
```

The Slack channel will receive a message like:
![message](https://github.com/linkall-labs/vanus-connect/blob/main/connectors/sink-slack/message.png?raw=true)

## Quick Start

In this section we will show you how to use Slack Sink to send a message to a Slack Channel.

### Prerequisites
- Have a container runtime (i.e., docker).
- Have a [Slack account](https://api.slack.com/apps).

### Create an App in Slack

  1. Create an app on slack.
![message.png](https://github.com/linkall-labs/vanus-connect/blob/main/connectors/sink-slack/createApp.png?raw=true)
  2. Select `From scratch`.
![message.png](https://github.com/linkall-labs/vanus-connect/blob/main/connectors/sink-slack/selectFromScratch.png?raw=true)
  3. Set the bot name and Workspace.
  4. Click on permissions in the central menu.
![message.png](https://github.com/linkall-labs/vanus-connect/blob/main/connectors/sink-slack/clickPerm.png?raw=true)
  5. Scopes 'Add OAuth Scope' `chat:write` and `chat:write.public`.
![message.png](https://github.com/linkall-labs/vanus-connect/blob/main/connectors/sink-slack/setPerm.png?raw=true)
  6. Install to workspace.
![message.png](https://github.com/linkall-labs/vanus-connect/blob/main/connectors/sink-slack/installWorkspace.png?raw=true)
  7. Set your configurations with the `Bot User OAuth Token` in OAuth & Permissions.
![message.png](https://github.com/linkall-labs/vanus-connect/blob/main/connectors/sink-slack/oath.png?raw=true)
  
### Create the config file

```shell
cat << EOF > config.yml
default: "test_app"
slack:
  - app_name: "test_app"
    token: "xoxp-422301774731343243235Example"
    default_channel: "#general"
EOF
```

| Name                     | Required  | Default | Description                                                                                         |
|:-------------------------|:---------:|:--------|:----------------------------------------------------------------------------------------------------|
| port                     |    NO     | 8080    | the port which Slack Sink listens on                                                                |
| default                  |    YES    |         | the default app name if event attribute doesn't have `xvslackapp`                                   |
| slack.[].app_name        |    YES    |         | custom slack app name as identifier                                                                 |
| slack.[].token           |    YES    |         | OAuth Token of this app, more visit: https://api.slack.com/legacy/oauth                             |
| slack.[].default_channel |    NO     |         | set default channel the messages send to if attribute was not be set, use `,` to separate multiples |

The Slack Sink tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

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
        "subject": "Test",
        "message": "Hello Slack!:wave: This is Sink Slack!"
    }
}'
```

Now, you should see in your slack channel your message.
![message.png](https://github.com/linkall-labs/vanus-connect/blob/main/connectors/sink-slack/message.png?raw=true)

### Clean

```shell
docker stop sink-slack
```

## Sink details

### Extension Attributes

The Slack Sink has additional options if the incoming CloudEvent contains the following[Extension Attributes](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#extension-context-attributes).

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

#### Sending a message from the default app to the default channel.

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

#### Sending a message from a specific app to a specific channel.

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

#### Sending a message from a specific app to multiple channels.

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
