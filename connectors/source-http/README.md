---
title: HTTP
---

# HTTP Source

## Introduction

The HTTP Source is a [Vance Connector](../vc) which aims to convert incoming HTTP Request to a CloudEvent.

For example, if the incoming HTTP Request looks like:

```bash
curl --location --request POST 'localhost:8080/webhook?source=123&id=abc&type=456&subject=def&test=demo' \
--header 'Content-Type: text/plain' \
--data-raw '{
    "test":"demo"
}'
```

which is converted to

```json
{
  "specversion":"1.0",
  "id":"abc",
  "source":"123",
  "type":"456",
  "subject":"def",
  "datacontenttype":"application/json",
  "data":{
    "body":{
      "xxxx":"aaa"
    },
    "headers":{
      "Accept":"*/*",
      "Accept-Encoding":"gzip, deflate, br",
      "Connection":"keep-alive",
      "Content-Length":"20",
      "Content-Type":"text/plain",
      "Host":"localhost:12321",
      "Postman-Token":"6abb2398-d4f3-4eb1-9c57-65b6934b84c1",
      "User-Agent":"PostmanRuntime/7.29.2"
    },
    "method":"POST",
    "path":"/webhook",
    "query_args":{
      "id":"abc",
      "source":"123",
      "subject":"def",
      "test":"demo",
      "type":"456"
    }
  },
  "xvhttpremoteip":"::1",
  "xvhttpremoteaddr":"[::1]:57822",
  "xvhttpbodyisjson":true
}
```

## Quick Start

in this section, we show how HTTP Source convert HTTP request(made by cURL) to CloudEvent.

### Create Config file

Assuming you use Vanus(https://github.com/linkall-labs/vanus) as CloudEvent receiver, if you have other receiver,
just set target to your endpoint.

```shell
cat << EOF > config.yml
# change url, port and eventbus to yours
target: http://<url>:<port>/gateway/<eventbus>
EOF
```

### Start Using Docker

mapping 8080 to 31080 in order to avoid port conflict.

```shell
docker run -d -p 31080:8080 --rm \
  -v ${PWD}:/vance/config \
  --name source-http public.ecr.aws/vanus/connector/source-http:latest
```

### Test

```shell
curl --location --request POST 'localhost:31080/webhook?source=123&id=abc&type=456&subject=def' \
--header 'Content-Type: text/plain' \
--data-raw '{
    "test":"demo"
}'
```

now, you could use `vsctl event get <eventbus>` to view event just sent. If you can't see event you sent,
try to use `--offset` to get event. (`vsctl` default retrieves event from earliest)

```
~> vsctl event get <eventbus>
+-----+----------------------------------------------------------------+
|     | Context Attributes,                                            |
|     |   specversion: 1.0                                             |
|     |   type: 456                                                    |
|     |   source: 123                                                  |
|     |   subject: def                                                 |
|     |   id: abc                                                      |
|     |   time: 2022-12-11T11:42:50.135762Z                            |
|     |   datacontenttype: application/json                            |
|     | Extensions,                                                    |
|     |   xvanuseventbus: wwf                                          |
|     |   xvanuslogoffset: AAAAAAAAAAQ=                                |
|     |   xvanusstime: 2022-12-11T11:42:50.874Z                        |
|     |   xvhttpbodyisjson: true                                       |
|     |   xvhttpremoteaddr: [::1]:57822                                |
|     |   xvhttpremoteip: ::1                                          |
|     | Data,                                                          |
|     |   {                                                            |
|     |     "body": {                                                  |
|     |       "xxxx": "aaa"                                            |
|     |     },                                                         |
|  0  |     "headers": {                                               |
|     |       "Accept": "*/*",                                         |
|     |       "Accept-Encoding": "gzip, deflate, br",                  |
|     |       "Connection": "keep-alive",                              |
|     |       "Content-Length": "20",                                  |
|     |       "Content-Type": "text/plain",                            |
|     |       "Host": "localhost:12321",                               |
|     |       "Postman-Token": "6abb2398-d4f3-4eb1-9c57-65b6934b84c1", |
|     |       "User-Agent": "PostmanRuntime/7.29.2"                    |
|     |     },                                                         |
|     |     "method": "POST",                                          |
|     |     "path": "/webhook",                                        |
|     |     "query_args": {                                            |
|     |       "id": "abc",                                             |
|     |       "source": "123",                                         |
|     |       "subject": "def",                                        |
|     |       "test": "demo",                                          |
|     |       "type": "456"                                            |
|     |     }                                                          |
|     |   }                                                            |
|     |                                                                |
+-----+----------------------------------------------------------------+
```

### Clean

```shell
docker stop source-http
```

## How to use

### Configuration

The default path is `/vance/config/config.yml`. if you want to change the default path, you can set env `CONNECTOR_CONFIG` to
tell HTTP Source.


| Name   | Required | Default | Description                         |
|:-------|:--------:|:-------:|-------------------------------------|
| target | **YES**  |    -    | the endpoint of CloudEvent sent to. |

### Attributes

#### Changing Default Required Attributes
if you want change default attributes of `id`,`source`, `type`, and `subject`(defined by [CloudEvents](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#required-attributes))
to you own, you could use query parameter to set them.

| Attribute  |      Default       | Query Parameter | Example                                 |
|:----------:|:------------------:|:----------------|:----------------------------------------|
|     id     |        UUID        | ?id=xxx         | http://url:port/webhook?id=xxxx         |
|   source   | vanus-http-source  | ?source=xxx     | http://url:port/webhook?source=xxxx     |
|    type    | naive-http-request | ?type=xxx       | http://url:port/webhook?type=xxxx       |
|  subject   |       empty        | ?subject=xxx    | http://url:port/webhook?subject=xxxx    |
| dataschema |       empty        | ?dataschema=xxx | http://url:port/webhook?dataschema=xxxx |

`datacontenttype` will be auto infer based on request body, if body can be converted to `JSON`, the `application/json` will be set,
otherwise `text/plain` will be set.

#### Extension Attributes
HTTP Source provides some [CloudEvents Extension Attributes](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#extension-context-attributes)

|    Attribute     |  Type   | Description                                                                                                                      |
|:----------------:|:-------:|:---------------------------------------------------------------------------------------------------------------------------------|
| xvhttpbodyisjson | boolean | HTTP Sink will validate if request body is JSON format data, if it is, this attribute is `true`, otherwise `false`               |
|  xvhttpremoteip  | string  | The IP of the request from where, if the request was through reverse-proxy like Nginx, the value may be not the original IP      |
| xvhttpremoteaddr | string  | The address of the request from where, if the request was through reverse-proxy like Nginx, the value may be not the original IP |


## Run in Kubernetes

```yaml
apiVersion: v1
kind: Service
metadata:
  name: source-http
  namespace: vanus
spec:
  selector:
    app: source-http
  type: ClusterIP
  ports:
    - port: 8080
      name: source-http
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: source-http
  namespace: vanus
data:
  config.yml: |-
    target: http://<url>:<port>/gateway/<eventbus>

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: source-http
  namespace: vanus
  labels:
    app: source-http
spec:
  selector:
    matchLabels:
      app: source-http
  replicas: 1
  template:
    metadata:
      labels:
        app: source-http
    spec:
      containers:
        - name: source-http
          image: public.ecr.aws/vanus/connector/source-http:latest
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "512Mi"
              cpu: "500m"
          imagePullPolicy: Always
          env:
            - name: LOG_LEVEL
              value: INFO
          volumeMounts:
            - name: config
              mountPath: /vance/config
      volumes:
        - name: config
          configMap:
            name: source-http
```

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md