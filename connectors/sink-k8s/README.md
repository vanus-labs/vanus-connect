---
title: Kubernetes
---

# Kubernetes Sink

## Introduction

The Kubernetes Sink is a [Vance Connector](https://github.com/linkall-labs/vance-docs/blob/main/docs/concept.md), 
which now supports create/delete operations for kubernetes resource.

For example, if the incoming CloudEvent looks like:

```json
{
  "specversion": "1.0",
  "id": "4395ffa3-f6de-443c-bf0e-bb9798d26a1d",
  "source": "vance.source.test",
  "type": "vance.type.test",
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

## Kubernetes Sink Configs

Users can specify their configs by either setting environments variables or mount a config.yml to
`/vance/config/config.yml` when they run the connector.

### Config Fields of Kubernetes Sink

| name | requirement | description                                                                    |
|------|-------------|--------------------------------------------------------------------------------|
| port | required    | port is used to specify the port Kubernetes Sink is listening on, default 8080 |

## Elasticsearch Sink Image

> public.ecr.aws/vanus/connector/sink-k8s

## Deploy

### using k8s(recommended)

```shell
kubectl apply -f https://github.com/linkall-labs/vance/blob/main/connectors/sink-k8s/sink-k8s.yml
```

## Examples

### create job

```shell
vsctl event put quick-start \
  --data-format json \
  --data "{\"kind\":\"Job\",\"apiVersion\":\"batch/v1\",\"metadata\":{\"name\":\"job-test\",\"namespace\":\"default\",\"annotations\":{\"operation\":\"create\"}},\"spec\":{\"template\":{\"spec\":{\"containers\":[{\"name\":\"container1\",\"image\":\"busybox:latest\",\"command\":[\"sleep\",\"60s\"]}],\"restartPolicy\":\"Never\"}},\"ttlSecondsAfterFinished\":100}}" \
  --id "1" \
  --source "quick-start" \
  --type "examples"
```

### delete job

```shell
vsctl event put quick-start \
  --data-format json \
  --data "{\"kind\":\"Job\",\"apiVersion\":\"batch/v1\",\"metadata\":{\"name\":\"job-test\",\"namespace\":\"default\",\"annotations\":{\"operation\":\"delete\"}},\"spec\":{\"template\":{\"spec\":{\"containers\":[{\"name\":\"container1\",\"image\":\"busybox:latest\",\"command\":[\"sleep\",\"60s\"]}],\"restartPolicy\":\"Never\"}},\"ttlSecondsAfterFinished\":100}}" \
  --id "1" \
  --source "quick-start" \
  --type "examples"
```
