#!/bin/bash
useradd -r dialoginsight &>/dev/null

set -e
mkdir -p "/etc/dialoginsight"
cp -f dialoginsight "/usr/bin/dialoginsight"
chown root:root "/usr/bin/dialoginsight"
cp -f dialoginsight.service /usr/lib/systemd/system/dialoginsight.service
chown root:root /usr/lib/systemd/system/dialoginsight.service
cp -n config.json "/etc/dialoginsight/config.json"
chown dialoginsight:dialoginsight -R "/etc/dialoginsight"

systemctl daemon-reload
