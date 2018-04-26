#!/bin/bash
# Connect the sidecar to SCF.
# Assumes to have `cf` with `cf-usb-plugin` installed.

set -o errexit
set -o nounset
#set -o xtrace

# Parameters, via Environment
# - CF_ADMIN_PASSWORD   (SCF cluster admin password)
# - CF_ADMIN_USER       (SCF cluster admin)
# - CF_CA_CERT          (PEM-encoded CA cert securing the TLS cert securing cf API comms)
# - CF_DOMAIN           (SCF base domain)
# - SERVICE_LOCATION    (Where SCF will find our service csm)
# - SERVICE_TYPE        (What SCF will use as its name for our service)
# - SIDECAR_API_KEY     (generated secret)
# - UAA_CA_CERT         (PEM-encoded CA cert securing the TLS cert securing uaa API comms)

# Default
SERVICE_TYPE="${SERVICE_TYPE:-mysql}"
BIN_DIR="$(cd $(dirname $0) && pwd)"

# Report progress to the user; use as printf
status() {
    local fmt="${1}"
    shift
    printf "\n%b${fmt}%b\n" "\033[0;32m" "$@" "\033[0m"
}

# Report problem to the user; use as printf
trouble() {
    local fmt="${1}"
    shift
    printf "\n%b${fmt}%b\n" "\033[0;31m" "$@" "\033[0m"
}

# helper function to retry a command several times, with a delay between trials
# usage: retry <max-tries> <delay> <command>...
function retry () {
    max=${1}
    delay=${2}
    i=0
    shift 2

    while test ${i} -lt ${max} ; do
        printf "Trying: %s\n" "$*"
        if "$@" ; then
            status ' SUCCESS'
            break
        fi
        trouble '  FAILED'
        status "Waiting ${delay} ..."
        sleep "${delay}"
        i="$(expr ${i} + 1)"
    done
}

status "Installing the CC CA certificate ..."
. "${BIN_DIR}/authorize_ca.sh"

cf install-plugin -f /usr/local/bin/cf-usb-plugin

status "Waiting for CC ..."
retry 240 30s cf api "api.${CF_DOMAIN}"

status "Logging in"
cf auth ${CF_ADMIN_USER} ${CF_ADMIN_PASSWORD}

status "Registering MySQL sidecar with CC"
cf usb create-driver-endpoint \
    "${SERVICE_TYPE}" \
    "${SERVICE_LOCATION}" \
    "${SIDECAR_API_KEY}" \
    -c ":"

status "MySQL sidevar configuration complete."
exit 0
