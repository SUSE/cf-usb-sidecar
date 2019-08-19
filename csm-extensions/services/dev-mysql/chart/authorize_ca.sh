#!/bin/bash
##
# This script installs certificate authorities:
# - The SCF CA securing the TLS cert securing cf <-> API comms

# The file is derived from
# scf.git:container-host-files/etc/scf/config/scripts/authorize_internal_ca.sh
# and modified to suit. As it is sourced by `scf-connector.sh` as an
# environment script we do things to ensure that we have an acceptable
# environment.

for ca_path in $(echo '
    /etc/pki/trust/anchors
    /usr/local/share/ca-certificates
') ; do
    if test -d "${ca_path}" ; then
        break
    fi
done
if ! test -d "${ca_path}" ; then
    printf "Error: Unable to find local CA certificate directory\n" >&2
    exit 1
fi

# We are sourced from a different script
if ! ( echo "${SHELLOPTS:-}" | tr ':' '\n' | grep --quiet errexit ) ; then
    printf "Error: errexit not set\n" >&2
    exit 1
fi

# SCF CA cert
if [ -r /etc/secrets/cf-ca-cert ]; then
    cp /etc/secrets/cf-ca-cert "${ca_path}"/cf-CA.crt
elif [ -n "${CF_CA_CERT:-}" ]; then
    printf "%b" "${CF_CA_CERT}" > "${ca_path}"/cf-CA.crt
fi

# UAA CA cert
if [ -r /etc/secrets/uaa-ca-cert ]; then
    cp /etc/secrets/uaa-ca-cert "${ca_path}"/uaa-CA.crt
elif [ -n "${UAA_CA_CERT:-}" ]; then
    printf "%b" "${UAA_CA_CERT}" > "${ca_path}"/uaa-CA.crt
fi

update-ca-certificates
