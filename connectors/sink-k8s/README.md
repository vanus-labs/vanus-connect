---
title: Kubernetes
---

# Kubernetes Sink

## Introduction

The Kubernetes Sink is a [Vanus Connector][vc] that aims to handle incoming CloudEvents in a way that extracts
the `data` part of the original event and operation kubernetes resource.

For example, if the incoming CloudEvent looks like:

```json
{
  "specversion": "1.0",
  "id": "4395ffa3-f6de-443c-bf0e-bb9798d26a1d",
  "source": "vanus.source.test",
  "type": "vanus.type.test",
  "datacontenttype": "application/json",
  "time": "2022-06-14T07:05:55.777689Z",
  "data": {
    "apiVersion": "batch/v1",
    "kind": "Job",
    "metadata": {
      "annotations": {
        "operation": "create"
      },
      "name": "job-test",
      "namespace": "default"
    },
    "spec": {
      "template": {
        "spec": {
          "containers": [
            {
              "command": [
                "sleep",
                "60s"
              ],
              "image": "busybox:latest",
              "name": "container1"
            }
          ],
          "restartPolicy": "Never"
        }
      },
      "ttlSecondsAfterFinished": 100
    }
  }
}
```

The Kubernetes Sink will extract `data` field write to Kubernetes cluster.

## Quickstart

### Prerequisites

- Have a Kubernetes cluster.

### Start with Kubernetes

```shell
kubectl apply -f sink-k8s.yaml
```

```yml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: sink-k8s-sa
  namespace: vanus
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: sink-k8s-cluster-role
rules:
  - apiGroups:
      - ""
      - "apps"
      - "batch"
    resources:
      - pods
      - jobs
      - cronjobs
      - daemonsets
      - deployments
      - statefulsets
    verbs:
      - create
      - get
      - list
      - watch
      - update
      - patch
      - delete
  - apiGroups:
      - apps
    resources:
      - deployments
      - statefulsets
    verbs:
      - get
      - list
      - create
      - update
      - patch
      - delete
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: sink-k8s-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: sink-k8s-cluster-role
subjects:
  - kind: ServiceAccount
    name: sink-k8s-sa
    namespace: vanus
---
apiVersion: v1
kind: Service
metadata:
  name: sink-k8s
  namespace: vanus
spec:
  selector:
    app: sink-k8s
  type: NodePort
  ports:
    - port: 8080
      nodePort: 31080
      name: http
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sink-k8s
  namespace: vanus
  labels:
    app: sink-k8s
spec:
  selector:
    matchLabels:
      app: sink-k8s
  replicas: 1
  template:
    metadata:
      labels:
        app: sink-k8s
    spec:
      containers:
        - name: sink-k8s
          image: public.ecr.aws/vanus/connector/sink-k8s
          imagePullPolicy: Always
          ports:
            - name: http
              containerPort: 8080
      serviceAccountName: sink-k8s-sa
```

### Test

Open a terminal and use following command to send a CloudEvent to the Sink.

Obtain your k8s cluster node INTERNAL-IP by following command and replace the following `<k8s node ip>`

```shell
kubectl get node -o wide
```


#### Example for create job

```shell
curl --location --request POST '<k8s node ip>:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
  "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
  "source": "quickstart",
  "specversion": "1.0",
  "type": "quickstart",
  "datacontenttype": "application/json",
  "time": "2022-10-26T10:38:29.345Z",
  "data": {
    "kind": "Job",
    "apiVersion": "batch/v1",
    "metadata": {
      "name": "job-test",
      "namespace": "default",
      "annotations": {
        "operation": "create"
      }
    },
    "spec": {
      "template": {
        "spec": {
          "containers": [
            {
              "name": "container1",
              "image": "busybox:latest",
              "command": [
                "sleep",
                "60s"
              ]
            }
          ],
          "restartPolicy": "Never"
        }
      },
      "ttlSecondsAfterFinished": 100
    }
  }
}'
```

### Example for delete job

```shell
curl --location --request POST '<k8s node ip>:31080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
  "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
  "source": "quickstart",
  "specversion": "1.0",
  "type": "quickstart",
  "datacontenttype": "application/json",
  "time": "2022-10-26T10:38:29.345Z",
  "data": {
    "kind": "Job",
    "apiVersion": "batch/v1",
    "metadata": {
      "name": "job-test",
      "namespace": "default",
      "annotations": {
        "operation": "delete"
      }
    },
    "spec": {
      "template": {
        "spec": {
          "containers": [
            {
              "name": "container1",
              "image": "busybox:latest",
              "command": [
                "sleep",
                "60s"
              ]
            }
          ],
          "restartPolicy": "Never"
        }
      },
      "ttlSecondsAfterFinished": 100
    }
  }
}'
```

### Clean resource

```shell
kubectl delete -f sink-k8s.yaml
```

## Integrate with Vanus

This section shows how a sink connector can receive CloudEvents from a
running [Vanus cluster](https://github.com/linkall-labs/vanus).

1. Run the sink-k8s.yaml

```shell
kubectl apply -f sink-k8s.yaml
```

2. Create an eventbus

```shell
vsctl eventbus create --name quick-start
```

3. Create a subscription (the sink should be specified as the sink service address or the host name with its port)

```shell
vsctl subscription create \
  --name quick-start \
  --eventbus quick-start \
  --sink 'http://sink-k8s:8080'
```

[vc]: https://www.vanus.dev/introduction/concepts#vanus-connect