#!/bin/bash
set -euo pipefail

# ==============================
# CONFIGURATION
# ==============================
NEW_USER="griffin"
SSH_PORT=22
DOMAIN="lebron-games.com"
PROJECTS_DIR="/srv"
SERVER_IP=$(curl -s https://api.ipify.org)
DNS_TIMEOUT=600  # seconds to wait for DNS to resolve to this server

# ==============================
# UPDATE & HARDEN BASE SYSTEM
# ==============================
echo "🔧 Updating system..."
apt update -y && apt upgrade -y
apt install -y curl wget ufw fail2ban debian-keyring debian-archive-keyring apt-transport-https unattended-upgrades

echo "🔒 Enabling automatic security updates..."
dpkg-reconfigure -plow unattended-upgrades

# ==============================
# CREATE NON-ROOT USER
# ==============================
if ! id "$NEW_USER" &>/dev/null; then
    echo "👤 Creating non-root user '$NEW_USER'..."
    adduser --disabled-password --gecos "" "$NEW_USER"
    usermod -aG sudo "$NEW_USER"

    # Allow passwordless sudo for this user
    echo "$NEW_USER ALL=(ALL) NOPASSWD:ALL" > /etc/sudoers.d/$NEW_USER
    chmod 440 /etc/sudoers.d/$NEW_USER
fi

mkdir -p /home/$NEW_USER/.ssh
if [ -f ~/.ssh/authorized_keys ]; then
    cp ~/.ssh/authorized_keys /home/$NEW_USER/.ssh/
fi
chown -R $NEW_USER:$NEW_USER /home/$NEW_USER/.ssh
chmod 700 /home/$NEW_USER/.ssh
chmod 600 /home/$NEW_USER/.ssh/authorized_keys || true

# ==============================
# SSH HARDENING
# ==============================
echo "🛡️ Hardening SSH..."
sed -i 's/^#\?PermitRootLogin.*/PermitRootLogin no/' /etc/ssh/sshd_config
sed -i 's/^#\?PasswordAuthentication.*/PasswordAuthentication no/' /etc/ssh/sshd_config
# Uncomment below to change SSH port if desired
# sed -i "s/^#\?Port.*/Port $SSH_PORT/" /etc/ssh/sshd_config
systemctl restart ssh

# ==============================
# FIREWALL
# ==============================
echo "🔥 Configuring firewall..."
ufw default deny incoming
ufw default allow outgoing
ufw allow $SSH_PORT/tcp
ufw allow 80,443/tcp
ufw --force enable

# ==============================
# FAIL2BAN
# ==============================
echo "🚨 Enabling Fail2Ban..."
systemctl enable --now fail2ban

# ==============================
# DNS PRE-FLIGHT CHECK
# ==============================
echo "🌍 Checking DNS for $DOMAIN..."
START_TIME=$(date +%s)
while true; do
    DNS_IP=$(dig +short $DOMAIN | tail -n1)
    if [ "$DNS_IP" == "$SERVER_IP" ]; then
        echo "✅ DNS points to this server ($DNS_IP)"
        break
    fi

    NOW=$(date +%s)
    ELAPSED=$((NOW - START_TIME))
    if [ $ELAPSED -ge $DNS_TIMEOUT ]; then
        echo "❌ DNS did not resolve to this server after $DNS_TIMEOUT seconds. Exiting."
        exit 1
    fi

    echo "⏳ Waiting for DNS propagation... (Current: $DNS_IP, Expected: $SERVER_IP)"
    sleep 30
done

# ==============================
# CADDY INSTALL
# ==============================
echo "🌐 Installing Caddy..."
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | tee /etc/apt/sources.list.d/caddy-stable.list
apt update && apt install -y caddy

# ==============================
# MULTI-PROJECT DIRECTORY SETUP
# ==============================
echo "📁 Creating project directory structure..."
mkdir -p $PROJECTS_DIR
chown -R $NEW_USER:$NEW_USER $PROJECTS_DIR

# Global Caddyfile that imports all project configs
cat <<EOF > /etc/caddy/Caddyfile
{
    email admin@$DOMAIN
}

import $PROJECTS_DIR/*/Caddyfile
EOF

systemctl reload caddy

# Example project for verification
PROJECT_PATH="$PROJECTS_DIR/default"
mkdir -p "$PROJECT_PATH/www"
cat <<EOF > "$PROJECT_PATH/www/index.html"
<!DOCTYPE html>
<html>
<head><title>Hello</title></head>
<body><h1>Server setup complete for $DOMAIN 🚀</h1></body>
</html>
EOF

cat <<EOF > "$PROJECT_PATH/Caddyfile"
$DOMAIN {
    root * $PROJECT_PATH/www
    encode gzip zstd
    file_server
}
EOF

systemctl reload caddy

# ==============================
# FINAL MESSAGE
# ==============================
echo ""
echo "✅ Secure server setup complete!"
echo "User: $NEW_USER"
echo "Domain: $DOMAIN"
echo "Projects directory: $PROJECTS_DIR"
echo "Caddy auto-import enabled: yes"
echo ""
echo "👉 Next steps:"
echo "1. Add new projects under /srv/<project>/"
echo "2. Each project should have its own Caddyfile."
echo "3. Caddy will auto-enable HTTPS for every new site."
