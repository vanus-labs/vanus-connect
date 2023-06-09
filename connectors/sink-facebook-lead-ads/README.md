---
title: Facebook Lead Ads
---

# Facebook Lead Ads Sink

## Introduction

The Facebook Lead Ads Sink is a [Vanus Connector][vc] that aims to handle incoming CloudEvents in a way that extracts the `data`
part of the original event and <must: description...>

For example, the incoming CloudEvent looks like this:

```json
{
  "id": "88767821-92c2-477d-9a6f-bfdfbed19c6a",
  "source": "quickstart",
  "specversion": "1.0",
  "type": "test",
  "time": "2022-07-08T03:17:03.139Z",
  "datacontenttype": "application/json",
  "data": {
    "name": "Vanus Community",
    "follow_up_action_url": "<follow_up_action_URL>",
    "questions": {
        {
            "type": "FIRST_NAME"
        },
        {
            "type": "LAST_NAME"
        },
        {
            "type": "EMAIL"
        },
    }
    "privacy_policy": {
        "url": "<Policy_URL>"
    }
  }
}
```

The Facebook Lead Ads Sink will create a form name ``Vanus Community`` with fields for ``first_name``, ``last_name`` and ``email address``.

## Quickstart

<optional prerequisites but recommended>
### Prerequisites
- Have a container runtime (i.e., docker).
- ...
</optional>

### Create the config file

<optional: explanation>

```shell
cat << EOF > config.yml
<example config content>
...
EOF
```

| Name | Required  | Default | Description                           |
|:-----|:---------:|:--------|---------------------------------------|
| port |    NO     | 8080    | the port which <name> Sink listens on |

...

The <name> Sink tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the
position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm \
  -p 31080:8080 \
  -v ${PWD}:/vanus-connect/config \
  --name sink-<name> public.ecr.aws/vanus/connector/sink-<name>
```

### Test

<option: explanation>.

Open a terminal and use following command to send a CloudEvent to the Sink.

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    <incoming event example>
}'
```

<show result with a paragraph>

### Clean resource

```shell
docker stop sink-<name>
```

## Sink details

<optional>
### Extension Attributes

The <name> Sink have additional reactions if the incoming CloudEvent contains
following[Extension Attributes](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#extension-context-attributes)
.

| Attribute      | Required | Examples  | Description                           |
|:---------------|:--------:|:----------|:--------------------------------------|

...
</optional>

### Data format

The <name> Sink requires following data format in CloudEvent's `data` field.

```json
{
  <full example>
}
```

<optional>
### Examples

#### Example1

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "quick-start",
    "specversion": "1.0",
    "type": "sink-<name>",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "data": {
        <data1>
    }
}'
```

#### example2

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "quick-start",
    "specversion": "1.0",
    "type": "sink-<name>",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "<extension_attr1>": "test",
    "<extension_attr2>": "test",
    "data": {
        <data2>
    }
}'
```

## Run in Kubernetes

```shell
kubectl apply -f sink-<name>.yaml
```

```yaml
<must: content>
```

[vc]: https://www.vanus.dev/introduction/concepts#vanus-connect
