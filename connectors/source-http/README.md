# HTTP Source

## Introduction

The HTTP Source is a [Vance Connector](../README.md) which aims to convert incoming HTTP Request to a CloudEvent.

For example, if the incoming HTTP Request looks like:

```bash
curl --location --request POST 'localhost:8080/webhook?source=123&id=abc&type=456&subject=def' \
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
    "test":"demo"
  },
  "xvhttpuseragent":"curl/7.77.0",
  "xvhttpremoteip":"::1",
  "xvhttpremoteaddr":"[::1]:62734"
}
```

## Quick Start

in this section, we show how to use HTTP Source push a message to your group chat.

### Create Config file

```shell
cat << EOF > config.yml
# Assuming you use Vanus(https://github.com/linkall-labs/vanus) as CloudEvent recevier, if you have other receiver, just set target to your endpoint.
target: http://<url>:<port>/gateway/<eventbus>
EOF
```

### Start Using Docker

```shell
# mapping 8080 to 31080 in order to avoid port conflict.
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

now, you could use `vsctl get event <eventbus> --number 10` to view event just sent

```
~> vsctl get event <eventbus> --number 10
+-----+-------------------------------------------------+
|     | Context Attributes,                             |
|     |   specversion: 1.0                              |
|     |   type: 456                                     |
|     |   source: 123                                   |
|     |   subject: def                                  |
|     |   id: abc                                       |
|     |   time: 2022-12-10T09:59:10.806608Z             |
|     |   datacontenttype: application/json             |
|     | Extensions,                                     |
|  0  |   xvanuseventbus: wwf                           |
|     |   xvanuslogoffset: AAAAAAAAAAg=                 |
|     |   xvanusstime: 2022-12-10T09:59:11.574Z         |
|     |   xvhttpremoteaddr: [::1]:62734                 |
|     |   xvhttpremoteip: ::1                           |
|     |   xvhttpuseragent: curl/7.77.0                  |
|     | Data,                                           |
|     |   {                                             |
|     |     "xxxx": "aaa"                               |
|     |   }                                             |
|     |                                                 |
+-----+-------------------------------------------------+
```

### Clean

```shell
docker stop sink-feishu
```

## Configuration

The default path is `/vance/config/config.yml`. if you want to change the default path, you can set env `CONNECTOR_CONFIG` to
tell HTTP Source.

| Name   | Required | Default | Description                         |
|:-------|:--------:|:-------:|-------------------------------------|
| target | **YES**  |    -    | the endpoint of CloudEvent sent to. |

## Run in Kubernetes

```yaml
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
          imagePullPolicy: Always
          volumeMounts:
            - name: config
              mountPath: /vance/config
      volumes:
        - name: config
          configMap:
            name: source-http
```
