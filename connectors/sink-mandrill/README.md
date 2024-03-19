---
title: Mandrill
---

# Mandrill Sink

## Introduction

The Mandrill Sink is a [Vanus Connector][vc] which aims to handle incoming CloudEvents in a way that extracts the `data`
part of the original event and deliver to the mandrill.

For example, if the incoming CloudEvent looks like:

```http
{
  "id" : "42d5b039-daef-4071-8584-e61df8fc1354",
  "source" : "quickstart",
  "specversion" : "1.0",
  "type" : "quickstart",
  "datacontenttype" : "application/json",
  "time" : "2023-01-26T10:38:29.345Z",
  "xvaction" : "messages.send-template",
  "xvtemplatename" : "test-template",
  "data" : {
    "to": [
        {"email":"test@vanus.ai","name":"test"}
    ]
  }
}
```

The Mandrill Sink will send an email to mandrill:

## Config

| Name          | Required | Default | Description                             |
|:--------------|:--------:|:--------|-----------------------------------------|
| port          |    NO    | 8080    | the port which Mandrill Sink listens on |
| api_key       |   YES    |         | the mandrill api key                    |
| template_name |    NO    |         | the default template name               |

The Mandrill Sink tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the
position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

## Sink details

### Extension Attributes

The Mandrill Sink has additional options if the incoming CloudEvent contains the following[Extension Attributes][ce_extension].

| Attribute      | Required | Examples      | Description   |
|:---------------|:--------:|:--------------|:--------------|
| xvaction       |   YES    | messages.send | action        |
| xvtemplatename |    NO    | test          | template name |

### Data Format

The Mandrill Sink call mandrill api to do, the event data must  meet api param, such as [messages.send](https://mailchimp.com/developer/transactional/api/messages/send-new-message/)

## Start with Docker

```shell
docker run -it --rm \
  -p 31080:8080 \
  -v ${PWD}:/vanus-connect/config \
  --name sink-mandrill public.ecr.aws/vanus/connector/sink-mandrill
```

### Test

Open a terminal and use the following command to send a CloudEvent to the Sink.

```shell
curl --location --request POST 'localhost:18080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
  "id" : "42d5b039-daef-4071-8584-e61df8fc1354",
  "source" : "quickstart",
  "specversion" : "1.0",
  "type" : "quickstart",
  "datacontenttype" : "application/json",
  "time" : "2023-01-26T10:38:29.345Z",
  "xvaction" : "messages.send-template",
  "xvtemplatename" : "test-template",
  "data" : {
    "to": [
        {"email":"test@vanus.ai","name":"test"}
    ]
  }
}'
```

Then your mandrill contacts will add one.

### Clean resource

```shell
docker stop sink-mandrill
```


[vc]: https://docs.vanus.ai/introduction/concepts#vanus-connect
[ce_extension]: https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#extension-context-attributes
