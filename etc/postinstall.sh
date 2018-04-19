<<<<<<< HEAD
#!/bin/bash
=======
#!/bin/sh
>>>>>>> dd2823f2432abda1e12873a2320197c4f5cc0158
set -e

if [ -d /run/systemd/system ]; then
    systemctl daemon-reload || :
    systemctl is-active --quiet supper && systemctl restart supper || :
fi
