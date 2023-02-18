#/usr/bin/env bash
set -e

DATE=$(date +%Y%m%d)
BRANCH=$(git rev-parse --abbrev-ref HEAD)
REVISION=$(git rev-parse --short HEAD)


#Cleanup any leftover
rm -rf dialoginsight_${BRANCH}_${DATE}_${REVISION}
rm -f dialoginsight_${BRANCH}_${DATE}_${REVISION}.tar.gz

mkdir dialoginsight_${BRANCH}_${DATE}_${REVISION}
cp -a packaging/* dialoginsight_${BRANCH}_${DATE}_${REVISION}/
cp Readme.md dialoginsight_${BRANCH}_${DATE}_${REVISION}/
CGO_ENABLED=0 go build -ldflags "-w -s" 
mv dialoginsight dialoginsight_${BRANCH}_${DATE}_${REVISION}/
tar -czf dialoginsight_${BRANCH}_${DATE}_${REVISION}.tar.gz dialoginsight_${BRANCH}_${DATE}_${REVISION}
rm -rf dialoginsight_${BRANCH}_${DATE}_${REVISION}
