# Display Sink 

## Introduction

The Display Sink is a [Vance Connector][vc] which prints received CloudEvents. This is commonly used as a logger to check incoming data.

For example, it will print the incoming CloudEvents looks like:

```json
{
  "id" : "42d5b039-daef-4071-8584-e61df8fc1354",
  "source" : "vance-http-source",
  "specversion" : "V1",
  "type" : "http",
  "datacontenttype" : "application/json",
  "time" : "2022-05-17T18:44:02.681+08:00",
  "data" : {
    "myData" : "simulation event data <1>"
  }
}
```

## Display Sink Configs

### Config

| Name | Required | Default | Description                                                   |
|:-----|:---------|:--------|:--------------------------------------------------------------|
| port | false    | 8080    | port is used to specify the port Display Sink is listening on |

## Display Sink Image

> vancehub/sink-display

## Deploy

### Docker

#### create config file

refer [config](#Config) to create `config.yaml`. for example:

```yaml
"port": 8080
```

#### run

```shell
 docker run --rm -v ${PWD}:/vance/config -v ${PWD}:/vance/secret vancehub/sink-display
```

### K8S

```shell
  kubectl apply -f sink-display.yaml
```

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md
[config]: https://github.com/linkall-labs/vance-docs/blob/main/docs/connector.md