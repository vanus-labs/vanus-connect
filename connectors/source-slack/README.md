---
title: Slack
---

# Slack Source

## Introduction

The Slack Source is a [Vanus Connector][vc] which aims to convert a Slack event to a CloudEvent.

For example, the Slack event is converted to:

```json
{
  "specversion": "1.0",
  "id": "5f986f77-1934-4f48-ad09-235f7e534cfd",
  "source": "https://github.com/vanus-labs/vanus-connect/connectors/source-slack",
  "type": "event_callback",
  "datacontenttype": "application/json",
  "time": "2023-04-24T05:57:11.047718Z",
  "eventtype": "message",
  "data": {
    "type": "event_callback",
    "token": "XXYYZZ",
    "team_id": "TXXXXXXXX",
    "api_app_id": "AXXXXXXXXX",
    "event_context": "EC12345",
    "event_id": "Ev08MFMKH6",
    "event_time": 1234567890,
    "authorizations": [
      {
        "enterprise_id": "E12345",
        "team_id": "T12345",
        "user_id": "U12345",
        "is_bot": false,
        "is_enterprise_install": false
      }
    ],
    "is_ext_shared_channel": false,
    "context_team_id": "TXXXXXXXX",
    "context_enterprise_id": null,
    "event": {
      "blocks": [
        {
          "block_id": "YrP",
          "elements": [
            {
              "elements": [
                {
                  "text": "just for test message",
                  "type": "text"
                }
              ],
              "type": "rich_text_section"
            }
          ],
          "type": "rich_text"
        }
      ],
      "channel": "CXXXXXX",
      "channel_type": "channel",
      "client_msg_id": "4606762f-bc02-4012-9eff-4222c683b0c9",
      "event_ts": "1682315827.695329",
      "team": "TXXXXXXXX",
      "text": "just for test message",
      "ts": "1682315827.695329",
      "type": "message",
      "user": "UXXXXXX"
    }
  }
}
```

## Quick Start

This section will show you how to use Slack Source to convert a Slack webhook event to a CloudEvent.

### Create Config file

```shell
cat << EOF > config.yml
target: http://localhost:31081
port: 8082
verify_token: xxxxxx
signing_secret: xxxxxx
EOF
```

| Name            | Required | Default | Description                                                       |
|:----------------|:--------:|:-------:|:------------------------------------------------------------------|
| target          |   YES    |         | the target URL to send CloudEvents                                |
| port            |    NO    |  8080   | the port to receive webhook event request                         |
| verify_token    |   YES    |         | the Slack webhook verify token                                    |
| signing_secret  |   YES    |         | the signing secret for check webhook header X-Slack-Signature     |

The Slack Source tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the
position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-slack public.ecr.aws/vanus/connector/source-slack
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

Configure Slack webhook and add scope and test, then the Display sink will receive a
CloudEvents.

### Clean

```shell
docker stop source-slack sink-display
```

## Source details

### Attributes

#### Extension Attributes

The Slack Source defines
following [CloudEvents Extension Attributes](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#extension-context-attributes)

| Attribute |  Type  | Description                               |
|:---------:|:------:|:------------------------------------------|
| eventtype | string | the slack event.type,for example: message |

[vc]: https://docs.vanus.ai/introduction/concepts#vanus-connect
[api] https://api.slack.com/apis/connections/events-api

