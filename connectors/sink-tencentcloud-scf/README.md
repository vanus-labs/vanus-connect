---
title: Tencent SCF
---

# Tencent SCF Sink

## Introduction

This connector for invoking SCF function with events.

## Quickstart

### create config file

```shell
cat << EOF > config.yml
port: 8080
function:
  name: "xxxxxxxxx"
  region: "ap-beijing"
  namespace: "default"
EOF
```

For full configuration, you can see [config](#config) section.

### create secret file

```shell
cat << EOF > secret.yml
secret_id: "xxxx"
secret_key: "xxxxx"
EOF
```

### start

```shell
docker run -d --rm \
  --network host \
  -v ${PWD}:/vance/config \
  -v ${PWD}:/vance/secret \
  --name sink-tencentcloud-scf public.ecr.aws/vanus/connector/sink-tencentcloud-scf:dev
```

### send an event to sink

```bash
curl --location --request POST 'localhost:8080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "local",
    "specversion": "1.0",
    "type": "xxxx",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "data": {
        "Records": [
            {
                "cos": "xxxx"
            },
            {
                "cos": "yyy"
            }
        ]
    }
}'
```

the `data` will be as payload to invoke function

### see logs in SCF console

![log.png](https://github.com/linkall-labs/vance/blob/main/connectors/sink-tencentcloud-scf/scf-log.png)

### clean resource

```shell
docker stop sink-tencentcloud-scf
```

## Configuration

### config

```yml
port: 8080
function:
  name: "vanus-cos-source-function-3513950818025804220"
  region: "ap-beijing"
  namespace: "default"
debug: false  
```

| Name               | Required | Default | Description                              |
|:-------------------|:--------:|:-------:|------------------------------------------|
| port               |    No    |  8080   | which port for listening                 |
| function.region    | **YES**  |    -    | which region the function was created    |
| function.name      | **YES**  |    -    | function name will be invoked            |
| function.namespace | **YES**  |    -    | which namespace the function was created |
| debug              |    NO    |  false  | if print debug log                       |

### secret


| Name       | Required | Default | Description                |
|:-----------|:--------:|:-------:|----------------------------|
| secret_id  | **YES**  |    -    | SecretID of Tencent Cloud  |
| secret_key | **YES**  |    -    | SecretKey of Tencent Cloud |

## Deploy

### using k8s(recommended)

```yml
apiVersion: v1
kind: Service
metadata:
  name: sink-tencentcloud-function
  namespace: vanus
spec:
  selector:
    app: sink-tencentcloud-function
  type: NodePort
  ports:
    - port: 8080
      targetPort: 8080
      nodePort: 32555
      name: sink-tencentcloud-function
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: sink-tencentcloud-function
  namespace: vanus
data:
  config.yml: |-
    port: 8080
    function:
      name: "xxxxxx"
      region: "ap-beijing"
      namespace: "default"
    debug: false

---
apiVersion: v1
kind: Secret
metadata:
  name: sink-tencentcloud-function
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
  name: sink-tencentcloud-function
  namespace: vanus
  labels:
    app: sink-tencentcloud-function
spec:
  selector:
    matchLabels:
      app: sink-tencentcloud-function
  replicas: 1
  template:
    metadata:
      labels:
        app: sink-tencentcloud-function
    spec:
      containers:
        - name: sink-tencentcloud-function
          image: public.ecr.aws/vanus/connector/sink-tencentcloud-function:dev
          imagePullPolicy: Always
          ports:
            - name: http
              containerPort: 8080
          volumeMounts:
            - name: config
              mountPath: /vance/config
            - name: secret
              mountPath: /vance/secret
      volumes:
        - name: secret
          secret:
            secretName: sink-tencentcloud-function
        - name: config
          configMap:
            name: sink-tencentcloud-function
```

### using vance Operator

coming soon