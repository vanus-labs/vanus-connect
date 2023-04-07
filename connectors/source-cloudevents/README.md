---
title: CloudEvents
---

# CloudEvents Source

## Introduction

The CloudEvents Source is a [Vanus Connector][vc] which aims to receive CloudEvents and send to target.

## Quick Start

This section shows how CloudEvents Source receive CloudEvents and send to target.

### Create Config file

```shell
cat << EOF > config.yml
target: http://localhost:31081
EOF
```

| Name          | Required | Default | Description                                                       |
|:--------------|:--------:|:-------:|:------------------------------------------------------------------|
| target        |   YES    |         | the target URL to send CloudEvents                                |
| port          |    NO    |  8080   | the port to receive CloudEvents                                   |
| path          |    NO    |         | the CloudEvents source http path to receive event                 |
| header        |    NO    |         | the CloudEvents source http header                                |
| auth.username |    NO    |         | the CloudEvents source http authentication by basic auth username |
| auth.password |    NO    |         | the CloudEvents source http authentication by basic auth password |

The CloudEvents Source tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify
the position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-cloudevents public.ecr.aws/vanus/connector/source-cloudevents
```

### Test

Open a terminal and use the following command to run a Display sink, which receives and prints CloudEvents.

```shell
docker run -it --rm \
  -p 31081:8080 \
  --name sink-display public.ecr.aws/vanus/connector/sink-display
```

Make sure the `target` value in your config file is `http://localhost:31081` so that the Source can send the CloudEvents
to the Display Sink.

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-cloudevents public.ecr.aws/vanus/connector/source-cloudevents
```

Open a terminal and use the following command to send an event request to CloudEvents Source

```shell
curl --location --request POST 'localhost:8080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
  "id" : "42d5b039-daef-4071-8584-e61df8fc1354",
  "source" : "quickstart",
  "specversion" : "1.0",
  "type" : "quickstart",
  "datacontenttype" : "application/json",
  "time" : "2023-01-26T10:38:29.345Z",
  "data" : "quickstart"
}'
```

Here is the sort of CloudEvent you should expect to receive in the Display Sink:

```json
{
  "id": "42d5b039-daef-4071-8584-e61df8fc1354",
  "source": "quickstart",
  "specversion": "1.0",
  "type": "quickstart",
  "datacontenttype": "application/json",
  "time": "2023-01-26T10:38:29.345Z",
  "data": "quickstart"
}
```

### Clean

```shell
docker stop source-cloudevents sink-display
```

## Run in Kubernetes

```shell
kubectl apply -f source-cloudevents.yaml
```

```yaml
apiVersion: v1
kind: Service
metadata:
  name: source-cloudevents
  namespace: vanus
spec:
  selector:
    app: source-cloudevents
  type: ClusterIP
  ports:
    - port: 8080
      name: source-cloudevents
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: source-cloudevents
  namespace: vanus
data:
  config.yml: |-
    target: http://<url>:<port>/gateway/<eventbus>

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: source-cloudevents
  namespace: vanus
  labels:
    app: source-cloudevents
spec:
  selector:
    matchLabels:
      app: source-cloudevents
  replicas: 1
  template:
    metadata:
      labels:
        app: source-cloudevents
    spec:
      containers:
        - name: source-cloudevents
          image: public.ecr.aws/vanus/connector/source-cloudevents:latest
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "512Mi"
              cpu: "500m"
          imagePullPolicy: Always
          volumeMounts:
            - name: config
              mountPath: /vanus-connect/config
      volumes:
        - name: config
          configMap:
            name: source-cloudevents
```

## Integrate with Vanus

This section shows how a source connector can send CloudEvents to a
running [Vanus cluster](https://github.com/vanus-labs/vanus).

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

3. Update the target config of the CloudEvents Source

```yaml
target: http://192.168.49.2:30001/gateway/quick-start
```

4. Run the CloudEvents Source

```shell
kubectl apply -f source-cloudevents.yaml
```

[vc]: https://docs.vanus.ai/introduction/concepts#vanus-connect
