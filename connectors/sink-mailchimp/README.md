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

## Create the config file


```shell
cat << EOF > config.yml
target: <you server endpoint>
EOF
```

| Name          | Required | Default | Description                                                                      |
|:--------------|:--------:|:--------|----------------------------------------------------------------------------------|
| port          |    NO    | 8080    | the port which HTTP Sink listens on                                              |
| api_key       |   YES    |         | the mailchimp [api key][api_key]                                                 |
| list_id       |    NO    |         | the mailchimp [audience id][list_id]                                             |

The Mailchimp Sink tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the
position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

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

```

[vc]: https://docs.vanus.ai/introduction/concepts#vanus-connect
[api_key]: https://mailchimp.com/help/about-api-keys
[list_id]: https://mailchimp.com/help/find-audience-id
