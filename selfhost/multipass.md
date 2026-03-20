# Multipass Setup

Reference: [Manage instances](https://documentation.ubuntu.com/multipass/latest/how-to-guides/manage-instances/)

> [!NOTE]
> This document describes the initialization of a lightweight Ubuntu Server environment using Multipass on a macOS-based Mac Mini M1.

The resulting instance is intended to operate as a headless server, accessed exclusively via SSH. Network connectivity is established through a bridged Wi-Fi interface with a manually assigned (static) IP address to ensure consistent reachability on the local network.

## MTU Configuration

```sh
# Check MTU value
ifconfig

# Find a workable number
ping -D -s 1480 google.com

# Update MTU value
sudo ifconfig eth0 mtu 1450
```

## Instance Lifecycle Management (Clean State)

The following commands are used to inspect, remove, and permanently purge Multipass instances in order to maintain a clean and reproducible server environment.

```sh
# List all existing instances
multipass list

# Mark a specific instance as deleted (recoverable state)
multipass delete <instance-name>

# Permanently delete a specific instance in a single step
multipass delete --purge <instance-name>

# Delete all instances
multipass delete --all
```

## Ubuntu Server Launch with Bridged Networking

A bridged network configuration is used so that the Ubuntu Server instance appears as a first-class device on the local network. This enables direct SSH access using a fixed IP address assigned at the network level.

```sh
# Display available network interfaces
multipass networks

# Inspect the currently configured bridged network
multipass get local.bridged-network

# Configure the bridged network to use the Wi-Fi interface (example: en1)
multipass set local.bridged-network=en1

# Launch an Ubuntu Server instance with bridged networking enabled
multipass launch --name server --bridged

# Access server instance
multipass shell server
```

```sh
# Get the number of logical cores
sysctl -n hw.logicalcpu

# Adding CPU cores, memory & disk
multipass set local.server.cpus=4
multipass set local.server.memory=4G
multipass set local.server.disk=128G
```

Server connectivity was verified under the following network configuration: the client address was 192.168.1.11, the multipass host (Mac Mini) address was 192.168.1.137, and the multipass server instance address was 192.168.1.2.

- A Redis container was started using Docker: docker run -d --name some-redis -p 6379:6379 redis.
- Connectivity was validated by testing access to redis://192.168.1.2:6379/0.

```sh
tuanta@M1 ~ % multipass info server
Name:           server
State:          Running
Snapshots:      0
IPv4:           192.168.64.7
                192.168.1.2
                172.17.0.1
Release:        Ubuntu 24.04.3 LTS
Image hash:     a40713938d74 (Ubuntu 24.04 LTS)
CPU(s):         4
Load:           0.00 0.00 0.00
Disk usage:     2.9GiB out of 123.9GiB
Memory usage:   326.9MiB out of 3.8GiB
Mounts:
```

## Next Steps

- Install Docker: [Install Docker Engine on Ubuntu](https://docs.docker.com/engine/install/ubuntu/)
- Install Kubectl (ARM64): [Install kubectl on Linux](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/)
