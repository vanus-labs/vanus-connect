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
  "datacontenttype": "application/json",
  "time": "2022-10-26T10:38:29.345Z",
  "data": {
    "subject": "vanus auto mail",
    "body": "this is vanus test email",
    "receiver": "example@example.com"
  }
}
```

then recipients will receive an email like:
![received.png](https://raw.githubusercontent.com/vanus-labs/vanus-connect/main/connectors/sink-email/received.png)

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

| Name              | Required | Default | Description                                                             |
|:------------------|:--------:|:-------:|-------------------------------------------------------------------------|
| port              |    NO    |  8080   | the port which Email Sink listens on                                    |
| default           |    NO    |    -    | the default account email, if not set it will be set to the first email |
| email.[].account  | **YES**  |    -    | email account address you want to use                                   |
| email.[].password | **YES**  |    -    | password for account authentication                                     |
| email.[].host     | **YES**  |    -    | SMTP server address                                                     |
| email.[].port     |  **NO**  |   25    | SMTP server port                                                        |
| email.[].format   |  **NO**  |  text   | `text` or `html`                                                        |

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
curl --location --request POST 'localhost:18080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "quick-start",
    "specversion": "1.0",
    "type": "quick-start",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "data": {
      "subject": "vanus auto mail",
      "body": "this is vanus test email",
      "receiver": "example@example.com"
    }
}'
```

now, you cloud see a new email in your mailbox.
![received.png](https://raw.githubusercontent.com/vanus-labs/vanus-connect/main/connectors/sink-email/received.png)

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
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "data": {
      "subject": "vanus auto mail",
      "body": "this is vanus test email",
      "receiver": "example@example.com"
    }
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
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "data": {
      "subject": "vanus auto mail",
      "body": "this is vanus test email",
      "receiver": "example1@example.com,example2@example.com"
    }
}'
```

[vc]: https://docs.vanus.ai/introduction/concepts#vanus-connect