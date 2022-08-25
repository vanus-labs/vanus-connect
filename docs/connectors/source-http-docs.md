# HTTP Source 

## Overview

A [Vance Connector][vc] which transforms HTTP requests to CloudEvents and deliver them to the target URL.

## User Guidelines

### Connector Introduction

The HTTP Source is a [Vance Connector][vc] which aims to generate CloudEvents in a way that wraps all headers and body of the 
original request into the `data` field of a new CloudEvent.

For example, if an original request looks like:

```http
> POST /payload HTTP/2

> Host: localhost:8080
> User-Agent: VanceCDK-HttpClient/1.0.0
> Content-Type: application/json
> Content-Length: 39

> {
>    "myData" : "simulation event data <1>"
> }
```

This POST HTTP request will be transformed into a CloudEvent looks like:

``` json
{
  "id" : "42d5b039-daef-4071-8584-e61df8fc1354",
  "source" : "vance-http-source",
  "specversion" : "V1",
  "type" : "http",
  "datacontenttype" : "application/json",
  "time" : "2022-05-17T18:44:02.681+08:00",
  "data" : {
    "headers" : {
      "user-agent" : "VanceCDK-HttpClient/1.0.0",
      "content-type" : "application/json",
      "content-length" : "39",
      "host" : "localhost:8080"
    },
    "body" : {
      "myData" : "simulation event data <1>"
    }
  }
}
```

## HTTP Source Configs

Users can specify their configs by either setting environments variables or mount a config.json to
`/vance/config/config.json` when they run the connector. Find examples of setting configs [here][config].

### Config Fields of the HTTP Source

| Configs   | Description                                                                     | Example                 |
|:----------|:--------------------------------------------------------------------------------|:------------------------|
| v_target  | v_target is used to specify the target URL HTTP Source will send CloudEvents to | "http://localhost:8081" |
| v_port    | v_port is used to specify the port HTTP Source is listening on                  | "8080"                  |

## HTTP Source Image

> docker.io/vancehub/source-http

## Local Development

You can run the source codes of the HTTP Source locally as well.

### Building via Maven

```shell
$ cd connectors/source-http
$ mvn clean package
```

### Running via Maven

```shell
$ mvn exec:java -Dexec.mainClass="com.linkall.source.http.Entrance"
```

⚠️ NOTE: For better local development and test, the connector can also read configs from `main/resources/config.json`. So, you don't need to 
declare any environment variables or mount a config file to `/vance/config/config.json`.

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md
[config]: https://github.com/linkall-labs/vance-docs/blob/main/docs/connector.md