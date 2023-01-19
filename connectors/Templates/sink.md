---
title: <name>
---

# <name> Sink

## Introduction

The <name> Sink is a [Vanus Connector](https://www.vanus.dev/introduction/concepts#vanus-connect) that aims to handle incoming CloudEvents in a way that extracts the `data` part of the
original event and <must: description...>

For example, the incoming CloudEvent looks like this:

```json
<incoming event example>
```

The <name> Sink will ... (eg. send a message to the Slack channel)


## Quickstart

<optional prerequisites but recommended>
### Prerequisites
- Have a container runtime (i.e., docker).
- ...
</optional>

### Create the config file

<optional: explanation>

Replace `<config1>`, `<config2>`, and `<config3>` to yours.

```shell
cat << EOF > config.yml
<example config content>
EOF
```

### start with Docker

```shell
docker run --rm \
  -p 31080:8080 \
  -v ${PWD}:/vanus-connect/config \
  --name sink-<name> public.ecr.aws/vanus/connector/sink-<name>:latest
```

### Test
<option: explanation>.

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    <incoming event example>
}'
```

<show result with a paragraph>

### clean resource

```shell
docker stop sink-<name>
```

## How to use

### Configuration

The default path is `/vance/config/config.yml`. if you want to change the default path, you can set env `CONNECTOR_CONFIG` to
tell Sink <Name>.


| Name                                 | Required | Default | Description                                                                                                                                       |
|:-------------------------------------|:--------:|:-------:|---------------------------------------------------------------------------------------------------------------------------------------------------|
| port                                 |    NO    |  8080   | the pot Sink MongoDB receives incoming events                                                                                                     |
...


```yaml
<full example>
```

### Extension Attributes

Sink <name> has defined a few [CloudEvents Extension Attribute](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#extension-context-attributes)
to determine how to process events.


| Attribute      | Required | Examples | Description                          |
|:---------------|:--------:|----------|--------------------------------------|
...

### Data

| Item                  | Required |   Type   | Default | Description                                                   |
|:----------------------|:--------:|:--------:|:-------:|---------------------------------------------------------------|
...


```json
{
  <full example>
}
```

### Examples

#### example1

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

#### example3

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "quick-start",
    "specversion": "1.0",
    "type": "sink-mongodb",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "<extension_attr3>": "test",
    "data": {
         <data2>
    }
}'
```

### Run in kubernetes
```yaml
<must: content>
```