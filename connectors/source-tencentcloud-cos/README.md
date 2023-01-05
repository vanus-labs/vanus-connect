---
title: Tencent COS
---


# Tencent COS Source Connector

## Introduction

This connector for capturing Tencent COS event

## Quickstart

### create config file

```shell
cat << EOF > config.yml
target: "35.87.170.130:8080"
eventbus: "xxx"
bucket:
  endpoint: "xxx.cos.<region>.myqcloud.com"
function:
  region: <region>
EOF
```
For full configuration, you can see config section.

### create secret file

```shell
cat << EOF > secret.yml
secret_id: "xxxx"
secret_key: "xxxxx"
EOF
```

### start using docker

```shell
docker run -d --rm \
  --network host \
  -v ${PWD}:/vance/config \
  -v ${PWD}:/vance/secret \
  --name source-tencentcloud-cos public.ecr.aws/vanus/connector/source-tencentcloud-cos:dev
```

### upload a file to your bucket

open COS console in browser, upload a file into your target bucket.

### see event was captured

An event like below is received:

```json
{
    "Records": [{
        "cos": {
            "cosSchemaVersion": "1.0",
            "cosObject": {
                "url": "http://testpic-1253970026.cos.ap-chengdu.myqcloud.com/testfile",
                "meta": {
                    "x-cos-request-id": "NWMxOWY4MGFfMjViMjU4NjRfMTUyMVxxxxxxxxx=",
                    "Content-Type": "",
                    "x-cos-meta-mykey": "myvalue"
                },
                "vid": "",
                "key": "/1253970026/testpic/testfile",
                "size": 1029
            },
            "cosBucket": {
                "region": "cd",
                "name": "testpic",
                "appid": "1253970026"
            },
            "cosNotificationId": "unkown"
        },
        "event": {
            "eventName": "cos:ObjectCreated:*",
            "eventVersion": "1.0",
            "eventTime": 1545205770,
            "eventSource": "qcs::cos",
            "requestParameters": {
                "requestSourceIP": "192.168.15.101",
                "requestHeaders": {
                    "Authorization": "q-sign-algorithm=sha1&q-ak=xxxxxxxxxxxxxx&q-sign-time=1545205709;1545215769&q-key-time=1545205709;1545215769&q-header-list=host;x-cos-storage-class&q-url-param-list=&q-signature=xxxxxxxxxxxxxxx"
                }
            },
            "eventQueue": "qcs:0:scf:cd:appid/1253970026:default.printevent.$LATEST",
            "reservedInfo": "",
            "reqid": 179398952
        }
    }]
}
```

please see [COS Trigger](https://cloud.tencent.com/document/product/583/9707) to understanding the structure of events.

### clean resource

```shell
docker stop source-tencentcloud-cos
```

## Configuration

### config

```yml
target: "x.x.x.x:8080"
eventbus: "xxxx"
bucket:
  endpoint: "xxxx.cos.ap-beijing.myqcloud.com"
function:
  region: "ap-beijing"
  name: "xxxx"
  namespace: "default"
  code:
    bucket: "vanus-1253760853"
    region: "ap-beijing"
    path: "/vanus/cos-source/dev/main.zip"
debug: false
secret_id: "xxxx"
secret_key: "xxxxx"
```

| Name                 | Required |              Default               | Description                                                                      |
|:---------------------|:--------:|:----------------------------------:|----------------------------------------------------------------------------------|
| target               | **YES**  |                 -                  | Target URL will send CloudEvents to                                              |
| eventbus             | **YES**  |                 -                  | target eventbus                                                                  |
| bucket.endpoint      | **YES**  |                 -                  | which bucket you want to capture.                                                |
| function.region      | **YES**  |                 -                  | which region the helper function will be deployed, suggest keep same with bucket |
| function.name        |    NO    | vanus-cos-source-function-<number> |                                                                                  |
| function.namespace   |    NO    |              default               | which namespace the function created to                                          |
| function.code.bucket |    NO    |          vanus-1253760853          | Not recommended to modify                                                        |
| function.code.region |    NO    |             ap-beijing             | Not recommended to modify                                                        |
| function.code.path   |    NO    |   /vanus/cos-source/dev/main.zip   | Not recommended to modify                                                        |
| debug                |    NO    |               false                | if print debug log                                                               |

### secret


| Name       | Required | Default | Description                |
|:-----------|:--------:|:-------:|----------------------------|
| secret_id  | **YES**  |    -    | SecretID of Tencent Cloud  |
| secret_key | **YES**  |    -    | SecretKey of Tencent Cloud |


## Deploy

### using k8s(recommended)

```yml
apiVersion: v1
kind: ConfigMap
metadata:
  name: source-tencentcloud-cos
  namespace: vanus
data:
  config.yml: |-
    target: "xxxx"
    eventbus: "xxxxx"
    bucket:
      endpoint: "xxxxx.cos.ap-beijing.myqcloud.com"
    function:
      region: "ap-beijing"

---
apiVersion: v1
kind: Secret
metadata:
  name: source-tencentcloud-cos
  namespace: vanus
type: Opaque
data:
  # cat secret.yml | base64
  secret.yml: |
    xxxxx
immutable: true
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: source-tencentcloud-cos
  namespace: vanus
  labels:
    app: source-tencentcloud-cos
spec:
  selector:
    matchLabels:
      app: source-tencentcloud-cos
  replicas: 1
  template:
    metadata:
      labels:
        app: source-tencentcloud-cos
    spec:
      containers:
        - name: source-tencentcloud-cos
          image: public.ecr.aws/vanus/connector/source-tencentcloud-cos:dev
          imagePullPolicy: Always
          volumeMounts:
            - name: config
              mountPath: /vance/config
            - name: secret
              mountPath: /vance/secret
      volumes:
        - name: secret
          secret:
            secretName: source-tencentcloud-cos
        - name: config
          configMap:
            name: source-tencentcloud-cos
```

### using vance Operator

coming soon

## Event Structure

| Field  | Required | Description                                                            |
|--------|:--------:|------------------------------------------------------------------------|
| id     | **YES**  | random UUID                                                            |
| source | **YES**  | function name                                                          |
| type   | **YES**  | tencent-cloud-cos-event                                                |
| time   | **YES**  | the time of this event generated with RFC3339 encoding                 |
| data   | **YES**  | see [COS Trigger](https://cloud.tencent.com/document/product/583/9707) |
