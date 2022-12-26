---
title: GitHub
---

# GitHub Source

## Overview
A Vance Connector which retrieves GitHub webhooks events, transform them into CloudEvents and deliver CloudEvents to the target URL.

## User Guidelines

## Connector Introduction
The GitHub Source is a [Vance Connector](https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md) designed to retrieves
GitHub webhooks events in various format, transform them into CloudEvents based on [CloudEvents Adapter specification](https://github.com/cloudevents/spec/blob/main/cloudevents/adapters/github.md) and wrap the body of the original request into the data of CloudEvents.

The original GitHub webhooks events look like:
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
This GitHub star event will be transformed into a CloudEvents like:
```JSON
CloudEvent:{
	id:"4ef226c0-08c7-11ed-998d-93772adf8abb", 
	source:"https://api.github.com/repos/XXXX/test-repo", 
	type:"com.github.watch.started", 
	datacontenttype:"application/json", 
	time:"2022-07-21T07:32:44.190Z", 
	data:JsonCloudEventData{
		"http request body"
	}
}
```
## GitHub Source Configs
Users can specify their configs by either setting environments variables or mount a config.json to /vance/config/config.json when they run the connector. Find examples of setting configs [here](https://github.com/linkall-labs/vance-docs/blob/main/docs/connector.md).
### Config Fields of the GitHub Source
|  Configs    |  Description    																  |  Example    			  |  Required    |
|  :----:     |  :----:         																  |  :----:     			  |  :----:      |
|  v_target   |  v_target is used to specify the target URL HTTP Source will send CloudEvents to  |  "http://localhost:8081"  |  YES  		 |
|  v_port     |  v_port is used to specify the port HTTP Source is listening on					  |  "8080"	                  |  YES         |
## GitHub Source Secrets
Users should set their sensitive data Base64 encoded in a secret file. And mount your local secret file to /vance/secret/secret.json when you run the connector.
### Encode your sensitive data
```Bash
$ echo -n ABCDEFG | base64
QUJDREVGRw==
```
Replace 'ABCDEFG' with your sensitive data.
### Set your local secret file
```Bash
$ cat secret.json
{
  "githubWebHookSecret": "${githubWebHookSecret}"
}
```
|  Secrets         		 |  Description    																  |  Example    			  |  Required    |
|  :----:     			 |  :----:         																  |  :----:     			  |  :----:      |
|  githubWebHookSecret   |  The githubWebHookSecret is used to verify your webhook secret key		      |  "YWJjZGU="				  |  YES  		 |
## GitHub Source Image
>    
### Run the GitHub-source image in a container
Mount your local config file and secret file to specific positions with -v flags.
```Bash
docker run -v $(pwd)/secret.json:/vance/secret/secret.json -v $(pwd)/config.json:/vance/config/config.json -p 8081:8081 source-github
```
