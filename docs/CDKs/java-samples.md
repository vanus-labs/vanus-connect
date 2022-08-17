---
title: http
parent: cdk-java
grand_parent: connectors
nav_order: 2
---

# Connector Samples

Now, you're getting more familiar with the basic concepts of vance APIs.

Let's go through the detail of provided examples (there are two connector samples under `examples` directory).

## Recap the Concept of Source

Before checking out the sample codes, let's recap the concept of a `Source` connector.

> A Source is a connector that implements the following functions:
> - Retrieves data from an underlying data producer. Vance doesn't limit the way a source retrieves data. (e.g. A source MAY pull data from a message queue or act as a HTTP server waiting for data to be sent to it).
> - Transforms retrieved data into CloudEvents.
> - Uses standard HTTP POST requests to send CloudEvents to the target URL specified in V_TARGET.

## Source Example

### com.vance.source.Entrance

`Entrance.java` is the entrance of the source connector.

```java
01 public class Entrance {
02     public static void main(String[] args) {
03         VanceApplication.run(MyConnector.class);
04     }
05 }
```

The `main` method uses a one-line code to easily launch the connector programme.

`VanceApplication.run()` only needs one parameter, which is the implementation of either a [Sink or Source](api.md#connector-interface) interface.

### com.vance.source.MySource

`MySource` implemented all methods of `Source` interface.

```java
01 @Override
02 public void start(){
03    // TODO Initialize your Adapter
04    MyAdapter adapter = (MyAdapter) getAdapter();
05
06    // TODO receive your original data and transform it into a CloudEvent via your Adapter
07    // In this sample, we use a String as the original data
08    for (int i = 0; i < NUM_EVENTS; i++) {
09        String data = "Event number " + i;
10        // TODO: construct CloudEvents
11        CloudEvent event = adapter.adapt(data);
12        // Use EnvUtil to get the target URL the source will send to
13        // You can replace the default sink URL with yours in resources/config.json
14        String sink = EnvUtil.getVanceSink();
15        // TODO: deliver CloudEvents to endpoint ${V_TARGET}
16        sendCloudEvent(event,sink);
17    }
18 }
```

The `start()` method is one of the methods declared in `Source` interface. It met all [requirements](#recap-the-concept-of-source) the `Source` asks:
1. Using a for loop to generate original data, which is a String consisted of "Event number" and the index of the loop
2. Using `adapt()` method from `MyAdapter` to transform a String into a CloudEvent
3. Sending CloudEvents to the URL which specified in `resources/config.json`

```java
01 public Adapter getAdapter() {
02     return new MyAdapter();
03 }
```

`getAdapter()` method is another method declared in `Source` interface. Its purpose is to return an instance of `Adapter` interface.

```java
01 private void sendCloudEvent(CloudEvent event, String targetURL){
02    Future<HttpResponse<Buffer>> responseFuture;
03    // Send CloudEvent to vance_sink
04    responseFuture = VertxMessageFactory.createWriter(webClient.postAbs(targetURL))
05            .writeStructured(event, JsonFormat.CONTENT_TYPE); // Use structured mode.
06    responseFuture.onSuccess(resp->{
07        LOGGER.info("send CloudEvent success");
08    }).onFailure(t-> LOGGER.info("send task failed"));
09 }
```

`sendCloudEvent()` is a method used to send a CloudEvent to the target URL. In this example, I use the Vert.x as the HTTP framework, but you can choose whatever you want to POST the HTTP request.

> It's strongly recommended to use one of the HTTP implementations CloudEvent-sdk provided.

### com.vance.source.MyAdapter

`MyAdapter` is the implementation to convert the original data into a CloudEvent.

>⚠️ Note: Don't directly implement `Adapter` interface️. Instead, implement `Adapter1` or `Adapter2` based on the number of types you need to construct a CloudEvent.

```java
01 public class MyAdapter implements Adapter1<String> {
02    private static final CloudEventBuilder template = CloudEventBuilder.v1();
03    @Override
04    public CloudEvent adapt(String data) {
05        template.withId(UUID.randomUUID().toString());
06        URI uri = URI.create("vance-http-source");
07        template.withSource(uri);
08        template.withType("http");
09        template.withDataContentType("application/json");
10        template.withTime(OffsetDateTime.now());
11        template.withData(data.getBytes());
12
13        return template.build();
14    }
15 }
```

In this example, MyAdapter chose `Adapter1` to implement since the only thing it needs to construct a CloudEvent is a String.

Codes between line 5 and 11 are trying to fill required fields of a CloudEvent. Learn more about [CloudEvents Specification][ce] and [CloudEvents java-sdk][ce-sdk]
if you are not familiar with them.

[ce]: https://github.com/cloudevents/spec
[ce-sdk]: https://github.com/cloudevents/sdk-java