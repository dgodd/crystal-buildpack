#!/usr/bin/env bash
# Runs the unit tests for this buildpack

set -euo pipefail

cd "$( dirname "${BASH_SOURCE[0]}" )/.."
source .envrc
./scripts/install_tools.sh

ginkgo -r -skipPackage=brats,integration,pkg .
