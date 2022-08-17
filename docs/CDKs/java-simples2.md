---
title: http
parent: cdk-java
grand_parent: connectors
nav_order: 3
---

# Connector Samples

After getting more familiar the source example, let's move to the detail of the sink example.

## Recap the Concept of Sink

Before checking out the sample codes, let's recap the concept of a `Sink` connector.

> A connector that receives CloudEvents and uses the data in specific logics. (e.g. A MySQL Sink extracts useful data from CloudEvents and writes them to a MySQL database).
## Sink Example

### com.vance.sink.Entrance

Like the source example, `Entrance.java` is again the entrance of the connector.

```java
01 public class Entrance {
02     public static void main(String[] args) {
03         VanceApplication.run(MySink.class);
04     }
05 }
```

The `main` method uses a one-line code to easily launch the connector programme.

`VanceApplication.run()` only needs one parameter, which is the implementation of either a [Sink or Source](api.md#connector-interface) interface.

### com.vance.sink.MySink

`MySink` implemented the `start()` method of `Sink` interface.

```java
01 @Override
02 public void start(){
03    // TODO write a HTTP Server which can handle requests based on CloudEvents format
04    HttpServer server = HttpServer.createHttpServer();
05    // Use ceHandler method to tell HttpServer logics you want to do with an incoming CloudEvent
06    server.ceHandler(event -> {
07        int num = eventNum.addAndGet(1);
08        // print number of received events
09        LOGGER.info("receive a new event, in total: "+num);
10        // Use JsonMapper to wrap a CloudEvent into a JsonObject for better printing
11        JsonObject js = JsonMapper.wrapCloudEvent(event);
12        LOGGER.info(js.encodePrettily());
13    });
14    server.listen();
15 }
```

The `start()` method above mainly:
1. create an HTTP server waiting for requests with CloudEvents
2. print the incoming CloudEvents in json format

- Line 4 created an HTTP server by using the `HttpServer` provided by the cdk
- The `ceHandler() method` registers a `Handler` to deal with CloudEvents for the server.
- Codes between line 6 and 13 demonstrated how to deal with an incoming CloudEvent request, which is simply to log the event in json format.
- Line 11 wrapped the incoming CloudEvent into a JsonObject for better printing

Utilities like `HttpServer` and `JsonMapper` are designed to reduce redundant works for developers.

You can write your own logics to handle incoming CloudEvents, for example, extract some fields from CloudEvents and insert them into a MySQL database. Then, the connector becomes a MySQL sink.

[ce]: https://github.com/cloudevents/spec
[ce-sdk]: https://github.com/cloudevents/sdk-java