---
layout: default
title: Home
nav_order: 1
description: "The Java Connector-Development Kit (CDK) is a collection of Java packages to help you to build a new [Vance Connector][vc] in minutes."
permalink: /
---

# Java CDK for Vance
{: .fs-9 }

## Introduction

The cdk aims to speed up the development of a vance connector by offering some utilities including:
- HTTP implementations (either to handle general HTTP requests or CloudEvents)
- Config implementation to load user-specific configs
- The ability to interact with the Vance operator

## Getting started

### Using the cdk-java

To use the cdk-java, add following dependency to your pom.xml

```
<dependency>
    <groupId>com.linkall</groupId>
    <artifactId>cdk-java</artifactId>
    <version>0.1.0</version>
</dependency>
```

In order to know how to create a new Vance Connector, check out the [API Documentation][api].

If you want to know more about the default HTTP implementation for handling requests or sending CloudEvents out, check out the 
[HTTP implementations][http].

Developers better have a basic familiarity with the [CloudEvents Specification][ce] and [CloudEvents java-sdk][ce-sdk] before they start writing a connector.

### Connector Examples

Here are some connector examples developed by cdk-java.

| Connector         | Type          | Description |
|:-------------|:------------------|:------|
| [sink-example]    | sink | Use this example to write your first sink connector.  |
| [source-example] | source   | Use this example to write your first source connector. |

Use them as samples when you want to write a sink or source connector.

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md
[api]: https://linkall-labs.github.io/cdk-java/api.html
[http]: https://linkall-labs.github.io/cdk-java/http.html
[sink-example]: https://github.com/linkall-labs/cdk-java/tree/main/examples/sink-example
[source-example]: https://github.com/linkall-labs/cdk-java/tree/main/examples/source-example
[ce]: https://github.com/cloudevents/spec
[ce-sdk]: https://github.com/cloudevents/sdk-java