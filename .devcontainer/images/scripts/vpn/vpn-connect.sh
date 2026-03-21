#!/bin/bash
set -euo pipefail

# Multi-protocol VPN connect helper
# Supports: OpenVPN, WireGuard, IPsec/StrongSwan, PPTP
# Protocol auto-detected from running state or defaults to OpenVPN

# --- Disconnect any active VPN first (only one VPN at a time) ---
vpn_disconnect_all() {
  local disconnected=false
  if pgrep -x openvpn &>/dev/null; then
    echo "Disconnecting active OpenVPN..."
    sudo killall openvpn 2>/dev/null || true
    disconnected=true
  fi
  if ip link show wg0 &>/dev/null 2>&1; then
    echo "Disconnecting active WireGuard..."
    sudo wg-quick down wg0 2>/dev/null || true
    disconnected=true
  fi
  if pgrep -x charon &>/dev/null; then
    echo "Disconnecting active IPsec..."
    sudo ipsec stop 2>/dev/null || true
    disconnected=true
  fi
  if pgrep -x pppd &>/dev/null; then
    echo "Disconnecting active PPTP..."
    sudo killall pppd 2>/dev/null || true
    disconnected=true
  fi
  if [ "$disconnected" = "true" ]; then
    sleep 1
    rm -f /tmp/vpn-auth.txt
  fi
}

# --- Detect protocol from argument or env ---
PROTOCOL="${VPN_PROTOCOL:-openvpn}"

vpn_disconnect_all

case "$PROTOCOL" in
  openvpn)
    OVPN_CONFIG="${OPENVPN_CONFIG:-/home/vscode/.config/openvpn/client.ovpn}"
    OVPN_AUTH="${OPENVPN_AUTH:-/tmp/vpn-auth.txt}"

    if [ ! -f "$OVPN_CONFIG" ]; then
      echo "No OpenVPN config at: $OVPN_CONFIG"
      exit 1
    fi

    VPN_ARGS=(
      --config "$OVPN_CONFIG"
      --daemon ovpn-client
      --log /tmp/openvpn.log
      --script-security 2
      --up /etc/openvpn/update-dns
      --down /etc/openvpn/update-dns
      --keepalive 10 60
      --connect-retry 5
      --connect-retry-max 0
      --persist-tun
      --persist-key
      --resolv-retry infinite
    )
    [ -f "$OVPN_AUTH" ] && [ -s "$OVPN_AUTH" ] && VPN_ARGS+=(--auth-user-pass "$OVPN_AUTH")

    echo "Starting OpenVPN..."
    sudo openvpn "${VPN_ARGS[@]}"

    A=0
    while [ "$A" -lt 15 ]; do
      if ip link show tun0 &>/dev/null; then
        VPN_IP=$(ip -4 addr show tun0 2>/dev/null | grep -oP 'inet \K[\d.]+' || echo "unknown")
        echo "Connected (tun0: $VPN_IP)"
        exit 0
      fi
      sleep 1
      ((A++))
    done
    echo "tun0 not detected after 15s (check /tmp/openvpn.log)"
    exit 1
    ;;

  wireguard)
    WG_CONFIG="${WIREGUARD_CONFIG:-/home/vscode/.config/wireguard/wg0.conf}"
    if [ ! -f "$WG_CONFIG" ]; then
      echo "No WireGuard config at: $WG_CONFIG"
      exit 1
    fi

    echo "Starting WireGuard..."
    sudo wg-quick up "$WG_CONFIG"

    A=0
    while [ "$A" -lt 10 ]; do
      if ip link show wg0 &>/dev/null; then
        VPN_IP=$(ip -4 addr show wg0 2>/dev/null | grep -oP 'inet \K[\d.]+' || echo "unknown")
        echo "Connected (wg0: $VPN_IP)"
        exit 0
      fi
      sleep 1
      ((A++))
    done
    echo "wg0 not detected after 10s"
    exit 1
    ;;

  ipsec)
    IPSEC_CONFIG="${IPSEC_CONFIG:-/home/vscode/.config/strongswan/ipsec.conf}"
    if [ ! -f "$IPSEC_CONFIG" ]; then
      echo "No IPsec config at: $IPSEC_CONFIG"
      exit 1
    fi

    echo "Starting IPsec..."
    sudo cp "$IPSEC_CONFIG" /etc/ipsec.d/profile.conf
    sudo ipsec restart 2>/dev/null
    sleep 2

    CONN_NAME=$(grep -oP '(?<=^conn )\S+' "$IPSEC_CONFIG" | head -1)
    if [ -z "$CONN_NAME" ]; then
      echo "No connection name found in config"
      exit 1
    fi

    if sudo ipsec up "$CONN_NAME" 2>/dev/null; then
      echo "Connected via IPsec ($CONN_NAME)"
      exit 0
    fi
    echo "IPsec connection failed"
    exit 1
    ;;

  pptp)
    PPTP_CONFIG="${PPTP_CONFIG:-/home/vscode/.config/pptp/tunnel.conf}"
    if [ ! -f "$PPTP_CONFIG" ]; then
      echo "No PPTP config at: $PPTP_CONFIG"
      exit 1
    fi

    echo "Starting PPTP..."
    sudo pppd call tunnel nodetach </dev/null 2>&1 | sudo tee /tmp/pptp.log >/dev/null &

    A=0
    while [ "$A" -lt 15 ]; do
      if ip link show ppp0 &>/dev/null; then
        VPN_IP=$(ip -4 addr show ppp0 2>/dev/null | grep -oP 'inet \K[\d.]+' || echo "unknown")
        echo "Connected (ppp0: $VPN_IP)"
        exit 0
      fi
      sleep 1
      ((A++))
    done
    echo "ppp0 not detected after 15s (check /tmp/pptp.log)"
    exit 1
    ;;

  *)
    echo "Unknown protocol: $PROTOCOL"
    echo "Supported: openvpn, wireguard, ipsec, pptp"
    exit 1
    ;;
esac
