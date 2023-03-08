---
title: Amazon SNS
---

# Amazon SNS Source

## Introduction

The Amazon SNS(Simple Notification Service) Source is a [Vanus Connector][vc] which is designed to subscribe to the SNS topic and receive messages published to the topic, 
and then transform them into CloudEvents based on [CloudEvents Adapter specification][ceas].

Push is adopted by Amazon SNS to deliver messages from SNS topics to the endpoints. Therefore, the Amazon SNS Source should subscribe to the SNS topics and start an endpoint to receive messages from the SNS topics. 

Original SNS message pushed to http/https endpoints looks like:
```HTTP
 POST / HTTP/1.1
x-amz-sns-message-type: Notification
x-amz-sns-message-id: 22b80b92-fdea-4c2c-8f9d-bdfb0c7bf324
x-amz-sns-topic-arn: arn:aws:sns:us-west-2:123456789012:MyTopic
x-amz-sns-subscription-arn: arn:aws:sns:us-west-2:123456789012:MyTopic:c9135db0-26c4-47ec-8998-413945fb5a96
Content-Length: 773
Content-Type: text/plain; charset=UTF-8
Host: myhost.example.com
Connection: Keep-Alive
User-Agent: Amazon Simple Notification Service Agent

{
  "Type" : "Notification",
  "MessageId" : "22b80b92-fdea-4c2c-8f9d-bdfb0c7bf324",
  "TopicArn" : "arn:aws:sns:us-west-2:123456789012:MyTopic",
  "Subject" : "My First Message",
  "Message" : "Hello world!",
  "Timestamp" : "2012-05-02T00:54:06.655Z",
  "SignatureVersion" : "1",
  "Signature" : "EXAMPLEw6JRN...",
  "SigningCertURL" : "https://sns.us-west-2.amazonaws.com/SimpleNotificationService-f3ecfb7224c7233fe7bb5f59f96de52f.pem",
  "UnsubscribeURL" : "https://sns.us-west-2.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-west-2:123456789012:MyTopic:c9135db0-26c4-47ec-8998-413945fb5a96"
}
```
which is converted to:

```json
{
  "id" : "22b80b92-fdea-4c2c-8f9d-bdfb0c7bf324",
  "source" : "arn:aws:sns:us-west-2:123456789012:MyTopic:c9135db0-26c4-47ec-8998-413945fb5a96",
  "specversion" : "V1",
  "type" : "com.amazonaws.sns.Notification",
  "datacontenttype" : "application/json",
  "subject" : "My First Message",
  "time" : "2022-08-18T06:00:04.638Z",
  "data" : {
    "Type" : "Notification",
    "MessageId" : "22b80b92-fdea-4c2c-8f9d-bdfb0c7bf324",
    "TopicArn" : "arn:aws:sns:us-west-2:123456789012:MyTopic",
    "Subject" : "My First Message",
    "Message" : "Hello world!",
    "Timestamp" : "2012-05-02T00:54:06.655Z",
    "SignatureVersion" : "1",
    "Signature" : "EXAMPLEw6JRN...",
    "SigningCertURL" : "https://sns.us-west-2.amazonaws.com/SimpleNotificationService-f3ecfb7224c7233fe7bb5f59f96de52f.pem",
	"UnsubscribeURL" : "https://sns.us-west-2.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-west-2:123456789012:MyTopic:c9135db0-26c4-47ec-8998-413945fb5a96"
  }
}
```

## Quick Start

This section will show you how to use Amazon S3 Source to converts S3 events to a CloudEvent.

### Prerequisites

- Have a container runtime (i.e., docker).
- Have an AWS SNS Topic.
- AWS IAM [Access Key][accessKey].
- AWS permissions for the IAM user:
    - sns:Subscribe 
    - sns:ConfirmSubscription
    - sns:Unsubscribe
  
- Have a tool to expose the Source to internet which SNS can push message to.
  
  We have designed for you a sandbox environment, removing the need to use your local
  machine. You can run Connectors directly and safely on the [Playground](https://play.linkall.com/).
   
   We've already exposed webhook to the internet if you're using the Playground. Go to GitHub-Twitter Scenario under Payload URL.
   ![Payload img](https://raw.githubusercontent.com/vanus-labs/vanus-connect/main/connectors/source-github/payload.png)

### Create the config file

```shell
cat << EOF > config.yml
target: http://localhost:31081
aws:
  access_key_id: AKIAIOSFODNN7EXAMPLE
  secret_access_key: wJalrXUtnFEMI/K7MDENG/bPxRfiCY`EXAMPLEKEY
port: 8082  
sns_arn: "arn:aws:sns:us-west-2:843378899134:myTopic"
endpoint: "http://{internet access endpoint}"
protocol: "http"
EOF
```

| Name                  | Required | Default | Description                                               |
|:----------------------|:--------:|:-------:|:----------------------------------------------------------|
| target                |   YES    |         | the target URL to send CloudEvents                        |
| port                  |   YES    |  8080   | the port to receive SNS message                           |
| aws.access_key_id     |   YES    |         | the AWS IAM [Access Key][accessKey]                       |
| aws.secret_access_key |   YES    |         | the AWS IAM [Secret Key][accessKey]                       |
| sns_arn               |   YES    |         | the arn of the SNS topic                                  |
| endpoint              |   YES    |         | the SNS Source export internet url of http/https endpoint |
| protocol              |   YES    |         | the protocol used to subscribe SNS topic                  |


The Amazon SNS Source tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm --network=host \
  -v ${PWD}:/vanus-connect/config \
  --name source-aws-sns public.ecr.aws/vanus/connector/source-aws-sns
```

### Test

Open a terminal and use the following command to run a Display sink, which receives and prints CloudEvents.

```shell
docker run -it --rm \
  -p 31081:8080 \
  --name sink-display public.ecr.aws/vanus/connector/sink-display
```

Make sure the `target` value in your config file is `http://localhost:31081` so that the Source can send CloudEvents to our Display Sink.

Open [AWS SNS Console](https://us-west-2.console.aws.amazon.com/sns/v3/home?region=us-west-2#/topics), select the topic and publish a message.

Here is the sort of CloudEvent you should expect to receive in the Display Sink:

```json
{
  "id" : "22b80b92-fdea-4c2c-8f9d-bdfb0c7bf324",
  "source" : "arn:aws:sns:us-west-2:123456789012:MyTopic:c9135db0-26c4-47ec-8998-413945fb5a96",
  "specversion" : "V1",
  "type" : "com.amazonaws.sns.Notification",
  "datacontenttype" : "application/json",
  "subject" : "My First Message",
  "time" : "2022-08-18T06:00:04.638Z",
  "data" : {
    "Type" : "Notification",
    "MessageId" : "22b80b92-fdea-4c2c-8f9d-bdfb0c7bf324",
    "TopicArn" : "arn:aws:sns:us-west-2:123456789012:MyTopic",
    "Subject" : "My First Message",
    "Message" : "Hello world!",
    "Timestamp" : "2012-05-02T00:54:06.655Z",
    "SignatureVersion" : "1",
    "Signature" : "EXAMPLEw6JRN...",
    "SigningCertURL" : "https://sns.us-west-2.amazonaws.com/SimpleNotificationService-f3ecfb7224c7233fe7bb5f59f96de52f.pem",
	"UnsubscribeURL" : "https://sns.us-west-2.amazonaws.com/?Action=Unsubscribe&SubscriptionArn=arn:aws:sns:us-west-2:123456789012:MyTopic:c9135db0-26c4-47ec-8998-413945fb5a96"
  }
}
```

### Clean

```shell
docker stop source-aws-sns sink-display
```


## Run in Kubernetes

```shell
kubectl apply -f source-github.yaml
```

```yaml
apiVersion: v1
kind: Service
metadata:
  name: source-aws-sns
  namespace: vanus
spec:
  selector:
    app: source-aws-sns
  type: ClusterIP
  ports:
    - port: 8080
      name: source-aws-sns
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: source-aws-sns
  namespace: vanus
data:
  config.yml: |-
    "target": "http://vanus-gateway.vanus:8080/gateway/quick_start"
    aws:
      access_key_id: AKIAIOSFODNN7EXAMPLE
      secret_access_Key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY
    port: 8080
    sns_arn: "arn:aws:sns:us-west-2:843378899134:myTopic"
    endpoint: "http://ip10-0-188-4-ce3k58kdjmeg0u4hla2g-8082.direct.play.linkall.com"
    protocol: "http"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: source-aws-sns
  namespace: vanus
  labels:
    app: source-aws-sns
spec:
  selector:
    matchLabels:
      app: source-aws-sns
  replicas: 1
  template:
    metadata:
      labels:
        app: source-aws-sns
    spec:
      containers:
        - name: source-aws-sns
          image: public.ecr.aws/vanus/connector/source-aws-sns
          imagePullPolicy: Always
          ports:
            - containerPort: 8080
                name: http
          volumeMounts:
            - name: config
              mountPath: /vanus-connect/config
      volumes:
        - name: config
          configMap:
            name: source-aws-sns
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

3. Update the target config of the Amazon SNS Source
```yaml
target: http://192.168.49.2:30001/gateway/quick-start
```

4. Run the Amazon SNS Source
```shell
kubectl apply -f source-aws-sns.yaml
```

[vc]: https://docs.vanus.ai/introduction/concepts#vanus-connect
[ceas]: https://github.com/cloudevents/spec/blob/main/cloudevents/adapters/aws-sns.md
[accessKey]: https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_access-keys.html
