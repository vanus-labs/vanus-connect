---
title: HttpServer-Sample
nav_exclude: true
---

# HttpServer Sample

An echo server that will receive http requests and scale via Vance.

## Pre-requisites

- Vance and [Vance Pre-requisites][vance-pre]

## Setup

This setup will go through installing an `ingress-nginx` Ingress Controller on the cluster
and deploying a httpserver sample in Vance. If you already have an Ingress Controller, you can use your existing one.

First you should clone the project:

```cli
$ git clone https://github.com/JieDing/vance-docs
$ cd vance-docs
```

### Install an ingress-nginx Ingress Controller

#### [Install Helm](https://helm.sh/docs/using_helm/)

#### Install ingress-nginx Ingress Controller via Helm

```cli
$ helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
$ helm repo update
$ helm install ingress-nginx ingress-nginx/ingress-nginx -n ingress-nginx --create-namespace
```

#### Wait for Ingress Controller to deploy

⚠️ Be sure to wait until the deployment has completed before continuing. ⚠️

```cli
$ kubectl get po -n ingress-nginx
NAME                                        READY   STATUS    RESTARTS   AGE
ingress-nginx-controller-7575567f98-kq9dm   1/1     Running   0          62s
```

### Deploying a HTTP sample

#### Deploy the sample
```cli
$ kubectl apply -f samples/sample-httpserver/httpserver-sample.yaml
```

#### Validate the httpServer has deployed
```cli
$ kubectl get deploy | grep http
httpserver-sample   0/0     0            0           37s
```

You should see `httpserver-sample` deployment with 0 pods as there currently aren't any traffics to the httpserver.
The pod number is scale to zero.

### Validating autoscaling

#### Curl the httpServer

Now that you have your application running, you can issue an HTTP request.
To do so, you'll need to know the IP address to request.

```cli
$ kubectl get svc -n ingress-nginx
NAME                                 TYPE           CLUSTER-IP       EXTERNAL-IP     PORT(S)                      AGE
ingress-nginx-controller             LoadBalancer   10.107.35.172    10.107.35.172   80:32521/TCP,443:31302/TCP   123m
$ curl -H "Host: myhost.com" <EXTERNAL-IP>
```

If you're using `minikube`, you may find that the `EXTERNAL-IP` of your ingres-controller is always `<pending>`.

```cli
$ kubectl get svc -n ingress-nginx
NAME                                 TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
ingress-nginx-controller             LoadBalancer   10.107.35.172    <pending>     80:32521/TCP,443:31302/TCP   124m
```

Try to use `minikube tunnel` in a separate terminal and keep it running to expose your `EXTERNAL-IP`.

```cli
$ minikube tunnel
Status:
machine: minikube
pid: 496077
route: 10.96.0.0/12 -> 192.168.49.2
minikube: Running
services: [httpserver-sample, ingress-nginx-controller]
errors:
minikube: no errors
router: no errors
loadbalancer emulator: no errors
......
```

#### Validate the pod scales

We can use [hey] to generate artificial HTTP load.
This generates 200 requests to the ingress controller.

```cli
$ kubectl delete -f samples/sample-httpserver/httpserver-sample.yaml
$ kubectl apply -f samples/sample-httpserver/httpserver-sample.yaml
$ hey -c 200 -n 200 -host myhost.com http://<EXTERNAL-IP>
```

The httpserver-sample assumes that each pod can process 20 requests,
thus with 200 concurrent requests, you can watch the pods scale out to 10 (200/20).

```cli
$ watch -n2 "kubectl get po | grep http"
Every 2.0s: kubectl get po | grep http
httpserver-sample-5bdf5c46d-6q4pt   1/1     Running   0          36s
httpserver-sample-5bdf5c46d-8x49f   1/1     Running   0          21s
httpserver-sample-5bdf5c46d-c4mz7   1/1     Running   0          36s
httpserver-sample-5bdf5c46d-dpjjv   1/1     Running   0          6s
httpserver-sample-5bdf5c46d-lj4tq   1/1     Running   0          21s
httpserver-sample-5bdf5c46d-lpt82   1/1     Running   0          6s
httpserver-sample-5bdf5c46d-mp4rb   1/1     Running   0          21s
httpserver-sample-5bdf5c46d-mxxvg   1/1     Running   0          38s
httpserver-sample-5bdf5c46d-nwtxl   1/1     Running   0          21s
httpserver-sample-5bdf5c46d-tfcph   1/1     Running   0          36s
```

After about 300s, you can find that the last replica will scale back down to zero.

## Cleanup resources

```cli
$ kubectl delete -f samples/sample-httpserver/httpserver-sample.yaml
```

[vance-pre]: deploy.md#pre-requisites
[hey]: https://github.com/rakyll/hey
