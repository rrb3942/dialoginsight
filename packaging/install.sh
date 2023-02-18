#!/bin/bash
useradd -r dialoginsight &>/dev/null

set -e
mkdir -p "/etc/dialogsight"
cp -f dialoginsight "/usr/bin/dialogsight"
chown root:root "/usr/bin/dialogsight"
cp -f dialoginsight.service /usr/lib/systemd/system/dialoginsight.service
chown root:root /usr/lib/systemd/system/dialoginsight.service
cp -n config.json "/etc/dialoginsight/config.json"
chown dialoginsight:dialoginsight -R "/etc/dialogsight"

systemctl daemon-reload
