#!/bin/bash
set +e
# Container-friendly DNS update for OpenVPN
# Writes directly to /etc/resolv.conf (no systemd-resolved/D-Bus dependency)
# Called by OpenVPN via --up and --down script hooks

[ "$script_type" ] || exit 0
[ "$dev" ] || exit 0

BACKUP="/etc/resolv.conf.ovpn-backup"

case "$script_type" in
  up)
    # Backup original resolv.conf (Docker DNS)
    [ ! -f "$BACKUP" ] && cp /etc/resolv.conf "$BACKUP"

    # Parse OpenVPN DHCP options
    NMSRVRS=""
    SRCHS=""
    # shellcheck disable=SC2086
    for optionvarname in $(printf '%s\n' ${!foreign_option_*} | sort -t _ -k 3 -g); do
      option="${!optionvarname}"
      # shellcheck disable=SC2086
      set -- $option
      if [ "$1" = "dhcp-option" ]; then
        [ "$2" = "DNS" ] && NMSRVRS="${NMSRVRS:+$NMSRVRS }$3"
        [ "$2" = "DOMAIN" ] && SRCHS="${SRCHS:+$SRCHS }$3"
      fi
    done

    # Write new resolv.conf (VPN DNS + Docker fallback)
    {
      [ -n "$SRCHS" ] && echo "search $SRCHS"
      for ns in $NMSRVRS; do echo "nameserver $ns"; done
      # Keep Docker DNS as fallback
      grep '^nameserver' "$BACKUP" 2>/dev/null || true
    } >/etc/resolv.conf
    ;;
  down)
    # Restore original resolv.conf
    [ -f "$BACKUP" ] && mv "$BACKUP" /etc/resolv.conf
    ;;
esac
