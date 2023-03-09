# Source Kafka Proposal

## Description

The Kafka Source is used to pull data from 1 or multiple topics in a broker, tranform them into CloudEvents and send the events to the target.
This page describes the design of the Kafka Source in detail.

## Programming Language

-[ ] Golang
-[x] Java

## Prerequisites

- A Kafka server
- A Kafka topic ctreated

## Connector Details

### Configuration

The Kafka Source needs following configurations to work properly.

| Name              | Required | Default | Description                                                |
| :---------------- | :------- | :-----: | :--------------------------------------------------------- |
| target            | YES      |         | the target URL which Kafka Source will send CloudEvents to |
| bootstrap_servers | YES      |         | the kafka cluster bootstrap servers                        |
| group_id          | YES      |         | the kafka cluster consumer group id                        |
| topics            | YES      |         | the kafka topics listened by kafka source                  |



```text
> { "name": "Jason", "age": "30" }
```

### Connector Behavior

If an incoming data looks like:

```text
> { "name": "Jason", "age": "30" }
```
The CloudEvent will look like:

```JSON
{
  "id" : "4ad0b59fc-3e1f-484d-8925-bd78aab15123",
  "source" : "kafka_bootstrap_servers.mytopic",
  "type" : "kafka.message",
  "datacontenttype" : "application/json",
  "time" : "2022-09-07T10:21:49.668Z",
  "data" : {
	          "name": "Jason",
	          "age": "30"
  }
}
```

In order to get the best result the Data should be a string in a json format, but the kafka connector can receive any type of messages.

### Used Libraries/APIs

[Apache Kafka SDK](https://mvnrepository.com/artifact/org.apache.kafka/kafka)
