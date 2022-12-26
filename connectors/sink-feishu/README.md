---
title: Feishu
---

# Feishu Sink

## Introduction

The Feishu Sink is a [Vance Connector](vc) which aims to handle incoming CloudEvents in a way that extracts the `data` part of the
original event and deliver these extracted `data` to  Feishu APIs.

For example, if the incoming CloudEvent looks like:

```http
{
  "id" : "42d5b039-daef-4071-8584-e61df8fc1354",
  "source" : "vance-http-source",
  "specversion" : "V1",
  "type" : "http",
  "datacontenttype" : "application/json",
  "time" : "2022-05-17T18:44:02.681+08:00",
  "vancefeishusinkservice": "bot",
  "data" : {
    ...
  }
}
```

### Supported Feishu Service

- Bot: pushing a message to Group Chat with message of text, post, share_chat, image, and interactive.

## Quick Start

in this section, we show how to use Feishu Sink push a text message to your group chat.

### Add a bot to your group chat

Go to your target group, click Chat Settings > Group Bots > Add Bot, and select Custom Bot to add the bot to the group chat.

Enter a name and description for your bot, or set up an avatar for the bot, and then click "Add".

![add-a-bot](https://github.com/linkall-labs/vance-docs/raw/main/resources/connectors/sink-feishu-bot/add-a-bot.gif)

You will get the webhook address of the bot in the following format:

```
https://open.feishu.cn/open-apis/bot/v2/hook/xxxxxxxxxxxxxxxxx
```

> ⚠️ Please keep this webhook address properly. Do not publish it on GitHub, blogs, and other publicly accessible sites to avoid it being maliciously called to send spam messages.

![bot-config](https://github.com/linkall-labs/vance-docs/raw/main/resources/connectors/sink-feishu-bot/feishu-config.png)

> ⚠️ You must set your signature verification to make sure push messages work.

### Create Config file

replace `chat_group`, `signature`, and `address` to yours. `chat_group` can be fill in any value as you want.

```shell
cat << EOF > config.yml
enable: ["bot"]
bot:
  webhooks:
    - chat_group: "bot1"
      signature: "xxxxxxx"
      address: "https://open.feishu.cn/open-apis/bot/v2/hook/xxxxxxx"
EOF
```

### Start Using Docker

mapping 8080 to 31080 in order to avoid port conflict.

```shell
docker run -d -p 31080:8080 --rm \
  -v ${PWD}:/vance/config \
  --name sink-feishu public.ecr.aws/vanus/connector/sink-feishu:latest
```

### Test

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

now, you cloud see a notice in your chat group.
![received-notification](https://github.com/linkall-labs/vance/blob/main/connectors/sink-feishu/received-message.png)

### Clean

```shell
docker stop sink-feishu
```

## How to use

### Configuration

The default path is `/vance/config/config.yml`. if you want to change the default path, you can set env `CONNECTOR_CONFIG` to
tell Feishu Sink.


| Name                       | Required | Default | Description                                                                       |
|:---------------------------|:--------:|:-------:|-----------------------------------------------------------------------------------|
| enable                     | **YES**  |    -    | service list you want Feishu Sink is enabled                                      |
| bot.webhooks               | **YES**  |    -    | list of chat-group's configuration                                                |
| bot.webhooks.[].chat_group | **YES**  |    -    | the chat_group name, you can set any value to it                                  |
| bot.webhooks.[].signature  | **YES**  |    -    | the signature to sign reqeust, you can get it when you create Chat Bot            |
| bot.webhooks.[].address    | **YES**  |    -    | the webhook address that message sent to, you can get it when you create Chat Bot |

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

In some cases, users can't make sure how many bots there have or wanner send one message to multiple groups, which means they need to dynamically send message to
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
          #          For China mainland
          #          image: linkall.tencentcloudcr.com/vanus/connector/sink-feishu:latest
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
              mountPath: /vance/config
      volumes:
        - name: config
          configMap:
            name: sink-feishu
```

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md