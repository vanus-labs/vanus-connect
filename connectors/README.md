# Vanus connectors

## Purpose

The article introduces how to add a Vanus connector.

## Requirements

To add a new connector you need to:

1. Develop the connector.
2. Document how to use the connector.

We provide CDK Java and CDK Go for you to quickly develop a connector.

## CDK Java

### Source

1. Create a java module name like source-example.
2. Import cdk java in the pom.xml.
3. Implement the interface `com.linkall.cdk.connector.Source`.

here is [example source](https://github.com/vanus-labs/cdk-java/tree/main/examples/source-example)

### Sink

1. Create a java module name like sink-example.
2. Import cdk java in the pom.xml.
3. Implement the interface `com.linkall.cdk.connector.Sink`.

here is [example sink](https://github.com/vanus-labs/cdk-java/tree/main/examples/sink-example)

## CDK Go

### Source

1. Create a go module name like source-example.
2. Import cdk go in the go.mod.
3. Implement the interface [Source](https://github.com/vanus-labs/cdk-go/blob/main/connector/source.go).

here is [example source](https://github.com/vanus-labs/cdk-go/tree/main/examples/source-example)

### Sink

1. Create a go module name like sink-example.
2. Import cdk java in the pom.xml.
3. Implement the interface [Sink](https://github.com/vanus-labs/cdk-go/blob/main/connector/sink.go).

here is [example sink](https://github.com/vanus-labs/cdk-go/tree/main/examples/sink-example)