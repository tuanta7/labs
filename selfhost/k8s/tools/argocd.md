# ArgoCD

Reference: [Getting Started](https://argo-cd.readthedocs.io/en/stable/getting_started/)

In most deployments, Argo CD is installed only once, usually in the `argocd` namespace, while applications are deployed into differents environments

```shell
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# grant access to other namespaces
kubectl label namespace dev  argocd.argoproj.io/managed-by=argocd
kubectl label namespace prod argocd.argoproj.io/managed-by=argocd

# get default admin password
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
```

By default, Argo CD's API server is not exposed outside the cluster. For development purposes, use kubectl port-forward to expose ArgoCD API Server

```sh
kubectl port-forward --address 0.0.0.0 svc/argocd-server -n argocd 8080:443
```

## 1. Ingress Configuration

## 2. Helm

### 2.1. Kustomize

### 2.2. Vault

Reference: [Vault on Kubernetes](https://developer.hashicorp.com/vault/tutorials/kubernetes/kubernetes-raft-deployment-guide)
