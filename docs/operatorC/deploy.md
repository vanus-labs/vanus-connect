
---
title: Vance-Installation
nav_exclude: true
---

# Vance Installation

Vance is composed of a set of APIs implemented as Kubernetes Custom Resource Definitions (CRDs) and a controller.

## Pre-requisites

- install a Kubernetes cluster
- install [KEDA](https://keda.sh/docs/2.7/deploy/)
- install [KEDA-http](https://github.com/kedacore/http-add-on/blob/main/docs/install.md)

## Install with YAML files

### Install

```
kubectl apply -f deploy/vance-1.0.0.yaml
```

### Uninstall

```
kubectl delete -f deploy/vance-1.0.0.yaml
```

### Verify the installation

The all-in-one YAML file will create CRDs and run the Vance controller in the `vance` namespace.

```
$ kubectl get crds | grep vance
connectors.vance.io                     2022-05-15T07:50:35Z
```

```
$ kubectl get po -n vance
NAME                                        READY   STATUS    RESTARTS      AGE
vance-controller-manager-6d454547f9-lscvv   2/2     Running   4 (80s ago)   11m
```
