# Feishu-bot Sink 

## Introduction

The Feishu-bot Sink is a [Vance Connector][vc] which aims to handle incoming CloudEvents in a way that extracts the `data` part of the
original event and deliver these extracted `data` to the webhook URL of a Feishu custom robot.

For example, if the incoming CloudEvent looks like:

```http
{
  "id" : "42d5b039-daef-4071-8584-e61df8fc1354",
  "source" : "vance-http-source",
  "specversion" : "V1",
  "type" : "http",
  "datacontenttype" : "application/json",
  "time" : "2022-05-17T18:44:02.681+08:00",
  "data" : {
    "myData" : "simulation event data <1>"
  }
}
```

The Feishu-bot Sink will POST an HTTP request to the Feishu bot which looks like:

``` json
> POST /payload HTTP/2

> Host: localhost:8080
> User-Agent: VanceCDK-HttpClient/1.0.0
> Content-Type: application/json
> Content-Length: 102

> {
>   "timestamp": "1599360473",
>   "sign": "xxxxxxxxxxxxxxxxxxxxx",
>   "msg_type": "text",
>   "content": {
>      "text": "simulation event data <1>"
>   }
> }
```

## Feishu-bot Sink Configs

Users can specify their configs by either setting environments variables or mount a config.json to
`/vance/config/config.json` when they run the connector. Find examples of setting configs [here][config].

### Config Fields of the Feishu-bot Sink

| Configs   | Description                                                            | Example                 |
|:----------|:-----------------------------------------------------------------------|:------------------------|
| v_target  | v_target is used to specify the target URL HTTP Sink will send data to. In this case, the target should be the webhook URL of the Feishu bot.| "https://open.feishu.cn/open-apis/bot/v2/hook/......" |
| v_port    | v_port is used to specify the port HTTP Sink is listening on           | "8080"                  |
| feishu_secret  | feishu_secret is used to specify the signature used to verify whether http requests are valid.         | "****************"                  |

## Add a Custom Bot to Your Feishu Group

Go to your target group, click Chat Settings > Group Bots > Add Bot, and select Custom Bot to add the bot to the group chat.

Enter a name and description for your bot, or set up an avatar for the bot, and then click "Add".

![add-a-bot](https://github.com/linkall-labs/vance-docs/raw/main/resources/connectors/sink-feishu-bot/add-a-bot.gif)

You will get the webhook address of the bot in the following format:

```
https://open.feishu.cn/open-apis/bot/v2/hook/xxxxxxxxxxxxxxxxx
```

The webhook URL should be used as the field of `v_target` for Feishu-bot configs.

> ⚠️ Please keep this webhook address properly. Do not publish it on GitHub, blogs, and other publicly accessible sites to avoid it being maliciously called to send spam messages.

![bot-config](https://github.com/linkall-labs/vance-docs/raw/main/resources/connectors/sink-feishu-bot/feishu-config.png)

>  ⚠️ You must set your signature verification to use this connector.

The signature should be used as the field of `feishu_secret` for Feishu-bot configs.

## Feishu-bot Sink Image

> docker.io/vancehub/sink-feishubot

## Local Development

You can run the sink codes of the Feishu-bot Sink locally as well.

### Building via Maven

```shell
$ cd sink-feishu-bot
$ mvn clean package
```

### Running via Maven

```shell
$ mvn exec:java -Dexec.mainClass="com.linkall.sink.feishubot.Entrance"
```

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md
[config]: https://github.com/linkall-labs/vance-docs/blob/main/docs/connector.md