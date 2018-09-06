#!/bin/bash
set -euo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."
source .envrc

if [ ! -f .bin/ginkgo ]; then
go build -o .bin/ginkgo github.com/onsi/ginkgo/ginkgo
fi
if [ ! -f .bin/buildpack-packager ]; then
go build -o .bin/buildpack-packager github.com/cloudfoundry/libbuildpack/packager/buildpack-packager
fi
