# gc-controller

[![release](https://img.shields.io/github/release/DoodleScheduling/gc-controller/all.svg)](https://github.com/DoodleScheduling/gc-controller/releases)
[![release](https://github.com/doodlescheduling/gc-controller/actions/workflows/release.yaml/badge.svg)](https://github.com/doodlescheduling/gc-controller/actions/workflows/release.yaml)
[![report](https://goreportcard.com/badge/github.com/DoodleScheduling/gc-controller)](https://goreportcard.com/report/github.com/DoodleScheduling/gc-controller)
[![Coverage Status](https://coveralls.io/repos/github/DoodleScheduling/gc-controller/badge.svg?branch=master)](https://coveralls.io/github/DoodleScheduling/gc-controller?branch=master)
[![license](https://img.shields.io/github/license/DoodleScheduling/gc-controller.svg)](https://github.com/DoodleScheduling/gc-controller/blob/master/LICENSE)

Pod garbage collector controller.
This controller cleans evicted/failed pods and can keep a configurable number of them.
Unlike the vanilla gc collector this controller is workload aware and collects evicted pods by namespaces and can keep
a number of evicted pods for each governing workload resource.
Despite this the vanilla one is configured by default to collect only if there are more than `12500` evicted pods `--terminated-pod-gc-threshold`.
See https://kubernetes.io/docs/reference/command-line-tools-reference/kube-controller-manager/. That said this flag is usually not not even configurable on
hosted kubernetes platforms since the control plane can't be modified.

## Installation

### Helm

Please see [chart/gc-controller](https://github.com/DoodleScheduling/gc-controller/tree/master/chart/gc-controller) for the helm chart docs.

### Kustomize

Alternatively you may get the bundled manifests in each release to deploy it using kustomize or use them directly.

## Configure controller

You may change some settings using command line args.
**Note**: by default the garbace collection keeps 2 (`--keep=2`) evicted pods by workload but deletes (`--max-age=168h`) any evicted pod older than 1 week.

```
--concurrent int                            The number of concurrent Pod reconciles. (default 4)
--enable-leader-election                    Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.
--graceful-shutdown-timeout duration        The duration given to the reconciler to finish before forcibly stopping. (default 10m0s)
--health-addr string                        The address the health endpoint binds to. (default ":9557")
--insecure-kubeconfig-exec                  Allow use of the user.exec section in kubeconfigs provided for remote apply.
--insecure-kubeconfig-tls                   Allow that kubeconfigs provided for remote apply can disable TLS verification.
--keep int                                  The number of pods to keep for each workload. (default 2)
--kube-api-burst int                        The maximum burst queries-per-second of requests sent to the Kubernetes API. (default 300)
--kube-api-qps float32                      The maximum queries-per-second of requests sent to the Kubernetes API. (default 50)
--leader-election-lease-duration duration   Interval at which non-leader candidates will wait to force acquire leadership (duration string). (default 35s)
--leader-election-release-on-cancel         Defines if the leader should step down voluntarily on controller manager shutdown. (default true)
--leader-election-renew-deadline duration   Duration that the leading controller manager will retry refreshing leadership before giving up (duration string). (default 30s)
--leader-election-retry-period duration     Duration the LeaderElector clients should wait between tries of actions (duration string). (default 5s)
--log-encoding string                       Log encoding format. Can be 'json' or 'console'. (default "json")
--log-level string                          Log verbosity level. Can be one of 'trace', 'debug', 'info', 'error'. (default "info")
--max-age duration                          The number of pods to keep for each workload. (default 168h0m0s)
--max-retry-delay duration                  The maximum amount of time for which an object being reconciled will have to wait before a retry. (default 15m0s)
--metrics-addr string                       The address the metric endpoint binds to. (default ":9556")
--min-retry-delay duration                  The minimum amount of time for which an object being reconciled will have to wait before a retry. (default 750ms)
--watch-all-namespaces                      Watch for resources in all namespaces, if set to false it will only watch the runtime namespace. (default true)
--watch-label-selector string               Watch for resources with matching labels e.g. 'sharding.fluxcd.io/shard=shard1'.
```
