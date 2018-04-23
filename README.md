# source-ranges-controller
[![Build Status](https://travis-ci.org/jeffersongirao/source-ranges-controller.png)](https://travis-ci.org/jeffersongirao/source-ranges-controller)
[![Go Report Card](https://goreportcard.com/badge/jeffersongirao/source-ranges-controller)](https://goreportcard.com/report/jeffersongirao/source-ranges-controller)


NOTE: This is an alpha-status project. We do regular tests on the code and functionality, but we can not assure a production-ready stability.

Source Ranges Controller sets loadBalancerSourceRanges to Kubernetes Services through a ConfigMap

## Requirements
Source Ranges controller is meant to be run on Kubernetes 1.8+.
All dependencies have been vendored, so there's no need to any additional download.