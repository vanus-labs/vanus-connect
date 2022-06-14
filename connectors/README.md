# vance-connectors

Vance aims to provide community-driven and reusable connectors for users to easily integrate with other services.
These out-of-box connectors generally act as proxies between external services and their applications.
Vance also provides fine-grained autoscaling (including to/from zero) for managed connectors in Kubernetes.

## Concepts

- **Connector** - A connector is an image-based program that interacts with a specific underlying data source
  (e.g. Databases or other web services) on behalf of user applications.
  In Vance, a connector is either a Source or a Sink.
- **Source** - A Source is a connector that implements the following functions:
    - Retrieves data from an underlying data producer. Vance doesn't limit the way a source retrieves data.
      (e.g. A source MAY pull data from a message queue or act as a HTTP server waiting for data to be sent to it).
    - Transforms retrieved data into CloudEvents.
    - Uses standard HTTP POST requests to send CloudEvents to the target URI specified in `V_TARGET`.
- **Sink** - A connector that receives CloudEvents and uses the data in specific logics. 

## Current Connectors

| Connectors<div style="width:200px"> | Description                                           |
| :---------------------------- | :---------------------------------------------- |
| sink-elasticsearch   |  |
| sink-http              | The HTTP Sink is a Vance Connector which aims to handle incoming CloudEvents in a way that extracts the data part of the original event and deliver these extracted data to the target URL.  |
| source-alicloud-billing         |                                          |
| source-aws-billing              |         |
| source-http    |    The HTTP Source is a Vance Connector which aims to generate CloudEvents in a way that wraps all headers and body of the original request into the data field of a new CloudEvent.   |
