---
title: Shopify App
---

# Shopify App Source

## Introduction

The Shopify App Source is a [Vanus Connector][vc] which aims to convert shopify data
to a CloudEvent. The Shopify App Source uses [Admin API][admin-rest] and pulls
data at a fixed time.

The shopify data is converted to:

```json
{
  "specversion": "1.0",
  "id": "026046e2-3cb0-4116-895e-c77877072dd2",
  "source": "shopfiy-source-shop-name",
  "datacontenttype": "application/json",
  "time": "2023-01-28T06:11:10.012579049Z",
  "data": {
  }
}
```

## Quick Start

This section shows how Shopify App Source converts data to a CloudEvent.

### Create the config file

```shell
cat << EOF > config.yml
target: http://localhost:31081
shop_name: shop_name
api_access_token: shpat_xxxxxx
sync_begin_date: 2023-04-01
EOF
```

| Name             | Required | Default | Description                        |
|:-----------------|:---------|:--------|:-----------------------------------|
| target           | YES      |         | the target URL to send CloudEvents |
| shop_name        | YES      |         | the shop name                      |
| api_access_token | YES      |         | the admin api access token         |
| sync_begin_date  | YES      |         | sync begin date                    |

The Shopify App Source tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify
the position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-shopify-app public.ecr.aws/vanus/connector/source-shopify-app
```

### Test

Before starting Shopify App Source use the following command to run a Display sink, which will receive and prints the
incoming CloudEvents.

```shell
docker run -it --rm \
  -p 31081:8080 \
  --name sink-display public.ecr.aws/vanus/connector/sink-display
```

Make sure the `target` value in your config file is `http://localhost:31081` so that the Source can send CloudEvents to
our Display Sink.

Now run Shopify App source to send CloudEvents

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-shopify-app public.ecr.aws/vanus/connector/source-shopify-app
```

Here is the sort of CloudEvent you should expect to receive in the Display Sink:

```json
{
  "specversion": "1.0",
  "id": "026046e2-3cb0-4116-895e-c77877072dd2",
  "source": "shopfiy-source-shop-name",
  "datacontenttype": "application/json",
  "time": "2023-01-28T06:11:10.012579049Z",
  "data": {
  }
}
```

To see the CloudEvents on your Terminal, run:

```shell
docker logs sink-display
```

### Clean

```shell
docker stop source-shopify-app sink-display
```

## Run in Kubernetes

```shell
kubectl apply -f source-shopify-app.yaml
```

```yaml
apiVersion: v1
kind: Service
metadata:
 name: source-shopify-app
 namespace: vanus
spec:
 selector:
  app: source-shopify-app
 type: ClusterIP
 ports:
  - port: 8080
    name: source-shopify-app
---
apiVersion: v1
kind: ConfigMap
metadata:
 name: source-shopify-app
 namespace: vanus
data:
 config.yml: |-
  target: http://<url>:<port>/gateway/<eventbus>
  shop_name: shop_name
  api_access_token: shpat_xxxxxx
  sync_begin_date: 2023-04-01
---
apiVersion: apps/v1
kind: Deployment
metadata:
 name: source-shopify-app
 namespace: vanus
 labels:
  app: source-shopify-app
spec:
 selector:
  matchLabels:
   app: source-shopify-app
 replicas: 1
 template:
  metadata:
   labels:
    app: source-shopify-app
  spec:
   containers:
    - name: source-shopify-app
      image: public.ecr.aws/vanus/connector/source-shopify-app:latest
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
       name: source-shopify-app
```

[vc]: https://docs.vanus.ai/introduction/concepts#vanus-connect
[admin-rest]: https://shopify.dev/api/admin-rest
