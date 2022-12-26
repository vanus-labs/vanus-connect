---
title: HTTP
---

# HTTP Sink 

## Overview

A [Vance Connector][vc] which receives CloudEvents and deliver specific data to the target URL.
## User guidelines

### Connector introduction

The HTTP Sink is a [Vance Connector][vc] which aims to handle incoming CloudEvents in a way that extracts the `data` part of the 
original event and deliver these extracted `data` to the target URL.

For example, if the incoming CloudEvent looks like:

```http
{
  "id" : "42d5b039-daef-4071-8584-e61df8fc1354",
  "source" : "vance-http-source",
  "specversion" : "V1",
  "type" : "http",
  "datacontenttype" : "application/json",
  "time" : "2022-05-17T18:44:02.681+08:00",
  "data" : {
    "myData" : "simulation event data <1>"
  }
}
```

The HTTP Sink will POST an HTTP request looks like:

``` json
> POST /payload HTTP/2

> Host: localhost:8080
> User-Agent: VanceCDK-HttpClient/1.0.0
> Content-Type: application/json
> Content-Length: 39

> {
>    "myData" : "simulation event data <1>"
> }
```

## HTTP Sink Configs

Users can specify their configs by either setting environments variables or mount a config.json to
`/vance/config/config.json` when they run the connector. Find examples of setting configs [here][config].

### Config Fields of the HTTP Sink

| Configs   | Description                                                            | Example                 |
|:----------|:-----------------------------------------------------------------------|:------------------------|
| v_target  | v_target is used to specify the target URL HTTP Sink will send data to | "http://localhost:8081" |
| v_port    | v_port is used to specify the port HTTP Sink is listening on           | "8080"                  |

## HTTP Sink Image

> docker.io/vancehub/sink-http

## Local Development

You can run the sink codes of the HTTP Sink locally as well.

### Building via Maven

```shell
$ cd sink-http
$ mvn clean package
```

### Running via Maven

```shell
$ mvn exec:java -Dexec.mainClass="com.linkall.sink.http.Entrance"
```

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md
[config]: https://github.com/linkall-labs/vance-docs/blob/main/docs/connector.md