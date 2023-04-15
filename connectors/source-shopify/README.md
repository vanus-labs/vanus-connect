---
title: Shopify
---
# Shopify Source

## Introduction

The Shopify Source is a [Vanus Connector][vc] which aims to convert an incoming Shopify Event Webhook Request to a CloudEvent.

For example, if type of incoming events is a `orders/create`:

which is converted to:
<details><summary><strong>Click to view</strong></summary>

```shell
+-----+-----------------------------------------------------------------------------------------------------------------------------------------------------------------------+
|     | Context Attributes,                                                                                                                                                   |
|     |   specversion: 1.0                                                                                                                                                    |
|     |   type: orders/create                                                                                                                                                 |
|     |   source: vanus-shopify-source                                                                                                                                        |
|     |   id: 20c54dff-409b-4350-ad4e-21184799800c                                                                                                                            |
|     |   time: 2023-03-17T08:34:07.358217Z                                                                                                                                   |
|     |   datacontenttype: application/json                                                                                                                                   |
|     | Extensions,                                                                                                                                                           |
|     |   xvanuseventbus: quick-start                                                                                                                                         |
|     |   xvanuslogoffset: AAAAAAAAAAE=                                                                                                                                       |
|     |   xvanusstime: 2023-03-17T08:34:07.383Z                                                                                                                               |
|     |   xvshopifyapiversion: 2023-01                                                                                                                                        |
|     |   xvshopifydomain: vanuscloudtest.myshopify.com                                                                                                                       |
|     |   xvshopifyorderid: 5277002105126                                                                                                                                     |
|     |   xvshopifytopic: orders/create                                                                                                                                       |
|     |   xvshopifywebhookid: 173b9052-6ceb-47b7-80d4-6ab4704beaa6                                                                                                            |
|     | Data,                                                                                                                                                                 |
|     |   {                                                                                                                                                                   |
|     |     "admin_graphql_api_id": "gid://shopify/Order/5277002105126",                                                                                                      |
|     |     "app_id": 1354745,                                                                                                                                                |
|     |     "browser_ip": "13.231.251.96",                                                                                                                                    |
|     |     "buyer_accepts_marketing": false,                                                                                                                                 |
|     |     "cancel_reason": null,                                                                                                                                            |
|     |     "cancelled_at": null,                                                                                                                                             |
|     |     "cart_token": null,                                                                                                                                               |
|     |     "checkout_id": 36657356603686,                                                                                                                                    |
|     |     "checkout_token": "23d74e275a0fdbb289bd4a377befa332",                                                                                                             |
|     |     "client_details": {                                                                                                                                               |
|     |       "accept_language": null,                                                                                                                                        |
|     |       "browser_height": null,                                                                                                                                         |
|     |       "browser_ip": "13.231.251.96",                                                                                                                                  |
|     |       "browser_width": null,                                                                                                                                          |
|     |       "session_hash": null,                                                                                                                                           |
|     |       "user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36"                           |
|     |     },                                                                                                                                                                |
|     |     "closed_at": null,                                                                                                                                                |
|     |     "company": null,                                                                                                                                                  |
|     |     "confirmed": true,                                                                                                                                                |
|     |     "created_at": "2023-03-17T04:34:03-04:00",                                                                                                                        |
|     |     "currency": "HKD",                                                                                                                                                |
|     |     "current_subtotal_price": "1025.00",                                                                                                                              |
|     |     "current_subtotal_price_set": {                                                                                                                                   |
|     |       "presentment_money": {                                                                                                                                          |
|     |         "amount": "1025.00",                                                                                                                                          |
|     |         "currency_code": "HKD"                                                                                                                                        |
|     |       },                                                                                                                                                              |
|     |       "shop_money": {                                                                                                                                                 |
|     |         "amount": "1025.00",                                                                                                                                          |
|     |         "currency_code": "HKD"                                                                                                                                        |
|     |       }                                                                                                                                                               |
|     |     },                                                                                                                                                                |
|     |     "current_total_discounts": "0.00",                                                                                                                                |
|     |     "current_total_discounts_set": {                                                                                                                                  |
|     |       "presentment_money": {                                                                                                                                          |
|     |         "amount": "0.00",                                                                                                                                             |
|     |         "currency_code": "HKD"                                                                                                                                        |
|     |       },                                                                                                                                                              |
|     |       "shop_money": {                                                                                                                                                 |
|     |         "amount": "0.00",                                                                                                                                             |
|     |         "currency_code": "HKD"                                                                                                                                        |
|     |       }                                                                                                                                                               |
|     |     },                                                                                                                                                                |
|     |     "current_total_duties_set": null,                                                                                                                                 |
|     |     "current_total_price": "1025.00",                                                                                                                                 |
|     |     "current_total_price_set": {                                                                                                                                      |
|     |       "presentment_money": {                                                                                                                                          |
|     |         "amount": "1025.00",                                                                                                                                          |
|     |         "currency_code": "HKD"                                                                                                                                        |
|     |       },                                                                                                                                                              |
|     |       "shop_money": {                                                                                                                                                 |
|     |         "amount": "1025.00",                                                                                                                                          |
|     |         "currency_code": "HKD"                                                                                                                                        |
|     |       }                                                                                                                                                               |
|     |     },                                                                                                                                                                |
|     |     "current_total_tax": "0.00",                                                                                                                                      |
|     |     "current_total_tax_set": {                                                                                                                                        |
|     |       "presentment_money": {                                                                                                                                          |
|     |         "amount": "0.00",                                                                                                                                             |
|     |         "currency_code": "HKD"                                                                                                                                        |
|     |       },                                                                                                                                                              |
|     |       "shop_money": {                                                                                                                                                 |
|     |         "amount": "0.00",                                                                                                                                             |
|     |         "currency_code": "HKD"                                                                                                                                        |
|     |       }                                                                                                                                                               |
|     |     },                                                                                                                                                                |
|     |     "customer_locale": "zh-CN",                                                                                                                                       |
|     |     "device_id": null,                                                                                                                                                |
|     |     "discount_applications": [],                                                                                                                                      |
|     |     "discount_codes": [],                                                                                                                                             |
|     |     "estimated_taxes": false,                                                                                                                                         |
|     |     "financial_status": "paid",                                                                                                                                       |
|     |     "fulfillment_status": null,                                                                                                                                       |
|     |     "fulfillments": [],                                                                                                                                               |
|     |     "gateway": "manual",                                                                                                                                              |
|     |     "id": 5277002105126,                                                                                                                                              |
|     |     "landing_site": null,                                                                                                                                             |
|     |     "landing_site_ref": null,                                                                                                                                         |
|     |     "line_items": [                                                                                                                                                   |
|     |       {                                                                                                                                                               |
|     |         "admin_graphql_api_id": "gid://shopify/LineItem/13782071116070",                                                                                              |
|     |         "discount_allocations": [],                                                                                                                                   |
|     |         "duties": [],                                                                                                                                                 |
|     |         "fulfillable_quantity": 1,                                                                                                                                    |
|     |         "fulfillment_service": "manual",                                                                                                                              |
|     |         "fulfillment_status": null,                                                                                                                                   |
|     |         "gift_card": false,                                                                                                                                           |
|     |         "grams": 0,                                                                                                                                                   |
|     |         "id": 13782071116070,                                                                                                                                         |
|     |         "name": "The Collection Snowboard: Oxygen",                                                                                                                   |
|     |         "price": "1025.00",                                                                                                                                           |
|     |         "price_set": {                                                                                                                                                |
|     |           "presentment_money": {                                                                                                                                      |
|     |             "amount": "1025.00",                                                                                                                                      |
|     |             "currency_code": "HKD"                                                                                                                                    |
|     |           },                                                                                                                                                          |
|     |           "shop_money": {                                                                                                                                             |
|     |             "amount": "1025.00",                                                                                                                                      |
|     |             "currency_code": "HKD"                                                                                                                                    |
|     |           }                                                                                                                                                           |
|     |         },                                                                                                                                                            |
|     |         "product_exists": true,                                                                                                                                       |
|     |         "product_id": 8187823194406,                                                                                                                                  |
|  1  |         "properties": [],                                                                                                                                             |
|     |         "quantity": 1,                                                                                                                                                |
|     |         "requires_shipping": true,                                                                                                                                    |
|     |         "sku": "",                                                                                                                                                    |
|     |         "tax_lines": [],                                                                                                                                              |
|     |         "taxable": true,                                                                                                                                              |
|     |         "title": "The Collection Snowboard: Oxygen",                                                                                                                  |
|     |         "total_discount": "0.00",                                                                                                                                     |
|     |         "total_discount_set": {                                                                                                                                       |
|     |           "presentment_money": {                                                                                                                                      |
|     |             "amount": "0.00",                                                                                                                                         |
|     |             "currency_code": "HKD"                                                                                                                                    |
|     |           },                                                                                                                                                          |
|     |           "shop_money": {                                                                                                                                             |
|     |             "amount": "0.00",                                                                                                                                         |
|     |             "currency_code": "HKD"                                                                                                                                    |
|     |           }                                                                                                                                                           |
|     |         },                                                                                                                                                            |
|     |         "variant_id": 44650297458982,                                                                                                                                 |
|     |         "variant_inventory_management": "shopify",                                                                                                                    |
|     |         "variant_title": null,                                                                                                                                        |
|     |         "vendor": "Hydrogen Vendor"                                                                                                                                   |
|     |       }                                                                                                                                                               |
|     |     ],                                                                                                                                                                |
|     |     "location_id": 79864037670,                                                                                                                                       |
|     |     "merchant_of_record_app_id": null,                                                                                                                                |
|     |     "name": "#1016",                                                                                                                                                  |
|     |     "note": null,                                                                                                                                                     |
|     |     "note_attributes": [],                                                                                                                                            |
|     |     "number": 16,                                                                                                                                                     |
|     |     "order_number": 1016,                                                                                                                                             |
|     |     "order_status_url": "https://vanuscloudtest.myshopify.com/73188213030/orders/4d73ec436313d5ff9feb8d3877686734/authenticate?key=5ccbcfe5a5743daa3c65fc6e1bb97b54", |
|     |     "original_total_duties_set": null,                                                                                                                                |
|     |     "payment_gateway_names": [                                                                                                                                        |
|     |       "manual"                                                                                                                                                        |
|     |     ],                                                                                                                                                                |
|     |     "payment_terms": null,                                                                                                                                            |
|     |     "presentment_currency": "HKD",                                                                                                                                    |
|     |     "processed_at": "2023-03-17T04:34:03-04:00",                                                                                                                      |
|     |     "processing_method": "manual",                                                                                                                                    |
|     |     "reference": "c3965dfb4e5cd89674942afe36d8e116",                                                                                                                  |
|     |     "referring_site": null,                                                                                                                                           |
|     |     "refunds": [],                                                                                                                                                    |
|     |     "shipping_lines": [],                                                                                                                                             |
|     |     "source_identifier": "c3965dfb4e5cd89674942afe36d8e116",                                                                                                          |
|     |     "source_name": "shopify_draft_order",                                                                                                                             |
|     |     "source_url": null,                                                                                                                                               |
|     |     "subtotal_price": "1025.00",                                                                                                                                      |
|     |     "subtotal_price_set": {                                                                                                                                           |
|     |       "presentment_money": {                                                                                                                                          |
|     |         "amount": "1025.00",                                                                                                                                          |
|     |         "currency_code": "HKD"                                                                                                                                        |
|     |       },                                                                                                                                                              |
|     |       "shop_money": {                                                                                                                                                 |
|     |         "amount": "1025.00",                                                                                                                                          |
|     |         "currency_code": "HKD"                                                                                                                                        |
|     |       }                                                                                                                                                               |
|     |     },                                                                                                                                                                |
|     |     "tags": "",                                                                                                                                                       |
|     |     "tax_lines": [],                                                                                                                                                  |
|     |     "taxes_included": false,                                                                                                                                          |
|     |     "test": false,                                                                                                                                                    |
|     |     "token": "4d73ec436313d5ff9feb8d3877686734",                                                                                                                      |
|     |     "total_discounts": "0.00",                                                                                                                                        |
|     |     "total_discounts_set": {                                                                                                                                          |
|     |       "presentment_money": {                                                                                                                                          |
|     |         "amount": "0.00",                                                                                                                                             |
|     |         "currency_code": "HKD"                                                                                                                                        |
|     |       },                                                                                                                                                              |
|     |       "shop_money": {                                                                                                                                                 |
|     |         "amount": "0.00",                                                                                                                                             |
|     |         "currency_code": "HKD"                                                                                                                                        |
|     |       }                                                                                                                                                               |
|     |     },                                                                                                                                                                |
|     |     "total_line_items_price": "1025.00",                                                                                                                              |
|     |     "total_line_items_price_set": {                                                                                                                                   |
|     |       "presentment_money": {                                                                                                                                          |
|     |         "amount": "1025.00",                                                                                                                                          |
|     |         "currency_code": "HKD"                                                                                                                                        |
|     |       },                                                                                                                                                              |
|     |       "shop_money": {                                                                                                                                                 |
|     |         "amount": "1025.00",                                                                                                                                          |
|     |         "currency_code": "HKD"                                                                                                                                        |
|     |       }                                                                                                                                                               |
|     |     },                                                                                                                                                                |
|     |     "total_outstanding": "0.00",                                                                                                                                      |
|     |     "total_price": "1025.00",                                                                                                                                         |
|     |     "total_price_set": {                                                                                                                                              |
|     |       "presentment_money": {                                                                                                                                          |
|     |         "amount": "1025.00",                                                                                                                                          |
|     |         "currency_code": "HKD"                                                                                                                                        |
|     |       },                                                                                                                                                              |
|     |       "shop_money": {                                                                                                                                                 |
|     |         "amount": "1025.00",                                                                                                                                          |
|     |         "currency_code": "HKD"                                                                                                                                        |
|     |       }                                                                                                                                                               |
|     |     },                                                                                                                                                                |
|     |     "total_shipping_price_set": {                                                                                                                                     |
|     |       "presentment_money": {                                                                                                                                          |
|     |         "amount": "0.00",                                                                                                                                             |
|     |         "currency_code": "HKD"                                                                                                                                        |
|     |       },                                                                                                                                                              |
|     |       "shop_money": {                                                                                                                                                 |
|     |         "amount": "0.00",                                                                                                                                             |
|     |         "currency_code": "HKD"                                                                                                                                        |
|     |       }                                                                                                                                                               |
|     |     },                                                                                                                                                                |
|     |     "total_tax": "0.00",                                                                                                                                              |
|     |     "total_tax_set": {                                                                                                                                                |
|     |       "presentment_money": {                                                                                                                                          |
|     |         "amount": "0.00",                                                                                                                                             |
|     |         "currency_code": "HKD"                                                                                                                                        |
|     |       },                                                                                                                                                              |
|     |       "shop_money": {                                                                                                                                                 |
|     |         "amount": "0.00",                                                                                                                                             |
|     |         "currency_code": "HKD"                                                                                                                                        |
|     |       }                                                                                                                                                               |
|     |     },                                                                                                                                                                |
|     |     "total_tip_received": "0.00",                                                                                                                                     |
|     |     "total_weight": 0,                                                                                                                                                |
|     |     "updated_at": "2023-03-17T04:34:04-04:00",                                                                                                                        |
|     |     "user_id": 93721493798                                                                                                                                            |
|     |   }                                                                                                                                                                   |
|     |                                                                                                                                                                       |
+-----+-----------------------------------------------------------------------------------------------------------------------------------------------------------------------+
```
</details>

## Quick Start

This section will show you how to use Shopify Source to convert a Shopify Order Create Webhook request to a CloudEvent.

### Prerequisites

- Have `Docker`
- Have `cURL`
- Have acknowledges on how to create Shopify Webhook. details can be found [here](https://shopify.dev/docs/apps/webhooks/configuration/shopifys)
- Have [Ngrok](https://ngrok.com/), this makes Shopify Source available on Internet

### Create Config file

```shell
cat << EOF > config.yml
target: http://localhost:31081
client_secret: <client_secret_of_your_app>
EOF
```


| Name          | Required | Default | Description                        |
| :-------------- | :--------: | :-------: | :----------------------------------- |
| target        |   YES   |        | the target URL to send CloudEvents |
| client_secret |   YES   |        | the client secret of your app      |

The Shopify Source tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

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
  "id": "aa5a3aa1-fbd2-4d76-961c-cc205f5b625b",
  "source": "vanus-shopify-source",
  "type": "orders/create",
  "datacontenttype": "application/json",
  "time": "2023-03-17T11:14:55.307657525Z",
  "data": {
    "admin_graphql_api_id": "gid://shopify/Order/5277067936038",
    "app_id": 1354745,
    "browser_ip": "3.34.96.234",
    "buyer_accepts_marketing": false,
    "cancel_reason": null,
    "cancelled_at": null,
    "cart_token": null,
    "checkout_id": 36657551147302,
    "checkout_token": "5e58edbdaa96ad5e2dc2d6d9a95f4b68",
    "client_details": {
      "accept_language": null,
      "browser_height": null,
      "browser_ip": "3.34.96.234",
      "browser_width": null,
      "session_hash": null,
      "user_agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36"
    },
    "closed_at": null,
    "company": null,
    "confirmed": true,
    "created_at": "2023-03-17T07:14:52-04:00",
    "currency": "HKD",
    "current_subtotal_price": "2629.95",
    "current_subtotal_price_set": {
      "presentment_money": {
        "amount": "2629.95",
        "currency_code": "HKD"
      },
      "shop_money": {
        "amount": "2629.95",
        "currency_code": "HKD"
      }
    },
    "current_total_discounts": "0.00",
    "current_total_discounts_set": {
      "presentment_money": {
        "amount": "0.00",
        "currency_code": "HKD"
      },
      "shop_money": {
        "amount": "0.00",
        "currency_code": "HKD"
      }
    },
    "current_total_duties_set": null,
    "current_total_price": "2629.95",
    "current_total_price_set": {
      "presentment_money": {
        "amount": "2629.95",
        "currency_code": "HKD"
      },
      "shop_money": {
        "amount": "2629.95",
        "currency_code": "HKD"
      }
    },
    "current_total_tax": "0.00",
    "current_total_tax_set": {
      "presentment_money": {
        "amount": "0.00",
        "currency_code": "HKD"
      },
      "shop_money": {
        "amount": "0.00",
        "currency_code": "HKD"
      }
    },
    "customer_locale": "en",
    "device_id": null,
    "discount_applications": [],
    "discount_codes": [],
    "estimated_taxes": false,
    "financial_status": "paid",
    "fulfillment_status": null,
    "fulfillments": [],
    "gateway": "manual",
    "id": 5277067936038,
    "landing_site": null,
    "landing_site_ref": null,
    "line_items": [
      {
        "admin_graphql_api_id": "gid://shopify/LineItem/13782232269094",
        "discount_allocations": [],
        "duties": [],
        "fulfillable_quantity": 1,
        "fulfillment_service": "snow-city-warehouse",
        "fulfillment_status": null,
        "gift_card": false,
        "grams": 0,
        "id": 13782232269094,
        "name": "The 3p Fulfilled Snowboard",
        "price": "2629.95",
        "price_set": {
          "presentment_money": {
            "amount": "2629.95",
            "currency_code": "HKD"
          },
          "shop_money": {
            "amount": "2629.95",
            "currency_code": "HKD"
          }
        },
        "product_exists": true,
        "product_id": 8187823063334,
        "properties": [],
        "quantity": 1,
        "requires_shipping": true,
        "sku": "sku-hosted-1",
        "tax_lines": [],
        "taxable": true,
        "title": "The 3p Fulfilled Snowboard",
        "total_discount": "0.00",
        "total_discount_set": {
          "presentment_money": {
            "amount": "0.00",
            "currency_code": "HKD"
          },
          "shop_money": {
            "amount": "0.00",
            "currency_code": "HKD"
          }
        },
        "variant_id": 44650297164070,
        "variant_inventory_management": "shopify",
        "variant_title": null,
        "vendor": "VanusCloudTest"
      }
    ],
    "location_id": 79864037670,
    "merchant_of_record_app_id": null,
    "name": "#1018",
    "note": null,
    "note_attributes": [],
    "number": 18,
    "order_number": 1018,
    "order_status_url": "https://vanuscloudtest.myshopify.com/73188213030/orders/67c021532f1acd9731573def31c96678/authenticate?key=ae755b290a6285ecf85dcb4a406e7f9d",
    "original_total_duties_set": null,
    "payment_gateway_names": [
      "manual"
    ],
    "payment_terms": null,
    "presentment_currency": "HKD",
    "processed_at": "2023-03-17T07:14:51-04:00",
    "processing_method": "manual",
    "reference": "d1487c532906960e20e32149cad2cb8d",
    "referring_site": null,
    "refunds": [],
    "shipping_lines": [],
    "source_identifier": "d1487c532906960e20e32149cad2cb8d",
    "source_name": "shopify_draft_order",
    "source_url": null,
    "subtotal_price": "2629.95",
    "subtotal_price_set": {
      "presentment_money": {
        "amount": "2629.95",
        "currency_code": "HKD"
      },
      "shop_money": {
        "amount": "2629.95",
        "currency_code": "HKD"
      }
    },
    "tags": "",
    "tax_lines": [],
    "taxes_included": false,
    "test": false,
    "token": "67c021532f1acd9731573def31c96678",
    "total_discounts": "0.00",
    "total_discounts_set": {
      "presentment_money": {
        "amount": "0.00",
        "currency_code": "HKD"
      },
      "shop_money": {
        "amount": "0.00",
        "currency_code": "HKD"
      }
    },
    "total_line_items_price": "2629.95",
    "total_line_items_price_set": {
      "presentment_money": {
        "amount": "2629.95",
        "currency_code": "HKD"
      },
      "shop_money": {
        "amount": "2629.95",
        "currency_code": "HKD"
      }
    },
    "total_outstanding": "0.00",
    "total_price": "2629.95",
    "total_price_set": {
      "presentment_money": {
        "amount": "2629.95",
        "currency_code": "HKD"
      },
      "shop_money": {
        "amount": "2629.95",
        "currency_code": "HKD"
      }
    },
    "total_shipping_price_set": {
      "presentment_money": {
        "amount": "0.00",
        "currency_code": "HKD"
      },
      "shop_money": {
        "amount": "0.00",
        "currency_code": "HKD"
      }
    },
    "total_tax": "0.00",
    "total_tax_set": {
      "presentment_money": {
        "amount": "0.00",
        "currency_code": "HKD"
      },
      "shop_money": {
        "amount": "0.00",
        "currency_code": "HKD"
      }
    },
    "total_tip_received": "0.00",
    "total_weight": 0,
    "updated_at": "2023-03-17T07:14:53-04:00",
    "user_id": 93721493798
  },
  "xvshopifyapiversion": "2023-01",
  "xvshopifytopic": "orders/create",
  "xvshopifyorderid": "5277067936038",
  "xvshopifydomain": "vanuscloudtest.myshopify.com",
  "xvshopifywebhookid": "ad6a6b00-2352-4dfc-9e57-aabe1d12f2da"
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

The Shopify Source defines following [CloudEvents Extension Attributes](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#extension-context-attributes)


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

This section shows how a source connector can send CloudEvents to a running [Vanus cluster](https://github.com/vanus-labs/vanus).

### Prerequisites

- Have a running K8s cluster
- Have a running Vanus cluster
- Vsctl Installed

1. Export the VANUS_GATEWAY environment variable (the ip should be a host-accessible address of the vanus-gateway service)

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
