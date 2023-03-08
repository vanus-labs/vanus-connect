---
title: Tencent Cloud SCF
---

# Tencent Cloud SCF Sink

## Introduction

The Tencent Cloud SCF Sink is a [Vanus Connector][vc] that aims to handle incoming CloudEvents in a way that extracts
the `data` part of the original event and Tencent Cloud invoke SCF function

## Quickstart

### Prerequisites

- Have a container runtime (i.e., docker).
- Have a Tencent Cloud Account and SCF function

### create config file

```shell
cat << EOF > config.yml
port: 8080
secret:
  secret_id: ABID570jkkngFWl7uY3QchbdUXVIuNisywoA
  secret_key: xxxxxx
function:
  name: "xxxxxxxxx"
  region: "ap-beijing"
  namespace: "default"
EOF
```

| Name               | Required | Default | Description                                          |
|:-------------------|:--------:|:-------:|------------------------------------------------------|
| port               |    No    |  8080   | the port which the Tencent Cloud SCF Sink listens on |
| secret.secret_id   |   YES    |         | the Tencent Cloud cam secretId                       |
| secret.secret_key  |   YES    |         | the Tencent Cloud SCF cam secretKey                  |
| function.region    | **YES**  |         | which region the function was created                |
| function.name      | **YES**  |         | function name will be invoked                        |
| function.namespace | **YES**  |         | which namespace the function was created             |

The Tencent Cloud SCF Sink tries to find the config file at `/vanus-connect/config/config.yml` by default. You can
specify the position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm \
  -p 31080:8080 \
  -v ${PWD}:/vanus-connect/config \
  --name sink-tencentcloud-scf public.ecr.aws/vanus/connector/sink-tencentcloud-scf
```

### Test

Open a terminal and use following command to send a CloudEvent to the Sink.

```shell
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

you can see logs in [SCF console](https://console.cloud.tencent.com/scf)

![log.png](https://github.com/vanus-labs/vanus-connect/blob/main/connectors/sink-tencentcloud-scf/scf-log.png?raw=true)

### Clean resource

```shell
docker stop sink-tencentcloud-scf
```

## Run in Kubernetes

```shell
kubectl apply -f sink-tencentcloud-scf.yaml
```

```yml
apiVersion: v1
kind: Service
metadata:
  name: sink-tencentcloud-scf
  namespace: vanus
spec:
  selector:
    app: sink-tencentcloud-scf
  type: ClusterIP
  ports:
    - port: 8080
      name: sink-tencentcloud-scf
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: sink-tencentcloud-scf
  namespace: vanus
data:
  config.yml: |-
    port: 8080
    secret:
      secret_id: ABID570jkkngFWl7uY3QchbdUXVIuNisywoA
      secret_key: xxxxxx
    function:
      name: "xxxxxx"
      region: "ap-beijing"
      namespace: "default"
    debug: false

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sink-tencentcloud-scf
  namespace: vanus
  labels:
    app: sink-tencentcloud-scf
spec:
  selector:
    matchLabels:
      app: sink-tencentcloud-scf
  replicas: 1
  template:
    metadata:
      labels:
        app: sink-tencentcloud-scf
    spec:
      containers:
        - name: sink-tencentcloud-scf
          image: public.ecr.aws/vanus/connector/sink-tencentcloud-scf
          imagePullPolicy: Always
          ports:
            - name: http
              containerPort: 8080
          volumeMounts:
            - name: config
              mountPath: /vanus-connect/config
      volumes:
        - name: config
          configMap:
            name: sink-tencentcloud-scf
```

## Integrate with Vanus

This section shows how a sink connector can receive CloudEvents from a
running [Vanus cluster](https://github.com/vanus-labs/vanus).

1. Run the sink-tencentcloud-scf.yaml

```shell
kubectl apply -f sink-tencentcloud-scf.yaml
```

2. Create an eventbus

```shell
vsctl eventbus create --name quick-start
```

3. Create a subscription (the sink should be specified as the sink service address or the host name with its port)

```shell
vsctl subscription create \
  --name quick-start \
  --eventbus quick-start \
  --sink 'http://sink-tencentcloud-scf:8080'
```

[vc]: https://docs.vanus.ai/introduction/concepts#vanus-connect
