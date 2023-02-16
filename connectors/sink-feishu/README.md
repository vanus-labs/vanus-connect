---
title: Feishu
---

# Feishu Sink

## Introduction

The Feishu Sink is a [Vanus Connector][vc] which aims to handle incoming CloudEvents in a way that extracts the `data` part of the
original event and deliver these extracted `data` to  the Feishu APIs. Now the Sink support Feishu Bot: pushing a message to Group Chat with message of text, post, share_chat, image, and interactive.

## Quick Start

In this section, we show how to use Feishu Sink push a text message to your group chat.

### Prerequisites

#### Add a bot to your group chat

Go to your target group, click Chat Settings > Group Bots > Add Bot, and select Custom Bot to add the bot to the group chat.

Enter a name and description for your bot, or set up an avatar for the bot, and then click "Add".

![add-a-bot](https://github.com/linkall-labs/vanus-connect/blob/main/connectors/sink-feishu/add-a-bot.gif?raw=true)

You will get the webhook address of the bot in the following format:

```
https://open.feishu.cn/open-apis/bot/v2/hook/xxxxxxxxxxxxxxxxx
```

> ⚠️ Please keep this webhook address properly. Do not publish it on GitHub, blogs, and other publicly accessible sites to avoid it being maliciously called to send spam messages.

![bot-config](https://github.com/linkall-labs/vanus-connect/blob/main/connectors/sink-feishu/feishu-config.png?raw=true)

> ⚠️ You must set your signature verification to make sure push messages work.

### Create the config file

Replace `chat_group`, `signature`, and `address` to yours. `chat_group` can be fill in any value as you want.

```shell
cat << EOF > config.yml
enable: ["bot"]
bot:
  webhooks:
    - chat_group: "bot1"
      signature: "xxxxxxx"
      url: "https://open.feishu.cn/open-apis/bot/v2/hook/xxxxxxx"
EOF
```


| Name                       | Required | Default | Description                                                                       |
|:---------------------------|:--------:|:-------:|-----------------------------------------------------------------------------------|
| enable                     | **YES**  |    -    | service list you want Feishu Sink is enabled                                      |
| bot.webhooks               | **YES**  |    -    | list of chat-group's configuration                                                |
| bot.webhooks.[].chat_group | **YES**  |    -    | the chat_group name, you can set any value to it                                  |
| bot.webhooks.[].signature  | **YES**  |    -    | the signature to sign request, you can get it when you create Chat Bot            |
| bot.webhooks.[].url        | **YES**  |    -    | the webhook address that message sent to, you can get it when you create Chat Bot |

The Feishu Sink tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm \
  -p 31080:8080 \
  -v ${PWD}:/vanus-connect/config \
  --name sink-feishu public.ecr.aws/vanus/connector/sink-feishu
```

### Test

Open a terminal and use THE following command to send a CloudEvent to the Sink.

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "sink-feishu-quickstart",
    "specversion": "1.0",
    "type": "quickstart",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "xvfeishuservice": "bot",
    "xvfeishuchatgroup": "bot1",
    "xvfeishumsgtype": "text",
    "data": "Hello Feishu"
}'
```

now, you can see a notification from your bot in your group chat.
![received-notification](https://github.com/linkall-labs/vanus-connect/blob/main/connectors/sink-feishu/received-message.png?raw=true)

### Clean

```shell
docker stop sink-feishu
```

## Sink details

### Extension Attributes

Feishu Sink has defined a few [CloudEvents Extension Attribute](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#extension-context-attributes)
to determine how to process event.

| Attribute         | Required | Examples               | Description                                                                                                                                                  |
|:------------------|:--------:|------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------|
| xvfeishuservice   | **YES**  | bot                    | which Feishu Service the event sent for                                                                                                                      |
| xvfeishumsgtype   | **YES**  | text                   | which Message Type the event convert to                                                                                                                      |
| xvfeishuchatgroup |    NO    | test_bot               | which Feishu chat-group the event sent for, the value should associate with you wrote in configuration, if `dynamic_route=false`, this attribute can't empty |
| xvfeishuboturls   |    NO    | bot1,bot2,bot3         | dynamic webhook urls, use  `,` to separate multiple urls.                                                                                                    |
| xvfeishubotsigns  |    NO    | signature1,,signature3 | dynamic webhook signatures, use  `,` to separate multiple signatures.                                                                                        |

**the number of urls represented by `xvfeishuboturls` must equal to the number of signatures represented by `xvfeishuboturls`**

### Chat Bot Dynamic Webhook

In some cases, users can't make sure how many bots there have or send one message to multiple groups, which means they need to dynamically send message to
Feishu Bot Service, `Chat Bot Dynamic Webhook` helps users do that.

in `config.yml`, set `dynamic_route=true` to enable this feature, otherwise `xvfeishuboturls` and `xvfeishubotsigns` will be ignored.

```yaml
enable: ["bot"]
bot:
  dynamic_route: true
  webhooks:
    - chat_group: "bot_predefined"
      url: "https://open.feishu.cn/open-apis/bot/v2/hook/xxxxx"
      signature: "xxxxxx"
```

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "sink-feishu-quickstart",
    "specversion": "1.0",
    "type": "quickstart",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "xvfeishuservice": "bot",
    "xvfeishumsgtype": "text",
    "xvfeishuboturls": "https://open.feishu.cn/open-apis/bot/v2/hook/xxx,https://open.feishu.cn/open-apis/bot/v2/hook/xxx,https://open.feishu.cn/open-apis/bot/v2/hook/xxx",
    "xvfeishubotsigns": "signature1,,signature3",
    "data": "Hello Feishu"
}'
```

this request means send a text message to three chat groups, and the second group hasn't signature.

Moreover, dynamic webhooks can work together with `xvfeishuchatgroup` attribute.

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "sink-feishu-quickstart",
    "specversion": "1.0",
    "type": "quickstart",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "xvfeishuservice": "bot",
    "xvfeishumsgtype": "text",
    "xvfeishuchatgroup": "bot_predefined",
    "xvfeishuboturls": "https://open.feishu.cn/open-apis/bot/v2/hook/xxx,https://open.feishu.cn/open-apis/bot/v2/hook/xxx,https://open.feishu.cn/open-apis/bot/v2/hook/xxx",
    "xvfeishubotsigns": "signature1,,signature3",
    "data": "Hello Feishu"
}'
```

this request will also send message to `bot_predefined`.

Note: Specified chat group was represented by `xvfeishuchatgroup` will be ignored if it wasn't be found in configuration.

## Examples

### Feishu Bot

you could find official docs of Feishu bot in https://open.feishu.cn/document/ukTMukTMukTM/ucTM5YjL3ETO24yNxkjN?lang=zh-CN#132a114c.

#### Text Message

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "sink-feishu-quickstart",
    "specversion": "1.0",
    "type": "quickstart",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "xvfeishuservice": "bot",
    "xvfeishuchatgroup": "bot1",
    "xvfeishumsgtype": "text",
    "data": "Hello Feishu"
}'
```

#### Post Message

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "sink-feishu-quickstart",
    "specversion": "1.0",
    "type": "quickstart",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "xvfeishuservice": "bot",
    "xvfeishuchatgroup": "bot1",
    "xvfeishumsgtype": "post",
    "data": {
        "zh_cn": {
            "title": "项目更新通知",
            "content": [
                [{
                        "tag": "text",
                        "text": "项目有更新: "
                    },
                    {
                        "tag": "a",
                        "text": "请查看",
                        "href": "http://www.baidu.com/"
                    },
                    {
                        "tag": "at",
                        "user_id": "abcdefgh"
                    }
                ]
            ]
        }
    }
}'
```

#### ShareChat Message

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "sink-feishu-quickstart",
    "specversion": "1.0",
    "type": "quickstart",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "xvfeishuservice": "bot",
    "xvfeishuchatgroup": "bot1",
    "xvfeishumsgtype": "share_chat",
    "data": "oc_ad6c99f9"
}'
```

#### Image Message

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "sink-feishu-quickstart",
    "specversion": "1.0",
    "type": "quickstart",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "xvfeishuservice": "bot",
    "xvfeishuchatgroup": "bot1",
    "xvfeishumsgtype": "image",
    "data": {
        "target": "feishu"
    }
}'
```

#### Interactive Message

```shell
curl --location --request POST 'localhost:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
    "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
    "source": "sink-feishu-quickstart",
    "specversion": "1.0",
    "type": "quickstart",
    "datacontenttype": "application/json",
    "time": "2022-10-26T10:38:29.345Z",
    "xvfeishuservice": "bot",
    "xvfeishuchatgroup": "bot1",
    "xvfeishumsgtype": "interactive",
    "data": {
        "elements": [
            {
                "tag": "div",
                "text": {
                    "content": "**西湖**，位于浙江省杭州市西湖区龙井路1号，杭州市区西部，景区总面积49平方千米，汇水面积为21.22平方千米，湖面面积为6.38平方千米。",
                    "tag": "lark_md"
                }
            },
            {
                "actions": [
                    {
                        "tag": "button",
                        "text": {
                            "content": "更多景点介绍 :玫瑰:",
                            "tag": "lark_md"
                        },
                        "url": "https://www.example.com",
                        "type": "default",
                        "value": {}
                    }
                ],
                "tag": "action"
            }
        ],
        "header": {
            "title": {
                "content": "今日旅游推荐",
                "tag": "plain_text"
            }
        }
    }
}'
```

## Run in Kubernetes

```shell
kubectl apply -f sink-feishu.yaml
```

```yaml
apiVersion: v1
kind: Service
metadata:
  name: sink-feishu
  namespace: vanus
spec:
  selector:
    app: sink-feishu
  type: ClusterIP
  ports:
    - port: 8080
      name: sink-feishu
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: sink-feishu
  namespace: vanus
data:
  config.yml: |-
    enable: ["bot"]
    bot:
      dynamic_route: false
      webhooks:
        - chat_group: "bot1"
          url: "https://open.feishu.cn/open-apis/bot/v2/hook/xxxxxxxxxx"
          signature: "xxxxxxxxxx"
        - chat_group: "bot2"
          signature: "xxxxxxxxxx"
          url: "https://open.feishu.cn/open-apis/bot/v2/hook/xxxxxxxxxx"

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sink-feishu
  namespace: vanus
  labels:
    app: sink-feishu
spec:
  selector:
    matchLabels:
      app: sink-feishu
  replicas: 1
  template:
    metadata:
      labels:
        app: sink-feishu
    spec:
      containers:
        - name: sink-feishu
          image: public.ecr.aws/vanus/connector/sink-feishu:latest
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "512Mi"
              cpu: "500m"
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
            name: sink-feishu
```

## Integrate with Vanus

This section shows how a sink connector can receive CloudEvents from a running [Vanus cluster](https://github.com/linkall-labs/vanus).

1. Run the sink-feishu.yaml
```shell
kubectl apply -f sink-feishu.yaml
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
  --sink 'http://sink-feishu:8080'
```

[vc]: https://www.vanus.ai/introduction/concepts#vanus-connect
