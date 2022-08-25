# Kafka Source 

## Overview


## User Guidelines

### Connector Introduction



## Kafka Source Configs


### Config Fields of the Kafka Source

## Kafka Source Secrets

Users should set their sensitive data Base64 encoded in a secret file. And mount your local secret file to `/vance/secret/secret.json` when you run the connector.

### Encode your sensitive data

```shell
$ echo -n ABCDEFG | base64
QUJDREVGRw==
```

Replace 'ABCDEFG' with your sensitive data.

### Set your local secret file


## Local Development

You can run the source codes of the Kafka Source locally as well.

### Building via Maven



### Running via Maven

⚠️ NOTE: For better local development and test, the connector can also read configs from `main/resources/config.json`. So, you don't need to 
declare any environment variables or mount a config file to `/vance/config/config.json`. Same logic applies to `main/resources/secret.json` as well.

[vc]: https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md
[config]: https://github.com/linkall-labs/vance-docs/blob/main/docs/connector.md
[ceas]: https://github.com/cloudevents/spec/blob/main/cloudevents/adapters/aws-s3.md