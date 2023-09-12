---
title: HTTP
---

# Mailchimp Sink

## Introduction

The Mailchimp Sink is a [Vanus Connector][vc] which aims to handle incoming CloudEvents in a way that extracts the `data`
part of the original event and deliver to the mailchimp.

For example, if the incoming CloudEvent looks like:

```http
{
  "id" : "42d5b039-daef-4071-8584-e61df8fc1354",
  "source" : "quickstart",
  "specversion" : "1.0",
  "type" : "quickstart",
  "datacontenttype" : "application/json",
  "time" : "2023-01-26T10:38:29.345Z",
  "xvaction" : "member.add",
  "data" : {
    "email_address": "text@linkall.com",
    "status": "subscribed"
  }
}
```

The Mailchimp Sink will add a member to mailchimp:

## Config

| Name        | Required | Default | Description                               |
|:------------|:--------:|:--------|-------------------------------------------|
| port        |    NO    | 8080    | the port which Mailchimp Sink listens on  |
| api_key     |   YES    |         | the mailchimp [api key][api_key]          |
| audience_id |   YES    |         | the mailchimp [audience id][audience_id]  |

The Mailchimp Sink tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the
position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

## Sink details

### Extension Attributes

The Mailchimp Sink has additional options if the incoming CloudEvent contains the following[Extension Attributes][ce_extension].

| Attribute | Required  | Examples   | Description |
|:----------|:---------:|:-----------|:------------|
| xvaction  |    NO     | member.put | action      |

### Data Format

The Mailchimp Sink call mailchimp api to do, the event data must  meet api param, such as [member.add](https://mailchimp.com/developer/marketing/api/list-members/add-member-to-list/)

## Start with Docker

```shell
docker run -it --rm \
  -p 31080:8080 \
  -v ${PWD}:/vanus-connect/config \
  --name sink-mailchimp public.ecr.aws/vanus/connector/sink-mailchimp
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
  "xvaction" : "member.add",
  "data" : {
    "email_address": "text@linkall.com",
    "status": "subscribed",
    "tags":["test"]
  }
}'
```

Then your mailchimp contacts will add one.

### Clean resource

```shell
docker stop sink-mailchimp
```


[vc]: https://docs.vanus.ai/introduction/concepts#vanus-connect
[api_key]: https://mailchimp.com/help/about-api-keys
[audience_id]: https://mailchimp.com/help/find-audience-id
[ce_extension]: https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#extension-context-attributes
