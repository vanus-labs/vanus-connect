---
title: Email
---

# Email Sink

## Introduction

The Email Sink is a [Vanus Connector][vc] which aims to handle incoming CloudEvents in a way that extracts the `data` part of the
original event and deliver these extracted `data` to SMTP server.

For example, the incoming CloudEvent looks like:

```json
{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "quick-start",
    "specversion": "1.0",
    "type": "quick-start",
    "datacontenttype": "text/plain",
    "time": "2022-10-26T10:38:29.345Z",
    "xvemailfrom": "example@example.com",
    "xvemailrecipients": "demo@demo.com",
    "xvemailsubject": "test",
    "data": "Hello Email Sink"
}
```

then recipients will receive an email like:
![received.png](https://raw.githubusercontent.com/linkall-labs/vanus-connect/main/connectors/sink-email/received.png)

## Quick Start

in this section, we show how to use Email Sink sends a text message to recipients.

### Create Config file

replace `<from_email_address>`, `<from_email_address>`, and `<smtp server address>` to yours.

```shell
cat << EOF > config.yml
default: "<from_email_address>"
email:
  - account: "<from_email_address>"
    password: "<password>"
    host: "<smtp server address>"
EOF
```

| Name              | Required | Default | Description                                                                                            |
|:------------------|:--------:|:-------:|--------------------------------------------------------------------------------------------------------|
| port              |    NO    |  8080   | the port which Email Sink listens on                                                                   |
| default           | **YES**  |    -    | Email Sink supports multiple email accounts as sender, you could set the default account by this field |
| email.[].account  | **YES**  |    -    | email account address you want to use                                                                  |
| email.[].password | **YES**  |    -    | password for account authentication                                                                    |
| email.[].host     | **YES**  |    -    | SMTP server address                                                                                    |
| email.[].port     |  **NO**  |   25    | SMTP server port                                                                                       |
| email.[].format   |  **NO**  |  text   | `text` or `html`                                                                                       |

The Email Sink tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm \
  -p 31080:8080 \
  -v ${PWD}:/vanus-connect/config \
  --name sink-email public.ecr.aws/vanus/connector/sink-email
```

### Test

replace `<from_email_address>`, `<from_email_address>`, and `<smtp server address>` to yours.

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "quick-start",
    "specversion": "1.0",
    "type": "quick-start",
    "datacontenttype": "text/plain",
    "time": "2022-10-26T10:38:29.345Z",
    "xvemailrecipients": "<to address>",
    "xvemailsubject": "Quick Start",
    "data": "Hello Email Sink"
}'
```

now, you cloud see a new email in your mailbox.
![received.png](https://raw.githubusercontent.com/linkall-labs/vanus-connect/main/connectors/sink-email/received.png)

### Clean

```shell
docker stop sink-email
```

## Sink details

### Extension Attributes

Email Sink has defined a few [CloudEvents Extension Attribute](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#extension-context-attributes)
to determine how to process event.

| Attribute         | Required | Examples                    | Description                                                               |
|:------------------|:--------:|-----------------------------|---------------------------------------------------------------------------|
| xvemailsubject    | **YES**  | Test Email Sink             | The subject of the email                                                  |
| xvemailrecipients | **YES**  | a@example.com,b@example.com | The recipients addresses where you want to send, use `,` to separate.     |
| xvemailfrom       |    NO    | example@example.com         | Which email account(from address) that configured in Sink you want to use |
| xvemailformat     |    NO    | text                        | what format of your email content, `text` or `html`                       |

### Examples

#### Sending email to single recipient with default account
```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "quick-start",
    "specversion": "1.0",
    "type": "quick-start",
    "datacontenttype": "text/plain",
    "time": "2022-10-26T10:38:29.345Z",
    "xvemailrecipients": "a@example.com",
    "xvemailsubject": "Quick Start",
    "data": "Hello Email Sink"
}'
```

#### Sending email to multiple recipients with default account
```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "quick-start",
    "specversion": "1.0",
    "type": "quick-start",
    "datacontenttype": "text/plain",
    "time": "2022-10-26T10:38:29.345Z",
    "xvemailrecipients": "a@example.com,b@example.com",
    "xvemailsubject": "Quick Start",
    "data": "Hello Email Sink"
}'
```

#### Sending email to multiple recipients with specified account
```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "quick-start",
    "specversion": "1.0",
    "type": "quick-start",
    "datacontenttype": "text/plain",
    "time": "2022-10-26T10:38:29.345Z",
    "xvemailrecipients": "a@example.com,b@example.com",
    "xvemailsubject": "Quick Start",
    "xvemailfrom": "example@example.com",
    "data": "Hello Email Sink"
}'
```

## Run in Kubernetes

```shell
kubectl apply -f sink-email.yaml
```

```yaml
apiVersion: v1
kind: Service
metadata:
  name: sink-email
  namespace: vanus
spec:
  selector:
    app: sink-email
  type: ClusterIP
  ports:
    - port: 8080
      name: sink-email
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: sink-email
  namespace: vanus
data:
  config.yml: |-
    default: "example@example.com"
    email:
      - account: "example@example.com"
        password: "xxxxxxxx"
        host: "smtp.example.com"
        format: "html"
      - account: "demo@demo.com"
        password: "xxxxxxxx"
        host: "smtp.demo.com"
        format: "text"

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sink-email
  namespace: vanus
  labels:
    app: sink-email
spec:
  selector:
    matchLabels:
      app: sink-email
  replicas: 1
  template:
    metadata:
      labels:
        app: sink-email
    spec:
      containers:
        - name: sink-email
          image: public.ecr.aws/vanus/connector/sink-email:latest
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
            name: sink-email
```

## Integrate with Vanus

This section shows how a sink connector can receive CloudEvents from a running [Vanus cluster](https://github.com/linkall-labs/vanus).

1. Run the sink-email.yaml
```shell
kubectl apply -f sink-email.yaml
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
  --sink 'http://sink-email:8080'
```

[vc]: https://www.vanus.ai/introduction/concepts#vanus-connect