#!/bin/bash
set -u

# Multi-protocol VPN status helper
# Checks all supported VPN protocols

CONNECTED=false

# OpenVPN
if pgrep -x openvpn &>/dev/null; then
  CONNECTED=true
  echo "Protocol : OpenVPN"
  echo "Process  : openvpn (PID: $(pgrep -x openvpn))"
  if ip link show tun0 &>/dev/null; then
    VPN_IP=$(ip -4 addr show tun0 2>/dev/null | grep -oP 'inet \K[\d.]+' || echo "unknown")
    echo "Interface: tun0 ($VPN_IP)"
  else
    echo "Interface: tun0 not found (connecting...)"
  fi
  if [ -f /tmp/openvpn.log ]; then
    echo "Recent logs:"
    sudo tail -5 /tmp/openvpn.log 2>/dev/null || tail -5 /tmp/openvpn.log 2>/dev/null || true
  fi
fi

# WireGuard
if ip link show wg0 &>/dev/null 2>&1; then
  CONNECTED=true
  VPN_IP=$(ip -4 addr show wg0 2>/dev/null | grep -oP 'inet \K[\d.]+' || echo "unknown")
  echo "Protocol : WireGuard"
  echo "Interface: wg0 ($VPN_IP)"
  sudo wg show wg0 2>/dev/null | head -5 || true
fi

# IPsec/StrongSwan
if pgrep -x charon &>/dev/null; then
  IPSEC_STATUS=$(sudo ipsec status 2>/dev/null || true)
  if echo "$IPSEC_STATUS" | grep -q ESTABLISHED; then
    CONNECTED=true
    echo "Protocol : IPsec/IKEv2"
    echo "$IPSEC_STATUS" | grep -E 'ESTABLISHED|INSTALLED' || true
  fi
fi

# PPTP
if pgrep -x pppd &>/dev/null; then
  CONNECTED=true
  echo "Protocol : PPTP"
  echo "Process  : pppd (PID: $(pgrep -x pppd))"
  if ip link show ppp0 &>/dev/null; then
    VPN_IP=$(ip -4 addr show ppp0 2>/dev/null | grep -oP 'inet \K[\d.]+' || echo "unknown")
    echo "Interface: ppp0 ($VPN_IP)"
  else
    echo "Interface: ppp0 not found (connecting...)"
  fi
fi

if [ "$CONNECTED" = "false" ]; then
  echo "Status: DISCONNECTED"
  echo "No VPN connections active."
fi
