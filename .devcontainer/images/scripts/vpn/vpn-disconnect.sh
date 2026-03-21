#!/bin/bash
set -euo pipefail

# Multi-protocol VPN disconnect helper
# Auto-detects and stops any active VPN connection

DISCONNECTED=false

# OpenVPN
if pgrep -x openvpn &>/dev/null; then
  echo "Stopping OpenVPN..."
  sudo killall openvpn 2>/dev/null || true
  A=0
  while [ "$A" -lt 10 ]; do
    pgrep -x openvpn &>/dev/null || break
    sleep 1
    ((A++))
  done
  echo "OpenVPN disconnected."
  DISCONNECTED=true
fi

# WireGuard
if ip link show wg0 &>/dev/null 2>&1; then
  echo "Stopping WireGuard..."
  sudo wg-quick down wg0 2>/dev/null \
    || sudo wg-quick down /home/vscode/.config/wireguard/wg0.conf 2>/dev/null || true
  echo "WireGuard disconnected."
  DISCONNECTED=true
fi

# IPsec/StrongSwan
if pgrep -x charon &>/dev/null; then
  echo "Stopping IPsec..."
  CONN_NAME=$(grep -oP '(?<=^conn )\S+' /etc/ipsec.d/profile.conf 2>/dev/null | head -1 || echo "")
  if [ -n "$CONN_NAME" ]; then
    sudo ipsec down "$CONN_NAME" 2>/dev/null || true
  fi
  sudo ipsec stop 2>/dev/null || true
  echo "IPsec disconnected."
  DISCONNECTED=true
fi

# PPTP
if pgrep -x pppd &>/dev/null; then
  echo "Stopping PPTP..."
  sudo killall pppd 2>/dev/null || true
  A=0
  while [ "$A" -lt 10 ]; do
    pgrep -x pppd &>/dev/null || break
    sleep 1
    ((A++))
  done
  echo "PPTP disconnected."
  DISCONNECTED=true
fi

# Cleanup
rm -f /tmp/vpn-auth.txt
sudo rm -f /etc/ipsec.d/profile.conf /etc/ipsec.d/profile.secrets 2>/dev/null || true

if [ "$DISCONNECTED" = "false" ]; then
  echo "No VPN is running."
fi
