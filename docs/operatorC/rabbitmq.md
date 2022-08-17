---
title: RabbitMQ-Sample
nav_exclude: true
---

# RabbitMQ Sample

A RabbitMQ connector that will receive messages from a RabbitMQ queue and scale via Vance.

## Pre-requisites

- Vance and [Vance Pre-requisites][vance-pre]

## Setup

This setup will go through creating a RabbitMQ queue on the cluster and deploying a rabbitmq connector in Vance.
If you already have RabbitMQ you can use your existing queues.

First you should clone the project:

```cli
$ git clone https://github.com/JieDing/vance-docs
$ cd vance-docs
```

### Creating a RabbitMQ queue

#### [Install Helm](https://helm.sh/docs/using_helm/)

#### Install RabbitMQ via Helm

```cli
$ helm repo add bitnami https://charts.bitnami.com/bitnami
$ helm install rabbitmq --set auth.username=user --set auth.password=PASSWORD bitnami/rabbitmq --wait
```

#### Wait for RabbitMQ to deploy

⚠️ Be sure to wait until the deployment has completed before continuing. ⚠️

```cli
$ kubectl get po | grep rabbitmq
NAME         READY   STATUS    RESTARTS   AGE
rabbitmq-0   1/1     Running   0          3m3s
```

### Deploying a RabbitMQ connector

#### Deploy a connector
```cli
$ kubectl apply -f samples/sample-rabbitmq/rabbitmq-sample.yaml
```

#### Validate the connector has deployed
```cli
$ kubectl get deploy | grep rabbitmq
rabbitmq-connector-sample              0/0     0            0           16s
```

You should see `rabbitmq-connector-sample` deployment with 0 pods as there currently aren't any queue messages.
The pod number is scale to zero.

### Validating autoscaling

#### Publish messages to the queue

The following job will publish 200 messages to the queue.

```cli
$ kubectl apply -f samples/sample-rabbitmq/rabbitmq-publisher-job.yaml
```

#### Validate the pod scales

The rabbitmq-connector assumes that each connector pod can process 50 messages,
thus with 200 messages available in the queue, you can watch the pods scale out to 4 (200/50).

```cli
$ watch -n2 "kubectl get po | grep rabbitmq-connector"
Every 2.0s: kubectl get po | grep rabbitmq-connector
rabbitmq-connector-sample-854495f585-5k6fw              1/1     Running     0          10s
rabbitmq-connector-sample-854495f585-8cr9s              1/1     Running     0          10s
rabbitmq-connector-sample-854495f585-chlq5              1/1     Running     0          21s
rabbitmq-connector-sample-854495f585-z4tfl              1/1     Running     0          10s
```

After the queue is drained, you can find that the last replica will scale back down to zero.

#### Update the connector config

You can specify each connector pod should process 20 messages by updating `scalerSpec.metadata.value=20`
without cleaning up current resources.

```cli
$ vim samples/sample-rabbitmq/rabbitmq-sample.yaml
......
scalerSpec:
    checkInterval: 5
    cooldownPeriod: 5
    metadata:
      secret: rabbitmq-consumer-secret
      queueName: hello
      mode: QueueLength
      value: "20"
......
$ kubectl apply -f samples/sample-rabbitmq/rabbitmq-sample.yaml
secret/rabbitmq-consumer-secret unchanged
connector.vance.io/rabbitmq-connector-sample configured
```

Then we resend 200 messages to the queue, and this time the pods should scale out to 10 (200/20).

```cli
$ kubectl delete job rabbitmq-publish
$ kubectl apply -f samples/sample-rabbitmq/rabbitmq-publisher-job.yaml
```

```cli
$ watch -n2 "kubectl get po | grep rabbitmq-connector"
Every 2.0s: kubectl get po | grep rabbitmq-connector
rabbitmq-connector-sample-854495f585-5r2nl              1/1     Running     0          20s
rabbitmq-connector-sample-854495f585-9rgrb              1/1     Running     0          20s
rabbitmq-connector-sample-854495f585-9ttlt              1/1     Running     0          20s
rabbitmq-connector-sample-854495f585-gps7n              1/1     Running     0          5s
rabbitmq-connector-sample-854495f585-kmkl4              1/1     Running     0          43s
rabbitmq-connector-sample-854495f585-lvl57              1/1     Running     0          35s
rabbitmq-connector-sample-854495f585-n4cc2              1/1     Running     0          5s
rabbitmq-connector-sample-854495f585-rjg86              1/1     Running     0          35s
rabbitmq-connector-sample-854495f585-wgvvw              1/1     Running     0          20s
rabbitmq-connector-sample-854495f585-zsz8d              1/1     Running     0          35s
```

## Cleanup resources

```cli
$ kubectl delete job rabbitmq-publish
$ kubectl delete -f samples/sample-rabbitmq/rabbitmq-sample.yaml
$ helm delete rabbitmq
```

[vance-pre]: deploy.md#pre-requisites
