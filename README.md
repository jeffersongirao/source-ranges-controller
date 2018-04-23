# source-ranges-controller
[![Build Status](https://travis-ci.org/jeffersongirao/source-ranges-controller.png)](https://travis-ci.org/jeffersongirao/source-ranges-controller)
[![Go Report Card](https://goreportcard.com/badge/jeffersongirao/source-ranges-controller)](https://goreportcard.com/report/jeffersongirao/source-ranges-controller)


NOTE: This is an alpha-status project. We do regular tests on the code and functionality, but we can not assure a production-ready stability.

Source Ranges Controller sets loadBalancerSourceRanges to Kubernetes Services through a ConfigMap

## Requirements
Source Ranges controller is meant to be run on Kubernetes 1.8+.
All dependencies have been vendored, so there's no need to any additional download.

## Setup

You can directly create the controller deployment with kubectl:

```console
$ kubectl create -f https://raw.githubusercontent.com/jeffersongirao/source-ranges-controller/master/example/controller.yaml
```

Next, run an application and expose it via a Kubernetes Service:

```console
$ kubectl run nginx --image=nginx --replicas=1 --port=80
$ kubectl expose deployment nginx --port=80 --target-port=80 --type=LoadBalancer
```

Create a ConfigMap with the desired source ranges. Make sure to change `10.4.12.0/22` to the desired CIDR.

```console
$ kubectl create configmap whitelist --from-literal=office-network-1=10.4.12.0/22
```

Annotate the Service with the config map name holding the source ranges.

```console
$ kubectl annotate service nginx "source-ranges.alpha.girao.net/config-map=whitelist"
```