---
title: ChatAi
---

# ChatAI Source

## Introduction

The ChatAI Source is a [Vanus Connector][vc] which aims to read an incoming request body, then call AI api and
convert the response content to a CloudEvent.

For example, the incoming Request looks like:

```bash
curl --location --request POST 'localhost:8080' \
--header '' \
--data-raw 'what is vanus'
```

which is converted to:

```json
{
  "specversion": "1.0",
  "id": "0effe4cc-06c7-4fe9-9180-aa7c3b30777e",
  "source": "vanus-chatAI-source",
  "type": "vanus-chatAI-type",
  "datacontenttype": "application/json",
  "time": "2023-03-28T09:15:10.70413Z",
  "data": {
    "message": "what is vanus", 
    "result": "vanus is a message queue"
  }
}
```

## Quick Start

This section will show you how to use ChatAI Source to read request body and call openai api to obtains response then
convert to a CloudEvent.

### Create Config file

```shell
cat << EOF > config.yml
target: http://localhost:31081
port: 8082
gpt
  token: xxxxxx
ernie_bot
  access_key: xxxxxx
  secret_key: xxxxxx
vanusai
  url: xxxxxx
  app_id: xxxxxx
EOF
```

| Name                 | Required | Default    | Description                                  |
|:---------------------|:--------:|:-----------|:---------------------------------------------|
| target               |   YES    |            | the target URL to send CloudEvents           |
| port                 |    NO    | 8080       | the port to receive HTTP request             |
| gpt.token            |   YES    |            | the ChatGPT auth token                       |
| ernie_bot.access_key |   YES    |            | the baidu ai [accessKey][ernie_bot]          |
| ernie_bot.secret_key |   YES    |            | the baidu ai [secretKey][ernie_bot]          |
| vanusai.url          |   YES    |            | vanus-ai url                                 |
| vanusai.app_id       |   YES    |            | vanus-ai application id                      |
| everyday_limit       |    NO    | 100        | the ChatAI Source call ai api count everyday |
| default_chat_mode    |    NO    | chatgpt    | chatgpt, wenxin, vanusai                     |

The ChatAI Source tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the
position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-chatai public.ecr.aws/vanus/connector/source-chatai
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
  --name source-chatai public.ecr.aws/vanus/connector/source-chatai
```

Open a terminal and use the following command to send a request to ChatAI Source

```shell
curl --location --request POST 'localhost:8082' \
--header 'Content-Type: text/plain' \
--data-raw 'what is vanus'
```

Here is the sort of CloudEvent you should expect to receive in the Display Sink:

```json
{
  "specversion": "1.0",
  "id": "0effe4cc-06c7-4fe9-9180-aa7c3b30777e",
  "source": "vanus-chatAI-source",
  "type": "vanus-chatAI-type",
  "datacontenttype": "application/json",
  "time": "2023-03-28T09:15:10.70413Z",
  "data": {
    "message": "what is vanus",
    "result": "vanus is a message queue"
  }
}
```

### Clean

```shell
docker stop source-chatai sink-display
```

## Source details

### Attributes

#### Changing Default Required Attributes

If you want to change the default attributes of `source`, `type` (defined
by [CloudEvents](https://github.com/cloudevents/spec/blob/main/cloudevents/spec.md#required-attributes)) to your own,
you could set the request header to set them.

| Header       | Description       |
|:-------------|:------------------|
| Vanus-Source | cloudevent source |
| Vanus-Type   | coudevent type    |

## Run in Kubernetes

```shell
kubectl apply -f source-chatai.yaml
```

```yaml
apiVersion: v1
kind: Service
metadata:
  name: source-chatai
  namespace: vanus
spec:
  selector:
    app: source-chatai
  type: ClusterIP
  ports:
    - port: 8080
      name: source-chatai
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: source-chatai
  namespace: vanus
data:
  config.yml: |-
    target: "http://localhost:18080"
    gpt:
      token: "sk-xxxxxx"
    ernie_bot:
      access_key: xxxxxxx
      secret_key: xxxxxx

---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: source-chatai
  namespace: vanus
  labels:
    app: source-chatai
spec:
  selector:
    matchLabels:
      app: source-chatai
  replicas: 1
  template:
    metadata:
      labels:
        app: source-chatai
    spec:
      containers:
        - name: source-chatai
          image: public.ecr.aws/vanus/connector/source-chatai:latest
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
            name: source-chatai
```

[vc]: https://docs.vanus.ai/introduction/concepts#vanus-connect
[ernie_bot]: https://ai.baidu.com/ai-doc/REFERENCE/Ck3dwjhhu
