---
title: <name>
---

# <name> Source

## Introduction

The <name> Source is a [Vanus Connector][vc] which aims to convert incoming data ...

<optional: incoming request/message example>

which is converted to

</optional>

```json
{
  <example
  event
  converted
  by
  this
  source>
}
```

## Quick Start

This section shows how <name> Source convert <xxxx> to a CloudEvent.

<optional prerequisites but recommended>

### Prerequisites

- Have a container runtime (i.e., docker).
- ...
  </optional>

### Create the config file

```shell
cat << EOF > config.yml
# use local Sink Display container to verify events
target: http://localhost:31081
<other_configs>
EOF
```

| Name   | Required | Default | Description                                                 |
|:-------|:---------|:--------|:------------------------------------------------------------|
| target | YES      | ""      | the target URL which <name> Source will send CloudEvents to |

...

The <name> Source tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the
position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-<name> public.ecr.aws/vanus/connector/source-<name>
```

### Test

Open a terminal and use the following command to run a Display sink, which receives and prints CloudEvents.

```shell
docker run -it --rm \
  -p 31081:8080 \
  --name sink-display public.ecr.aws/vanus/connector/sink-display
```

Make sure the `target` value in your config file is `http://localhost:31081` so that the Source can send CloudEvents to
our Display Sink.

<do some operation>

Here is the sort of CloudEvent you should expect to receive in the Display Sink:

```json
{
  "id": "ef26ed7b-9377-4bf5-b8d4-4fc6347e4fa2",
  "source": "kafka.host.docker.internal.topic1",
  "specversion": "1.0",
  "type": "kafka.message",
  "datacontenttype": "plain/text",
  "time": "2022-12-05T09:00:42.618Z",
  "data": "Hello world!"
}
```

### Clean

```shell
docker stop source-<name> sink-display
```

## Source details

<optional>
### Extension Attributes
The <name> Source defines following [CloudEvents Extension Attributes](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#extension-context-attributes)

| Attribute | Type | Description                                                                                                                      |
|:----------|:-----|:---------------------------------------------------------------------------------------------------------------------------------|

...
</optional>

<optional>
### Data 
<optional the structure of data>
</optional>

## Run in Kubernetes

```shell
kubectl apply -f source-<name>.yaml
```

```yaml
<content>
```

## Integrate with Vanus

This section shows how a source connector can send CloudEvents to a
running [Vanus cluster](https://github.com/linkall-labs/vanus).

### Prerequisites

- Have a running K8s cluster
- Have a running Vanus cluster
- Vsctl Installed

1. Export the VANUS_GATEWAY environment variable (the ip should be a host-accessible address of the vanus-gateway
   service)

```shell
export VANUS_GATEWAY=192.168.49.2:30001
```

2. Create an eventbus

```shell
vsctl eventbus create --name quick-start
```

3. Update the target config of the <name> Source

```yaml
target: http://192.168.49.2:30001/gateway/quick-start
```

4. Run the <name> Source

```shell
docker run --network=host \
  --rm \
  -v ${PWD}:/vanus-connect/config \
  --name source-<name> public.ecr.aws/vanus/connector/source-<name>:latest
```

[vc]: https://www.vanus.dev/introduction/concepts#vanus-connect
