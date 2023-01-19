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

```shell
cat << EOF > config.yml
# use local Sink Display container to verify events
target: http://localhost:31081
<other_configs>
EOF
```


### Start with Docker

```shell
docker run --rm \
  -p 31080:8080 \
  -v ${PWD}:/vanus-connect/config \
  --name source-<name> public.ecr.aws/vanus/connector/source-<name>:latest
```

### Test

Open a terminal and use the following command to run a Display sink, which receives and prints CloudEvents out.
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

## Source details




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