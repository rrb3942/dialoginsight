#!/bin/bash
set -e

PROJECT="dialoginsight"
PREFIX="/"
BIN="usr/bin"
ETC="etc/${PROJECT}"
DOC="usr/share/doc/${PROJECT}"
MAN="usr/share/man/man1"
INIT="usr/lib/systemd/system"

mkdir -p "${PREFIX}${BIN}"
mkdir -p "${PREFIX}${ETC}"
mkdir -p "${PREFIX}${DOC}"
mkdir -p "${PREFIX}${MAN}"
mkdir -p "${PREFIX}${INIT}"

cp -f ${BIN}/* "${PREFIX}${BIN}"
cp -f ${DOC}/* "${PREFIX}${DOC}"
cp -f ${INIT}/* "${PREFIX}${INIT}"
cp -f ${MAN}/* "${PREFIX}${MAN}"
cp -n ${ETC}/* "${PREFIX}${ETC}"

for FILE in ${BIN}/*; do
	echo "${PREFIX}${BIN}/${FILE##*/}"
	chown root:root "${PREFIX}${BIN}/${FILE##*/}"
done

for FILE in ${MAN}/*; do
	chown root:root "${PREFIX}${MAN}/${FILE##*/}"
done

for FILE in ${INIT}/*; do
	chown root:root "${PREFIX}${INIT}/${FILE##*/}"
done

chown root:root -R "${PREFIX}${ETC}"
chown root:root -R "${PREFIX}${DOC}"

adduser --system dialoginsight --shell /sbin/nologin &> /dev/null
systemctl daemon-reload
