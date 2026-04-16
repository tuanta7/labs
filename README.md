# Code Labs

A collection of code samples, prototyping, and experimental projects.

Theory notes and documentation are available in these repositories:

- Go Language: [intervievv/golang](https://github.com/intervievv/golang)
- System Design Concepts: [intervievv/system-design](https://github.com/intervievv/system-design)

## Covered Topics

## Quickstart Environment Setup

### Golang

Reference: [Download and install Go](https://go.dev/doc/install)

```shell
# go to home directory
cd ~

# install Go
wget https://go.dev/dl/go1.25.0.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.25.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
go version

# back to the previous directory
cd -
```

### NodeJS (nvm)

```shell
# install nvm
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash

export NVM_DIR="$HOME/.nvm"
[ -s "$NVM_DIR/nvm.sh" ] && \. "$NVM_DIR/nvm.sh" # This loads nvm
[ -s "$NVM_DIR/bash_completion" ] && \. "$NVM_DIR/bash_completion" # This loads nvm bash_completion

source ~/.bashrc
nvm --version

# install node
nvm install --lts
```

### Docker

Reference: [Install Docker Engine](https://docs.docker.com/engine/install/)

Using the `apt` repository

```sh
# uninstall all conflicting packages
sudo apt remove $(dpkg --get-selections docker.io docker-compose docker-compose-v2 docker-doc podman-docker containerd runc | cut -f1)

# Add Docker's official GPG key:
sudo apt update
sudo apt install ca-certificates curl
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to apt sources:
sudo tee /etc/apt/sources.list.d/docker.sources <<EOF
Types: deb
URIs: https://download.docker.com/linux/ubuntu
Suites: $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}")
Components: stable
Signed-By: /etc/apt/keyrings/docker.asc
EOF

sudo apt update

# install the latest version
sudo apt install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# verify the installation
sudo docker run hello-world
```

Notes: The Docker daemon binds to a Unix socket, not a TCP port. By default the socket is owned by root; add users to the `docker` group to allow non-root use.

```sh
# create the docker group
sudo groupadd docker

# add user to the docker group.
sudo usermod -aG docker $USER

# create new user if needed
sudo adduser newuser
```