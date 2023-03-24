#!/bin/bash
set -e

PROJECT="dialoginsight"
PREFIX="/"
BIN="usr/bin"
ETC="etc/${PROJECT}"
DOC="usr/share/doc/${PROJECT}"
MAN="usr/share/man/man1"
INIT="usr/lib/systemd/system"

mkdir -p "${PREFIX}/${BIN}"
mkdir -p "${PREFIX}/${ETC}"
mkdir -p "${PREFIX}/${DOC}"
mkdir -p "${PREFIX}/${MAN}"
mkdir -p "${PREFIX}/${INIT}"

cp -a -f "${BIN}" "${PREFIX}/${BIN}"
cp -a -f "${DOC}" "${PREFIX}/${DOC}"
cp -a -f "${INIT}" "${PREFIX}/${INIT}"
cp -a -f "${MAN}" "${PREFIX}/${MAN}"
cp -a -n "${ETC}" "${PREFIX}/${ETC}"

for FILE in ${BIN}/*; do
	chown root:root "${PREFIX}/${BIN}/${FILE}"
done

for FILE in ${MAN}/*; do
	chown root:root "${PREFIX}/${MAN}/${FILE}"
done

for FILE in ${INIT}/*; do
	chown root:root "${PREFIX}/${INIT}/${FILE}"
done

chown root:root -R "${PREFIX}/${ETC}"
chown root:root -R "${PREFIX}/${DOC}"

adduser --system --quiet dialoginsight
systemctl daemon-reload
