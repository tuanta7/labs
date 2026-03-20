# K3s

Reference: [Quick-Start Guide](https://docs.k3s.io/quick-start)

K3s provides an installation script that is a convenient way to install it as a service on systemd or openrc based systems.

- Additional utilities will be installed including kubectl
- A single-node server installation is a fully-functional Kubernetes cluster, including all the datastore, control-plane, kubelet, and container runtime components necessary to host workload pods

```sh
curl -sfL https://get.k3s.io | sh -
```

To install additional agent nodes and add them to the cluster, run the installation script with the K3S_URL and K3S_TOKEN environment variables.

```sh
sudo cat /var/lib/rancher/k3s/server/node-token
curl -sfL https://get.k3s.io | K3S_URL=https://myserver:6443 K3S_TOKEN=mynodetoken sh -
```