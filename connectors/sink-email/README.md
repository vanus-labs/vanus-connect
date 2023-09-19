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
    "recipients": "example@example.com",
    "body_type": "text"
  }
}
```

then recipients will receive an email
## Quick Start

in this section, we show how to use Email Sink sends a text message to recipients.

### Create Config file

```shell
cat << EOF > config.yml
email:
    account: "example@example.com"
    password: "123456"
    host: "smtp.gmail.com"
EOF
```

| Name             | Required | Default | Description                                                             |
|:-----------------|:--------:|:-------:|-------------------------------------------------------------------------|
| port             |    NO    |  8080   | the port which Email Sink listens on                                    |
| email.account    | **YES**  |    -    | email account address you want to use                                   |
| email.password   | **YES**  |    -    | password for account authentication                                     |
| email.host       | **YES**  |    -    | SMTP server address                                                     |
| email.port       |  **NO**  |   25    | SMTP server port                                                        |

The Email Sink tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm \
  -p 31080:8080 \
  -v ${PWD}:/vanus-connect/config \
  --name sink-email public.ecr.aws/vanus/connector/sink-email
```

### Test


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
      "recipients": "example@example.com"
    }
}'
```

now, you cloud see a new email in your mailbox.

### Clean

```shell
docker stop sink-email
```

## Sink details

### Event Data Schema


| Attribute  | Required | Description                                       |
|:-----------|:--------:|---------------------------------------------------|
| subject    |   YES    | Email subject                                     |
| body       |   YES    | Email body                                        |
| recipients |   YES    | Email recipients                                  |
| body_type  |    NO    | Email body type, `text` or `html`, default `text` |

[vc]: https://docs.vanus.ai/introduction/concepts#vanus-connect