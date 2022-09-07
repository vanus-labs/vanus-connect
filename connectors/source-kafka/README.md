# Kafka Source 

## Overview

A [Vance Connector][vc] which transforms Kafka messages from topics to CloudEvents and deliver them to the target URL.

## User Guidelines

### Connector Introduction

The Kafka Source is a [Vance Connector][vc] which aims to generate CloudEvents in a way that wraps the body of the 
original message into the `data` field of a new CloudEvent.
## The ideal message
The ideal type of event for the Kafka source is a String in a JSON format. But it can handle any other type of data provided by Kafka. 
> JSON Formatted String
> String = "{ "name": "Jason", "age": "30"}"
>

For example, if an original message looks like:
... json
> { "name": "Jason", "age": "30" }
```

A Kafka message transformed into a CloudEvent looks like:

``` JSON
{
  "id" : "4ad0b59fc-3e1f-484d-8925-bd78aab15123",
  "source" : "kafka.localhost.topic2",
  "type" : "kafka.message",
  "datacontenttype" : "application/json or Plain/text",
  "time" : "2022-09-07T10:21:49.668Z",
  "data" : {
	 "name": "Jason",
	 "age": "30"
	 }
}
```

## Kafka Source Configs

Users can specify their configs by either setting environments variables or mounting a config.json to
`/vance/config/config.json` when they run the connector. Find examples of setting configs [here][config].

### Config Fields of the kafka Source

| Configs   | Description                                                                     | Example                 |
|:----------|:--------------------------------------------------------------------------------|:------------------------|
| v_target  | v_target is used to specify the target URL HTTP Source will send CloudEvents to | "http://localhost:8081" |
| KAFKA_SERVER_URL    | The URL of the Kafka Cluster the Kafka Source is listening on                  | "8080"                  |
| KAFKA_SERVER_PORT    | v_port is used to specify the port Kafka Source is listening on                  | "8080"                  |
| CLIENT_ID    |  An optional identifier for multiple Kafka Sources that is passed to a Kafka broker with every request.                  | "kafkaSource"                  |
| TOPIC_LIST    | The source will listen to the topic or topics specified.                   | "topic1"  or "topic1, topic2, topic3"                 |

## Kafka Source Image

> vancehub/source-kafka

## Local Development

You can run the source codes of the Kafka Source locally as well.

### Building via Maven

```shell
$ cd connectors/source-Kafka
$ mvn clean package
```

### Running via Maven

```shell
$ mvn exec:java -Dexec.mainClass="com.linkall.source.Kafka.Entrance"
```

⚠️ NOTE: For better local development and test, the connector can also read configs from `main/resources/config.json`. So, you don't need to 
declare any environment variables or mount a config file to `/vance/config/config.json`.

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md
[config]: https://github.com/linkall-labs/vance-docs/blob/main/docs/connector.md
