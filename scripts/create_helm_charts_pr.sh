#!/bin/bash

# This script is used to open a PR against this repository:
# https://github.com/SUSE/kubernetes-charts-suse-com
# It will fetch the bundle specified by SIDECAR_BUNDLE and will update the helm
# charts in the above repo.

set -e

if [ -z "$GITHUB_USER"  ]; then
  echo "GITHUB_USER environment variable not set"
  exit 1
fi

if [ -z "$GITHUB_PASSWORD" ] && [ -z "${GITHUB_TOKEN}" ] ; then
  echo "GITHUB_PASSWORD environment variable not set"
  exit 1
fi

if [[ -z "${SIDECAR_BUNDLE}" ]]; then
  echo "SIDECAR_BUNDLE not set"
  exit 1
fi

# GitHub organization can be overridden to test the script
GITHUB_ORGANIZATION="${GITHUB_ORGANIZATION:-SUSE}"

TMP_WORKDIR=$(mktemp -d -p "${PWD}")
trap "rm -rf '${TMP_WORKDIR}'" EXIT

pushd "${TMP_WORKDIR}"

# Get the "hub" cli
# https://hub.github.com/hub.1.html
if type -p hub ; then
  HUB="$(type -p hub)"
else
  wget -O - https://github.com/github/hub/releases/download/v2.3.0-pre10/hub-linux-amd64-2.3.0-pre10.tgz | tar xvz --wildcards --strip-components=2 '*/bin/hub'
  HUB="${PWD}/hub"
  chmod +x "${HUB}"
fi

# Clone the kubernetes-charts-suse-com github repo
"${HUB}" clone "git@github.com:${GITHUB_ORGANIZATION}/kubernetes-charts-suse-com.git"

# The bundle name is <helm chart name>-<version>.tgz
# where the helm chart name is "cf-usb-sidecar-<thing>"
# So we need to recover that... somehow

wget -O bundle.tgz "${SIDECAR_BUNDLE}"
versioned_bundle="$(basename "${SIDECAR_BUNDLE}" .tgz)"
bundle="${versioned_bundle#cf-usb-sidecar-}"
bundle="${bundle%%-*}"
bundle="cf-usb-sidecar-${bundle}"

cd kubernetes-charts-suse-com
# Remove old charts
rm -rf "stable/${bundle}"
# Place the new ones
mkdir "stable/${bundle}"
tar -x -C "stable/${bundle}" -f ../bundle.tgz
# We can't ship the automatic database provisioning, because we don't track
# license provenance and it's not based on SLE
rm -f "stable/${bundle}/templates/db.yaml"

# Fix up the registry host name to the prod server
sed -i 's@^\(\s\+\)hostname:\s\+".*"$@\1hostname: "registry.suse.com"@' "stable/${bundle}/values.yaml"

$HUB config user.email "cf-ci-bot@suse.de"
$HUB config user.name "${GITHUB_USER}"
$HUB checkout -b "${versioned_bundle}"
$HUB -c core.fileMode=false add .
$HUB commit -m "Submitting ${versioned_bundle}"
$HUB push origin "${versioned_bundle}"

# Open a Pull Request, head: current branch, base: master
$HUB pull-request -m "${versioned_bundle} submitted by ${SOURCE_BUILD:-<unknown build>}" -b master
