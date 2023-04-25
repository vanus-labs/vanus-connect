---
title: Schedule
---

# Schedule Source

## Introduction

The Schedule Source is a [Vanus Connector][vc] which aims to schedule make a CloudEvent which likes

```json
{
  "id": "ef26ed7b-9377-4bf5-b8d4-4fc6347e4fa2",
  "source": "vanus.ai/schedule",
  "specversion": "1.0",
  "type": "schedule",
  "datacontenttype": "application/json",
  "time": "2022-12-05T09:00:42.618Z",
  "data": {}
}
```

## Quick Start

This section shows how Schedule Source make a CloudEvent.

### Create the config file

```shell
cat << EOF > config.yml
target: http://localhost:31081
cron: "* * * * * *"
EOF
```

| Name      | Required | Default | Description                                                            |
|:----------|:---------|:--------|:-----------------------------------------------------------------------|
| target    | YES      |         | the target URL which Schedule Source will send CloudEvents to          |
| cron      | YES      |         | the schedule [cron], second,minute,hour,day of month,month,day of week |
| time_zone | NO       |         | the schedule [time zone][tz], default is your server zone              |

...

The Schedule Source tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the
position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-schedule public.ecr.aws/vanus/connector/source-schedule
```

### Test

Open a terminal and use the following command to run a Display sink, which receives and prints CloudEvents.

```shell
docker run -it --rm \
  -p 31081:8080 \
  --name sink-display public.ecr.aws/vanus/connector/sink-display
```

Make sure the `target` value in your config file is `http://localhost:31081` so that the Source can send CloudEvents to
our Display Sink.

Here is the sort of CloudEvent you should expect to receive in the Display Sink:

```json
{
  "id": "ef26ed7b-9377-4bf5-b8d4-4fc6347e4fa2",
  "source": "vanus.ai/schedule",
  "specversion": "1.0",
  "type": "schedule",
  "datacontenttype": "application/json",
  "time": "2022-12-05T09:00:42.618Z",
  "data": {}
}
```

### Clean

```shell
docker stop source-schedule sink-display
```

[vc]: https://www.vanus.dev/introduction/concepts#vanus-connect
[cron]: https://en.wikipedia.org/wiki/Cron
[tz]: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones