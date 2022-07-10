# Display Sink 

## Overview

A [Vance Connector][vc] which prints received CloudEvents. This is commonly used as a logger to check incoming data.

## User Guidelines

### Connector introduction

The Display Sink is a single function [Connector][vc] which aims to print incoming CloudEvents in JSON format.

For example, it will print the incoming CloudEvent looks like:

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

## Display Sink Configs

Users can specify their configs by either setting environments variables or mount a config.json to
`/vance/config/config.json` when they run the connector. Find examples of setting configs [here][config].

### Config Fields of the Display Sink

| Configs   | Required | Description                                                            | Example                 |
|:----------|:----|:-----------------------------------------------------------------------|:------------------------|
| v_port    |   false   | v_port is used to specify the port Display Sink is listening on           | "8080"                  |

## Display Sink Image

> docker.io/vancehub/display

## Local Development

You can run the sink codes of the Display Sink locally as well.

### Building via Maven

```shell
$ cd sink-http
$ mvn clean package
```

### Running via Maven

```shell
$ mvn exec:java -Dexec.mainClass="com.linkall.sink.display.Entrance"
```

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md
[config]: https://github.com/linkall-labs/vance-docs/blob/main/docs/connector.md