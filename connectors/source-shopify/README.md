---
title: Shopify
---

# Shopify Source

## Introduction

The Shopify Source is a [Vanus Connector][vc] which aims to convert an incoming Shopify Event Webhook Request to a
CloudEvent.

For example, if type of incoming events is a `orders/create`:

which is converted to:
<details><summary><strong>Click to view</strong></summary>

```json
{
  "specversion": "1.0",
  "id": "bf6672f8-40f5-4b87-91c3-5eccbf4aec89",
  "source": "vanus-shopify-source",
  "type": "products/update",
  "datacontenttype": "application/json",
  "time": "2023-10-23T09:24:40.720817Z",
  "data": {
    "admin_graphql_api_id": "gid://shopify/Product/788032119674292922",
    "body_html": "An example T-Shirt",
    "created_at": null,
    "handle": "example-t-shirt",
    "id": 788032119674292900,
    "image": null,
    "images": [],
    "options": [],
    "product_type": "Shirts",
    "published_at": "2023-10-23T05:24:40-04:00",
    "published_scope": "web",
    "status": "active",
    "tags": "example, mens, t-shirt",
    "template_suffix": null,
    "title": "Example T-Shirt",
    "updated_at": "2023-10-23T05:24:40-04:00",
    "variants": [
      {
        "admin_graphql_api_id": "gid://shopify/ProductVariant/642667041472713922",
        "barcode": null,
        "compare_at_price": "24.99",
        "created_at": null,
        "fulfillment_service": "manual",
        "grams": 200,
        "id": 642667041472714000,
        "image_id": null,
        "inventory_item_id": null,
        "inventory_management": "shopify",
        "inventory_policy": "deny",
        "inventory_quantity": 75,
        "old_inventory_quantity": 75,
        "option1": "Small",
        "option2": null,
        "option3": null,
        "position": 0,
        "price": "19.99",
        "product_id": 788032119674292900,
        "requires_shipping": true,
        "sku": "example-shirt-s",
        "taxable": true,
        "title": "",
        "updated_at": null,
        "weight": 200,
        "weight_unit": "g"
      },
      {
        "admin_graphql_api_id": "gid://shopify/ProductVariant/757650484644203962",
        "barcode": null,
        "compare_at_price": "24.99",
        "created_at": null,
        "fulfillment_service": "manual",
        "grams": 200,
        "id": 757650484644203900,
        "image_id": null,
        "inventory_item_id": null,
        "inventory_management": "shopify",
        "inventory_policy": "deny",
        "inventory_quantity": 50,
        "old_inventory_quantity": 50,
        "option1": "Medium",
        "option2": null,
        "option3": null,
        "position": 0,
        "price": "19.99",
        "product_id": 788032119674292900,
        "requires_shipping": true,
        "sku": "example-shirt-m",
        "taxable": true,
        "title": "",
        "updated_at": null,
        "weight": 200,
        "weight_unit": "g"
      }
    ],
    "vendor": "Acme"
  },
  "xvshopifywebhookid": "614966e9-2907-476e-bce5-65f6fc5c1f9e",
  "xvshopifydomain": "vanusfashion.myshopify.com",
  "xvshopifyapiversion": "2023-10",
  "xvshopifytriggeredat": "2023-10-23T09:24:40.066647883Z"
}
```

</details>

## Quick Start

This section will show you how to use Shopify Source to convert a Shopify Order Create Webhook request to a CloudEvent.

### Prerequisites

- Have `Docker`
- Have `curl`
- Have acknowledges on how to create Shopify Webhook. details can be
  found [here](https://shopify.dev/docs/apps/webhooks/configuration/shopifys)
- Have [Ngrok](https://ngrok.com/), this makes Shopify Source available on Internet

### Create Config file

```shell
cat << EOF > config.yml
target: http://localhost:31081
client_secret: <client_secret_of_your_app>
EOF
```

|      Name      | Required  | Default  |            Description             |
|:--------------:|:---------:|:--------:|:----------------------------------:|
|     target     |    YES    |          | the target URL to send CloudEvents |
| client_secret  |    YES    |          |   the client secret of your app    |

The Shopify Source tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the
position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-shopify public.ecr.aws/vanus/connector/source-shopify
```

### Run Ngrok to expose Shopify Source to internet

```shell
ngrok http 8080
```

you will get a `Forwarding` URL like `https://xxxx.xxxx.ngrok.io` after ngrok started

![ngrok.png](ngrok.png)

### Create a Shopify webhook via cURL

replace `<your_shop_name>, <your_shop_access_token>, <Forwarding URL>` to yours.

```shell
curl --location --request POST 'https://<your_shop_name>.myshopify.com/admin/api/2023-01/webhooks.json' \
--header 'X-Shopify-Access-Token: <your_shop_access_token>' \
--header 'Content-Type: application/json' \
--data-raw '{
    "webhook": {
        "address": "<Forwarding URL>",
        "topic": "orders/create",
        "format": "json"
    }
}'
```

### Test

Open a terminal and use the following command to run a Display sink, which receives and prints CloudEvents.

```shell
docker run -it --rm \
  -p 31081:8080 \
  --name sink-display public.ecr.aws/vanus/connector/sink-display
```

Open the browser and create a test order in your shop,

![shopify.png](shopify.png)

Here is the sort of CloudEvent you should expect to receive in the Display Sink:

<details><summary><strong>Click to view</strong></summary>

```json
{
  "specversion": "1.0",
  "id": "bf6672f8-40f5-4b87-91c3-5eccbf4aec89",
  "source": "vanus-shopify-source",
  "type": "products/update",
  "datacontenttype": "application/json",
  "time": "2023-10-23T09:24:40.720817Z",
  "data": {
    "admin_graphql_api_id": "gid://shopify/Product/788032119674292922",
    "body_html": "An example T-Shirt",
    "created_at": null,
    "handle": "example-t-shirt",
    "id": 788032119674292900,
    "image": null,
    "images": [],
    "options": [],
    "product_type": "Shirts",
    "published_at": "2023-10-23T05:24:40-04:00",
    "published_scope": "web",
    "status": "active",
    "tags": "example, mens, t-shirt",
    "template_suffix": null,
    "title": "Example T-Shirt",
    "updated_at": "2023-10-23T05:24:40-04:00",
    "variants": [
      {
        "admin_graphql_api_id": "gid://shopify/ProductVariant/642667041472713922",
        "barcode": null,
        "compare_at_price": "24.99",
        "created_at": null,
        "fulfillment_service": "manual",
        "grams": 200,
        "id": 642667041472714000,
        "image_id": null,
        "inventory_item_id": null,
        "inventory_management": "shopify",
        "inventory_policy": "deny",
        "inventory_quantity": 75,
        "old_inventory_quantity": 75,
        "option1": "Small",
        "option2": null,
        "option3": null,
        "position": 0,
        "price": "19.99",
        "product_id": 788032119674292900,
        "requires_shipping": true,
        "sku": "example-shirt-s",
        "taxable": true,
        "title": "",
        "updated_at": null,
        "weight": 200,
        "weight_unit": "g"
      },
      {
        "admin_graphql_api_id": "gid://shopify/ProductVariant/757650484644203962",
        "barcode": null,
        "compare_at_price": "24.99",
        "created_at": null,
        "fulfillment_service": "manual",
        "grams": 200,
        "id": 757650484644203900,
        "image_id": null,
        "inventory_item_id": null,
        "inventory_management": "shopify",
        "inventory_policy": "deny",
        "inventory_quantity": 50,
        "old_inventory_quantity": 50,
        "option1": "Medium",
        "option2": null,
        "option3": null,
        "position": 0,
        "price": "19.99",
        "product_id": 788032119674292900,
        "requires_shipping": true,
        "sku": "example-shirt-m",
        "taxable": true,
        "title": "",
        "updated_at": null,
        "weight": 200,
        "weight_unit": "g"
      }
    ],
    "vendor": "Acme"
  },
  "xvshopifywebhookid": "614966e9-2907-476e-bce5-65f6fc5c1f9e",
  "xvshopifydomain": "vanusfashion.myshopify.com",
  "xvshopifyapiversion": "2023-10",
  "xvshopifytriggeredat": "2023-10-23T09:24:40.066647883Z"
}
```

</details>

### Clean

```shell
docker stop source-shopify sink-display
```

## Source details

### Attributes

| Attribute |        Default        |
| :---------: | :----------------------: |
|    id    |          UUID          |
|  source  |  vanus-shopify-source  |
|   type   | {the topic of request} |

#### Extension Attributes

The Shopify Source defines
following [CloudEvents Extension Attributes](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#extension-context-attributes)

|      Attribute      |  Type  | Description                                                                                                                                      |
| :-------------------: | :------: | :------------------------------------------------------------------------------------------------------------------------------------------------- |
|  xvshopifyorderid  | string | Order id if have                                                                                                                                 |
|   xvshopifytopic   | string | The topic of incoming request, full topics can be found in[here](https://shopify.dev/docs/api/admin-rest/2023-01/resources/webhook#event-topics) |
| xvshopifywebhookid | string | The webhook id of incoming request belongs to                                                                                                    |
|   xvshopifydomain   | string | The shop name of incoming request belongs to                                                                                                     |
| xvshopifyapiversion | string | The Shopify Request API Version                                                                                                                  |

## Run in Kubernetes

```shell
kubectl apply -f source-shopify.yaml
```

```yaml
apiVersion: v1
kind: Service
metadata:
  name: source-shopify
  namespace: vanus
spec:
  selector:
    app: source-shopify
  type: ClusterIP
  ports:
    - port: 8080
      name: source-shopify
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: source-shopify
  namespace: vanus
data:
  config.yml: |-
    target: http://<url>:<port>/gateway/<eventbus>
    client_secret: "xxxxxxx"

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: source-shopify
  namespace: vanus
  labels:
    app: source-shopify
spec:
  selector:
    matchLabels:
      app: source-shopify
  replicas: 1
  template:
    metadata:
      labels:
        app: source-shopify
    spec:
      containers:
        - name: source-shopify
          image: public.ecr.aws/vanus/connector/source-shopify:latest
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
            name: source-shopify
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

3. Update the target config of the Shopify Source

```yaml
target: http://192.168.49.2:30001/gateway/quick-start
```

4. Run the Shopify Source

```shell
kubectl apply -f source-shopify.yaml
```

[vc]: https://docs.vanus.ai/introduction/concepts#vanus-connect
