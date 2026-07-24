# KubeHealth Operator

KubeHealth Operator is a Kubernetes Operator built using **Kubebuilder** and **Go** that continuously monitors the health of Pods in a Kubernetes namespace.

The operator watches a custom resource named **KubeHealth** and periodically checks the Pods in the configured namespace. Based on the Pod status and container restart counts, it updates the Custom Resource status with a real-time health summary.

---

## Features

- Monitor Pods in any Kubernetes namespace
- Configurable restart threshold
- Counts healthy Pods
- Counts failed/unhealthy Pods
- Detects frequently restarting Pods
- Updates Custom Resource status automatically
- Periodic reconciliation every 30 seconds
- Built using Kubernetes Operator pattern

---

## Architecture

```
          +-----------------------+
          |   KubeHealth CR       |
          +----------+------------+
                     |
                     |
                     v
         +------------------------+
         |  KubeHealth Operator   |
         |     (Controller)       |
         +-----------+------------+
                     |
          Lists Pods from Namespace
                     |
                     v
        +--------------------------+
        | Kubernetes API Server    |
        +--------------------------+
                     |
                     v
               Pod Health Report
                     |
                     v
        Updates KubeHealth Status
```

---

## Custom Resource Example

```yaml
apiVersion: apps.abhishek.dev/v1
kind: KubeHealth
metadata:
  name: kubehealth-sample

spec:
  namespace: default
  restartThreshold: 3
```

### Spec Fields

| Field | Description |
|-------|-------------|
| namespace | Namespace to monitor |
| restartThreshold | Maximum allowed restart count before a Pod is considered unhealthy |

---

## Status Fields

The operator continuously updates the CR status.

| Field | Description |
|-------|-------------|
| healthyPods | Number of healthy Pods |
| failedPods | Number of failed/unhealthy Pods |
| restartingPods | Pods whose restart count exceeds the configured threshold |
| clusterStatus | Healthy / Degraded |
| lastChecked | Timestamp of the latest health check |

---

## Prerequisites

- Go 1.24+
- Docker
- kubectl
- Kubernetes or OpenShift cluster
- Kubebuilder
- Operator SDK

---

## Install CRD

```bash
make install
```

---

## Deploy the Operator

```bash
make deploy IMG=<your-image>
```

---

## Create a KubeHealth Resource

```bash
kubectl apply -f config/samples/apps_v1_kubehealth.yaml
```

---

## Verify

List the Custom Resource

```bash
kubectl get kubehealth
```

View complete status

```bash
kubectl describe kubehealth kubehealth-sample
```

or

```bash
kubectl get kubehealth kubehealth-sample -o yaml
```

---

## Example Expected Status

> **Note:** The following status is an example of the values updated by the operator after reconciliation.

```yaml
status:
  healthyPods: 5
  failedPods: 1
  restartingPods: 2
  clusterStatus: Degraded
  lastChecked: "2026-07-23T11:25:00Z"
```

---

## Project Structure

```
api/
    v1/
        kubehealth_types.go

cmd/
    main.go

internal/
    controller/
        kubehealth_controller.go

config/
    crd/
    manager/
    rbac/
    samples/
```

---

## Future Improvements

- Pod Events monitoring
- CrashLoopBackOff detection
- Prometheus metrics
- Alerting through Slack / Email
- Namespace selector support
- Multi-namespace monitoring
- Grafana dashboard integration

---

## Tech Stack

- Go
- Kubernetes
- Kubebuilder
- Operator SDK
- Controller Runtime
- OpenShift

---

## License

Licensed under the Apache License 2.0.