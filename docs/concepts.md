---
title: Vance CDKs
nav_order: 2
---

# Vance Concepts

Vance aims to provide community-driven and reusable connectors for users to easily integrate with other services.
These out-of-box connectors generally act as proxies between external services and their applications.
Vance also provides fine-grained autoscaling (including to/from zero) for managed connectors in Kubernetes.

Vance defines the following terms:

- **Connector** - A connector is an image-based program that interacts with a specific underlying data source
  (e.g. Databases or other web services) on behalf of user applications.
  In Vance, a connector is either a Source or a Sink.
- **Source** - A Source is a connector that implements the following functions:
    - Retrieves data from an underlying data producer. Vance doesn't limit the way a source retrieves data.
      (e.g. A source MAY pull data from a message queue or act as a HTTP server waiting for data to be sent to it).
    - Transforms retrieved data into CloudEvents.
    - Uses standard HTTP POST requests to send CloudEvents to the target URI specified in `V_TARGET`.
- **Sink** - A connector that receives CloudEvents and uses the data in specific logics.
  (e.g. A MySQL Sink extracts useful data from CloudEvents and writes them to a MySQL database).
- **CloudEvents** - A specification for describing event data in common formats to provide interoperability
  across services, platforms and systems.
- **Autoscaling** - The ability to automatically scale managed connectors to match demand.
  (e.g. deploying more pods reponds to increased load or scales resources down if the load decreases).
- **Autoscaling Criteria** - Vance connectors can dynamically scale based on the following rules:
    - HTTP traffic
    - CPU or memory load
    - KEDA-supported scalers
- **KEDA** - A Kubernetes-based Event Driven Autoscaling component. It provides event driven scale
  for any container running in Kubernetes.
