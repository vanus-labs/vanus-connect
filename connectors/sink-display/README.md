# Display Sink 

## Overview

A [Vance Connector][vc] which prints received CloudEvents. This is commonly used as a logger to check incoming data.

## Introduction

The Display Sink is a single function [Connector][vc] which aims to print incoming CloudEvents in JSON format.

For example, it will print the incoming CloudEvent looks like:

```http
{
  "specversion": "1.0",
  "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
  "source": "quickstart",
  "type": "quickstart",
  "datacontenttype": "application/json",
  "time": "2022-10-26T10:38:29.345Z",
  "data": {
    "myData": "simulation event data"
  }
}
```

## Quick Start

### Start Using Docker

mapping 8080 to 31080 in order to avoid port conflict.

```shell
docker run -d -p 31080:8080 --rm \
  -v ${PWD}:/vance/config \
  --name sink-display public.ecr.aws/vanus/connector/sink-display:latest
```

### Test
1. make a HTTP request
```shell
curl --location --request POST 'localhost:8080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "quickstart",
    "specversion": "1.0",
    "type": "quickstart",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "data": {
        "myData": "simulation event data"
    }
}'
```

2. view logs
```shell
docker logs sink-display
```

```shell
receive a new event, in total: 1
{
  "specversion": "1.0",
  "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
  "source": "quickstart",
  "type": "quickstart",
  "datacontenttype": "application/json",
  "time": "2022-10-26T10:38:29.345Z",
  "data": {
    "myData": "simulation event data"
  }
}
```

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md