# Ollama Setup on Mac

This guide details the steps required to set up **Ollama** (to run the Gemma model), configure a **reverse proxy using
Nginx**, and expose the service securely via a **Cloudflare Tunnel**.

## 1. Ollama and Gemma Model Setup

This [link](https://ai.google.dev/gemma/docs/integrations/ollama) provides the necessary command for installing Ollama and pulling the Gemma model. Once installed, you can verify the setup by running the Gemma model:

```sh
ollama run gemma:2b
```

⁉️ (Not sure) The Ollama API is designed to operate on a local machine. Public exposure is discouraged. Because of this constraint, any request that appears to originate from an external domain may be rejected pre-emptively. To make Ollama perceive the request as local, the proxy must override the Host header before forwarding.

## 2. Nginx Reverse Proxy Setup

Nginx will be used to run a reverse proxy for your Ollama service. This is often required for domain mapping, SSL
termination, and better security/management.

```sh
# install Nginx (using Homebrew on macOS)
brew install nginx

# start the Nginx service
brew services start nginx

# check service status
brew services
```

### Configuration

Use `nano` (or your preferred editor) to modify the main configuration file. The path below is typical for Homebrew
installations.

> [!TIP]
> To show hidden folders on a Mac, open Finder and press `Command + Shift + .` (period) to toggle their visibility.

```sh
nano /opt/homebrew/etc/nginx/nginx.conf
```

Inside the `http` block, you will typically define a `server` block to proxy requests to the default Ollama port (
`11434`).

```nginx
server {
    listen 80; # Or a custom port, e.g., 8080
    server_name your.local.domain.name; # Replace with your domain or IP

    location /o {
        rewrite ^/o/(.*)$ /$1 break;
        proxy_pass http://localhost:11434; # Ollama's default port
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

After saving the changes, test the configuration syntax and reload Nginx.

```sh
nginx -t
brew services reload nginx
```

## 3. Cloudflare Tunnel Setup

The setup requires the creation of several components in Cloudflare:

- A tunnel to serve as the secured connectivity channel.
- Service tokens, which will be used as non-interactive credentials for automated access.
- An Access policy configured with the Service Auth action; otherwise, Cloudflare Access will initiate an identity-provider authentication flow.
- An application configured to apply the above policy in order to protect the designated tunnel domain.

Install and run the tunnel service

```sh
# install cloudflared (using Homebrew on macOS)
brew install cloudflared 

# remove all existing tunnels
sudo cloudflared service uninstall

# create a new tunnel
sudo cloudflared service install eyJhI..xyz
```
