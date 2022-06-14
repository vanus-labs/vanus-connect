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

## Vance Connector Configs

Users can specify their configs by either setting environments variables or mount a config.json to
`/vance/config/config.json` when they run the connector.

### Set environments variables for HTTP Sink

```
//use V_TARGET to specify the target URI HTTP Sink will send CloudEvents to
--env "V_TARGET"="http://localhost:8081"

//use V_PORT to specify the port HTTP Sink is listening on
--env "V_PORT"="8080"
```

⚠️ **NOTE: ENV keys MUST be uppercase** ⚠️

### Set config.json and mount it on `/vance/config/config.json`

```json
{
  //use V_TARGET to specify the target URI HTTP Sink will send CloudEvents to.
  //use V_PORT to specify the port HTTP Sink is listening on.
  //JSON standard does not allow comments. Remember to delete these comments when you copy configs.
  "v_target": "http://localhost:8081",
  "v_port": "8080"
}
```

⚠️ **NOTE: json keys MUST be lowercase** ⚠️

## HTTP Source Image

> docker.io/vancehub/sink-http

## Local Development

You can run the source codes of the HTTP Source locally as well.

### Building via Maven

```shell
$ cd sink-http
$ mvn clean package
```

### Running via Maven

```shell
$ mvn exec:java -Dexec.mainClass="com.linkall.sink.http.Entrance"
```

[vc]: https://github.com/JieDing/vance-docs/blob/main/docs/concept.md