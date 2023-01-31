---
title: GitHub
---

# GitHub Source

## Introduction

The GitHub Source is a [Vanus Connector](https://www.vanus.dev/introduction/concepts#vanus-connect) which aims to retrieve GitHub webhook events and transform them into CloudEvents based on the [CloudEvents Adapter specification](https://github.com/cloudevents/spec/blob/main/cloudevents/adapters/github.md) 
by wrapping the body of the original request into the data field.

An original GitHub webhook event looks like:
```JSON
{
  "action": "created",
  "starred_at": "2022-07-21T06:28:23Z",
  "repository": {
    "id": 513353059,
    "node_id": "R_kgDOHpklYw",
    "name": "test-repo",
    "full_name": "XXXX/test-repo",
    "private": false,
    "owner": {
      "login": "XXXX",
      "type": "User",
      "site_admin": false
    },
    "visibility": "public",
    "forks": 0,
    "open_issues": 2,
    "watchers": 1,
    "default_branch": "main"
  },
  "sender": {
    "login": "XXXX",
    "id": 75800782,
    "node_id": "MDQ6VXNlcjc1ODAwNzgy",
    "avatar_url": "https://avatars.githubusercontent.com/u/75800782?v=4",
    "gravatar_id": "",
    "url": "https://api.github.com/users/XXXX",
    "html_url": "https://github.com/XXXX",
    "followers_url": "https://api.github.com/users/XXXX/followers",
    "following_url": "https://api.github.com/users/XXXX/following{/other_user}",
    "gists_url": "https://api.github.com/users/XXXX/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/XXXX/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/XXXX/subscriptions",
    "organizations_url": "https://api.github.com/users/XXXX/orgs",
    "repos_url": "https://api.github.com/users/XXXX/repos",
    "events_url": "https://api.github.com/users/XXXX/events{/privacy}",
    "received_events_url": "https://api.github.com/users/XXXX/received_events",
    "type": "User",
    "site_admin": false
  }
}
```

which is converted to

```JSON
{
	id:"4ef226c0-08c7-11ed-998d-93772adf8abb", 
	source:"https://api.github.com/repos/XXXX/test-repo", 
	type:"com.github.watch.started", 
	datacontenttype:"application/json", 
	time:"2022-07-21T07:32:44.190Z", 
	data: {
       "action": "created", 
       ...
	}
}
```

## Quick Start

This section will teach you how to use GitHub Source to convert events from GitHub webhook into CloudEvents.

### Prerequisites

- Have a container runtime (i.e., docker).
- Have a GitHub Repository.

### Create the config file

```shell
cat << EOF > config.yml
target: http://localhost:31081
port: 8082
secret:
  github_webhook_secret: ""
EOF
```

| Name                     | Required | Default | Description                              |
|:-------------------------|:---------|:--------|:-----------------------------------------|
| target                   | YES      |         | the target URL to send CloudEvents       |
| port                     | YES      | 8080    | the port to receive GitHub webhook event |
| github_webhook_secret    | NO       |         | the GitHub webhook secret                |

The GitHub Source tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-github public.ecr.aws/vanus/connector/source-github
```

### Test
We have designed for you a sandbox environment, removing the need to use your local
machine. You can run Connectors directly and safely on the [Playground](https://play.linkall.com/).

1. We've already exposed the GitHub Source to the internet if you're using the Playground. Go to GitHub-Twitter Scenario under Payload URL.
![Payload img](https://raw.githubusercontent.com/linkall-labs/vanus-connect/main/connectors/source-github/payload.png)
2. Create a GitHub webhook for you repository.
   1. Create a webhook under the Settings tab inside your GitHub repository.
   2. Set the configuration for your webhook.

3. Open a terminal and use the following command to run a Display sink, which receives and prints CloudEvents.

```shell
docker run -it --rm \
  -p 31081:8080 \
  --name sink-display public.ecr.aws/vanus/connector/sink-display
```

Make sure the `target` value in your config file is `http://localhost:31081` so that the Source can send CloudEvents to our Display Sink.

4. Star your GitHub repository. 

Here is the sort of CloudEvent you should expect to receive in the Display Sink:
```json
{
  id:"4ef226c0-08c7-11ed-998d-93772adf8abb",
  source:"https://api.github.com/repos/XXXX/test-repo",
  type:"com.github.star.started",
  datacontenttype:"application/json",
  time:"2022-07-21T07:32:44.190Z",
  data: {
     "action": "created", 
     ...
  }
}
```

### Clean

```shell
docker stop source-github sink-display
```

## Run in Kubernetes

```shell
kubectl apply -f source-github.yaml
```

```yaml
apiVersion: v1
kind: Service
metadata:
   name: source-github
   namespace: vanus
spec:
   selector:
      app: source-github
   type: ClusterIP
   ports:
      - port: 8080
        name: source-github
---
apiVersion: v1
kind: ConfigMap
metadata:
   name: source-github
   namespace: vanus
data:
   config.yml: |-
      target: "http://vanus-gateway.vanus:8080/gateway/quick_start"
      port: 8080
      secret:
        github_webhook_secret: ""
---
apiVersion: apps/v1
kind: Deployment
metadata:
   name: source-github
   namespace: vanus
   labels:
      app: source-github
spec:
   selector:
      matchLabels:
         app: source-github
   replicas: 1
   template:
      metadata:
         labels:
            app: source-github
      spec:
         containers:
            - name: source-github
              image: public.ecr.aws/vanus/connector/source-github
              imagePullPolicy: Always
              ports:
                 - containerPort: 8080
                      name: github
              volumeMounts:
                 - name: source-github-config
                   mountPath: /vanus-connector/config
         volumes:
            - name: source-github-config
              configMap:
                 name: source-github
```

## Integrate with Vanus

This section shows how a source connector can send CloudEvents to a running [Vanus cluster](https://github.com/linkall-labs/vanus).

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

3. Update the target config of the GitHub Source
```yaml
target: http://192.168.49.2:30001/gateway/quick-start
```

4. Run the GitHub Source
```shell
docker run --network=host \
  --rm \
  -v ${PWD}:/vanus-connect/config \
  --name source-github public.ecr.aws/vanus/connector/source-github
  ```
