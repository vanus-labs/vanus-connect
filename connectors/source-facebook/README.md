---
title: Facebook
---

# Facebook Source

## Introduction

The Facebook Source is a [Vanus Connector][vc] which aims to convert a Facebook event to a CloudEvent.

For example, the incoming Facebook page event looks like:

```json
[
  {
    "entry": [
      {
        "changes": [
          {
            "field": "feed",
            "value": {
              "from": {
                "id": "987654321",
                "name": "Cinderella Hoover"
              },
              "item": "post",
              "post_id": "{page-post-id}",
              "verb": "add",
              "created_time": 1520544814,
              "is_hidden": false,
              "message": "It is Thursday and I want to eat cake."
            }
          }
        ],
        "id": "123456789",
        "time": 1520544816
      }
    ],
    "object": "page"
  }
]
```

which is converted to:

```json
{
  "specversion": "1.0",
  "id": "f7c7821d-a79f-45cb-a78d-b35cae51027c",
  "source": "vanus.facebook",
  "type": "page",
  "datacontenttype": "application/json",
  "time": "2023-03-29T10:47:14.90043Z",
  "changeitem": "post",
  "pageid": "123456789",
  "fields": "feed",
  "data": {
    "field": "feed",
    "value": {
      "created_time": 1520544814,
      "from": {
        "id": "987654321",
        "name": "Cinderella Hoover"
      },
      "is_hidden": false,
      "item": "post",
      "message": "It is Thursday and I want to eat cake.",
      "post_id": "{page-post-id}",
      "verb": "add"
    }
  }
}
```

## Quick Start

This section will show you how to use Facebook Source to convert a facebook webhook event to a CloudEvent.

### Create Config file

```shell
cat << EOF > config.yml
target: http://localhost:31081
port: 8082
verify_token: verify_token
EOF
```

| Name         | Required | Default | Description                                                 |
|:-------------|:--------:|:-------:|:------------------------------------------------------------|
| target       |   YES    |         | the target URL to send CloudEvents                          |
| port         |    NO    |  8080   | the port to receive webhook event request                   |
| verify_token |   YES    |         | the Facebook webhook verify token                           |
| app_secret   |   YES    |         | the app secret for check webhook header X-Hub-Signature-256 |

The Facebook Source tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the
position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-facebook public.ecr.aws/vanus/connector/source-facebook
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

Configure [Facebook webhook][fb-webhook] and add Page webhook then click test, then the Display sink will receive a
CloudEvents.

### Clean

```shell
docker stop source-fackbook sink-display
```

## Source details

### Attributes

#### Extension Attributes

The Facebook Source defines
following [CloudEvents Extension Attributes](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#extension-context-attributes)

| Attribute  |  Type  | Description                                       |
|:----------:|:------:|:--------------------------------------------------|
|   pageid   | string | the facebook page id                              |
|   fields   | string | the facebook webhook subscribed fields            |
| changeitem | string | the changes item when the facebook fields is feed |

## Run in Kubernetes

```shell
kubectl apply -f source-facebook.yaml
```

```yaml
apiVersion: v1
kind: Service
metadata:
  name: source-facebook
  namespace: vanus
spec:
  selector:
    app: source-facebook
  type: ClusterIP
  ports:
    - port: 8080
      name: source-facebook
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: source-facebook
  namespace: vanus
data:
  config.yml: |-
    target: http://<url>:<port>/gateway/<eventbus>
    verify_token: test

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: source-facebook
  namespace: vanus
  labels:
    app: source-facebook
spec:
  selector:
    matchLabels:
      app: source-facebook
  replicas: 1
  template:
    metadata:
      labels:
        app: source-facebook
    spec:
      containers:
        - name: source-facebook
          image: public.ecr.aws/vanus/connector/source-facebook:latest
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
            name: source-facebook
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

3. Update the target config of the Facebook Source

```yaml
target: http://192.168.49.2:30001/gateway/quick-start
```

4. Run the Facebook Source

```shell
kubectl apply -f source-facebook.yaml
```

[vc]: https://docs.vanus.ai/introduction/concepts#vanus-connect
[fb-webhook]: https://developers.facebook.com/docs/graph-api/webhooks/getting-started?locale=en_US#configure-webhooks-product
