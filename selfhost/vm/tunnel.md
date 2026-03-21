# Cloudflare Tunnel

First, add the Cloudflare apt repository and GPG key so cloudflared can be installed via apt.

```sh
# Add Cloudflare GPG key
sudo mkdir -p --mode=0755 /usr/share/keyrings
curl -fsSL https://pkg.cloudflare.com/cloudflare-public-v2.gpg \
  | sudo tee /usr/share/keyrings/cloudflare-public-v2.gpg >/dev/null

# Add Cloudflare repo to apt sources
echo 'deb [signed-by=/usr/share/keyrings/cloudflare-public-v2.gpg] https://pkg.cloudflare.com/cloudflared any main' | sudo tee /etc/apt/sources.list.d/cloudflared.list
```

Then install cloudflared from the repository and quickly verify that it is available.

```sh
# Install cloudflared
sudo apt-get update && sudo apt-get install cloudflared

# Quick version check
cloudflared --version

# Optional: check the systemd service status
sudo systemctl status cloudflared
```

To have something to expose through the tunnel, start a simple nginx container and confirm it is listening.

```sh
# Create a demo HTTP endpoint on port 18080
sudo docker run --name test-endpoint -p 18080:80 -d nginx

# Verify the container is reachable
ss -tunlp | grep 18080
curl http://localhost:18080
```

Next, log in to [Cloudflare](https://dash.cloudflare.com/) and create a named tunnel that will later be mapped to hostnames; then run the tunnel in the foreground for testing or as a systemd service for long‑running usage.

```sh
# Run tunnel in the foreground (good for testing)
cloudflared tunnel run --token eyJhI...

# Install as a systemd service
sudo cloudflared service install eyJhI...

# Start and enable the service
sudo systemctl start cloudflared
sudo systemctl enable cloudflared

# Tail logs while debugging
sudo journalctl -u cloudflared -f
```

If you no longer need the tunnel on this machine, you can remove the service and optionally clean up local configuration.

```sh
# Stop and remove the systemd service
sudo systemctl stop cloudflared
sudo cloudflared service uninstall

# Optional: remove local Cloudflare tunnel config
sudo rm -rf ~/.cloudflared/
```

For the WARP client on Linux, package resolution issues are often due to the Ubuntu codename used in the repository configuration.

```sh
Error: Unable to locate package cloudflare-warp
```

To fix this, try changing the Ubuntu codename in the WARP apt source (for example, replacing `victoria` or `questing` with `jammy`) as suggested in this [Cloudflare community thread](https://community.cloudflare.com/t/warp-cli-linux-error/549725).
