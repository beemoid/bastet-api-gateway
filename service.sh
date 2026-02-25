#!/usr/bin/env bash
# =============================================================================
# service.sh — Manage api-gateway as a systemd service
#
# Usage:
#   sudo ./service.sh install    Install and enable the service
#   sudo ./service.sh uninstall  Remove the service
#        ./service.sh start      Start the service
#        ./service.sh stop       Stop the service
#        ./service.sh restart    Restart the service
#        ./service.sh status     Show running status
#        ./service.sh log        Tail live logs (Ctrl+C to exit)
#        ./service.sh log-all    Show full log history
# =============================================================================

set -euo pipefail

SERVICE_NAME="bastet-api-gateway"
UNIT_FILE="/etc/systemd/system/${SERVICE_NAME}.service"

# Resolve the directory this script lives in (the deployment root)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINARY="${SCRIPT_DIR}/${SERVICE_NAME}"
ENV_FILE="${SCRIPT_DIR}/.env"

# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------
need_root() {
    if [[ $EUID -ne 0 ]]; then
        echo "✗  This command must be run as root (sudo ./service.sh $1)"
        exit 1
    fi
}

check_binary() {
    if [[ ! -x "${BINARY}" ]]; then
        echo "✗  Binary not found or not executable: ${BINARY}"
        echo "   Run 'make dist' and copy the dist package to this server."
        exit 1
    fi
}

check_env() {
    if [[ ! -f "${ENV_FILE}" ]]; then
        echo "✗  .env file not found: ${ENV_FILE}"
        echo "   Copy .env.example → .env and fill in your values."
        exit 1
    fi
}

# ---------------------------------------------------------------------------
# install
# ---------------------------------------------------------------------------
cmd_install() {
    need_root install
    check_binary
    check_env

    echo "→ Installing ${SERVICE_NAME} as a systemd service..."

    # Determine the user that should run the service
    # Defaults to the owner of the binary; override by setting SERVICE_USER env var
    SERVICE_USER="${SERVICE_USER:-$(stat -c '%U' "${BINARY}")}"

    cat > "${UNIT_FILE}" <<EOF
[Unit]
Description=BASTET API Gateway
After=network.target
Wants=network-online.target

[Service]
Type=simple
User=${SERVICE_USER}
WorkingDirectory=${SCRIPT_DIR}
EnvironmentFile=${ENV_FILE}
ExecStart=${BINARY}
Restart=on-failure
RestartSec=5s
StandardOutput=journal
StandardError=journal
SyslogIdentifier=${SERVICE_NAME}

# Hardening (optional — comment out if they cause issues)
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable "${SERVICE_NAME}"

    echo "✓  Service installed and enabled."
    echo "   Run: sudo ./service.sh start"
}

# ---------------------------------------------------------------------------
# uninstall
# ---------------------------------------------------------------------------
cmd_uninstall() {
    need_root uninstall

    if systemctl is-active --quiet "${SERVICE_NAME}"; then
        systemctl stop "${SERVICE_NAME}"
    fi

    systemctl disable "${SERVICE_NAME}" 2>/dev/null || true
    rm -f "${UNIT_FILE}"
    systemctl daemon-reload

    echo "✓  Service removed."
}

# ---------------------------------------------------------------------------
# start / stop / restart
# ---------------------------------------------------------------------------
cmd_start() {
    check_binary
    check_env
    systemctl start "${SERVICE_NAME}"
    echo "✓  ${SERVICE_NAME} started."
    sleep 1
    systemctl status "${SERVICE_NAME}" --no-pager || true
}

cmd_stop() {
    systemctl stop "${SERVICE_NAME}"
    echo "✓  ${SERVICE_NAME} stopped."
}

cmd_restart() {
    check_binary
    check_env
    systemctl restart "${SERVICE_NAME}"
    echo "✓  ${SERVICE_NAME} restarted."
    sleep 1
    systemctl status "${SERVICE_NAME}" --no-pager || true
}

# ---------------------------------------------------------------------------
# status
# ---------------------------------------------------------------------------
cmd_status() {
    systemctl status "${SERVICE_NAME}" --no-pager || true
}

# ---------------------------------------------------------------------------
# log
# ---------------------------------------------------------------------------
cmd_log() {
    echo "→ Live logs for ${SERVICE_NAME} (Ctrl+C to exit)..."
    journalctl -u "${SERVICE_NAME}" -f
}

cmd_log_all() {
    journalctl -u "${SERVICE_NAME}" --no-pager
}

# ---------------------------------------------------------------------------
# Dispatch
# ---------------------------------------------------------------------------
COMMAND="${1:-help}"

case "${COMMAND}" in
    install)   cmd_install   ;;
    uninstall) cmd_uninstall ;;
    start)     cmd_start     ;;
    stop)      cmd_stop      ;;
    restart)   cmd_restart   ;;
    status)    cmd_status    ;;
    log)       cmd_log       ;;
    log-all)   cmd_log_all   ;;
    *)
        echo ""
        echo "Usage: ./service.sh <command>"
        echo ""
        echo "  install    Install and enable as a systemd service (requires sudo)"
        echo "  uninstall  Remove the systemd service (requires sudo)"
        echo "  start      Start the service"
        echo "  stop       Stop the service"
        echo "  restart    Restart the service"
        echo "  status     Show current status"
        echo "  log        Tail live logs (Ctrl+C to exit)"
        echo "  log-all    Print full log history"
        echo ""
        ;;
esac
