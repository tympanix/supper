#!/bin/bash
set -e

if [ -d /run/systemd/system ]; then
    systemctl daemon-reload || :
    systemctl is-active --quiet supper && systemctl restart supper || :
fi
