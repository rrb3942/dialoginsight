#/usr/bin/env bash
set -e

DATE=$(date -u)
BRANCH=$(git rev-parse --abbrev-ref HEAD)
REVISION=$(git rev-parse --short HEAD)
PROJECT="dialoginsight"
VERSION="0.11"
ARCH=$(dpkg --print-architecture)
PACKAGING="packaging"
BUILDDIR="build"
BUILDNAME="${PROJECT}_${VERSION}_${ARCH}"
BUILDROOT="${BUILDDIR}/${BUILDNAME}"
BUILDBIN="${BUILDROOT}/usr/bin"
BUILDETC="${BUILDROOT}/etc/dialoginsight"
BUILDDOC="${BUILDROOT}/usr/share/doc/dialoginsight"
BUILDMAN="${BUILDROOT}/usr/share/man/man1"
BUILDINIT="${BUILDROOT}/usr/lib/systemd/system"

echo "Cleaning ${BUILDDIR}"
rm -rf ${BUILDDIR}
echo "Making directories"
mkdir -p ${BUILDDIR}
mkdir -p ${BUILDROOT}
mkdir -p ${BUILDBIN}
mkdir -p ${BUILDETC}
mkdir -p ${BUILDDOC}
mkdir -p ${BUILDMAN}
mkdir -p ${BUILDINIT}

echo "Copying debian packaging"
cp -a ${PACKAGING}/DEBIAN ${BUILDROOT}/
sed -i "s/VERSION_NUMBER_TOKEN/${VERSION}/g" ${BUILDROOT}/DEBIAN/control
echo "Copying configuration files"
cp ${PACKAGING}/config.json ${BUILDETC}/
cp ${PACKAGING}/dialoginsight.service ${BUILDINIT}/
echo "Copying documentation"
cp Readme.md ${BUILDDOC}/
cp LICENSE ${BUILDDOC}/
pandoc --standalone --to man Readme.md | gzip --best > ${BUILDMAN}/dialoginsight.1.gz
echo "Building binary"
CGO_ENABLED=0 go build -ldflags "-w -s -X 'main.Version=${VERSION}' -X 'main.Built=${DATE}' -X 'main.Branch=${BRANCH}' -X 'main.Revision=${REVISION}' -X 'main.Compiler=$(go version)'" -o ${BUILDBIN}/${PROJECT}
echo "Building .deb package"
dpkg-deb --root-owner-group --build ${BUILDROOT}

echo "Building simple tar package"
cp ${PACKAGING}/install.sh ${BUILDROOT}/
rm -rf ${BUILDROOT}/DEBIAN
cd ${BUILDDIR}
tar -czf "${BUILDNAME}.tar.gz" ${BUILDNAME}
echo "${BUILDNAME}.tar.gz generated"

echo "Converting .deb to .rpm"
fpm -s deb -t rpm --rpm-rpmbuild-define "_buildhost localhost" *.deb

echo "Cleaning up ${BUILDROOT}"
rm -rf ${BUILDNAME}
