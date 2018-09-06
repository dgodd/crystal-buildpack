#!/bin/bash
set -euo pipefail

GO_VERSION="1.11"
export GoInstallDir="/tmp/go$GO_VERSION"
mkdir -p $GoInstallDir

if [ ! -f $GoInstallDir/go/bin/go ]; then
GO_SHA256="02d3f9cea35ef31b1f8869151a50be3d057d01271adcfd16ce888cd1c44cd1c3"
URL=https://buildpacks.cloudfoundry.org/dependencies/go/go${GO_VERSION}.linux-amd64-cflinuxfs2-${GO_SHA256:0:8}.tar.gz

echo "-----> Download go ${GO_VERSION}"
curl -s -L --retry 15 --retry-delay 2 $URL -o /tmp/go.tar.gz

DOWNLOAD_SHA256=$(shasum -a256 /tmp/go.tar.gz | cut -d ' ' -f 1)
if [[ $DOWNLOAD_SHA256 != $GO_SHA256 ]]; then
echo "       **ERROR** SHA256 mismatch: got $DOWNLOAD_SHA256 expected $GO_SHA256"
exit 1
fi

tar xzf /tmp/go.tar.gz -C $GoInstallDir
rm /tmp/go.tar.gz
fi
if [ ! -f $GoInstallDir/go/bin/go ]; then
echo "       **ERROR** Could not download go"
exit 1
fi
