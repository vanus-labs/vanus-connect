---
title: <name>
---

# <name> Source

## Introduction

The <name> Source is a [Vanus Connector](https://www.vanus.dev/introduction/concepts#vanus-connect) which aims to convert incoming data ...

<optional: incoming request/message example>

which is converted to

</optional>

```json
{
 <example event converted by this source>
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

Assuming you use [Vanus](https://github.com/linkall-labs/vanus) as your CloudEvent receiver, if you have other receiver,
config the target to your own address

Replace `<vanus_gateway_url>`, `<port>`, and `<eventbus>` to yours.

```shell
cat << EOF > config.yml
# use local Sink Display container to verify events
target: http://localhost:31081
<other_configs>
EOF
```

<optional if have secret configs>
### create the secret file
...
</optional>

### start with Docker

<mapping 8080(container port) to 31080(host port) to avoid conflict.>

```shell
docker run -d --rm \
  -p 31080:8080 \
  -v ${PWD}:/vance/config \
  --name source-<name> public.ecr.aws/vanus/connector/source-<name>:latest
```

### Test

Start display Sink with the following command, which received events that this source made:
```shell
docker run -d --rm \
  -p 31081:8080 \
  --name sink-display public.ecr.aws/vanus/connector/sink-display
```

<do some operation>

use `docker logs sink-display` to view events

```json
{
 "id" : "ef26ed7b-9377-4bf5-b8d4-4fc6347e4fa2",
 "source" : "kafka.host.docker.internal.topic1",
 "specversion" : "V1",
 "type" : "kafka.message",
 "datacontenttype" : "plain/text",
 "time" : "2022-12-05T09:00:42.618Z",
 "data" : "Hello world!"
}
```

### Clean

```shell
docker stop source-<name> sink-display
```

## How to use




<optional>
### Extension Attributes
The <name> Source defines following [CloudEvents Extension Attributes](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#extension-context-attributes)

|    Attribute     |  Type   | Description                                                                                                                      |
|:----------------:|:-------:|:---------------------------------------------------------------------------------------------------------------------------------|
...

<optional>
### Data 
<explain the structure of data>

### Run in Kubernetes

```yaml
<content>
```