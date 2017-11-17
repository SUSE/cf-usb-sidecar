#!/bin/bash
##
# This script installs certificate authorities:
# - The SCF CA securing the TLS cert securing cf <-> API comms

# The file is derived from
# scf.git:container-host-files/etc/hcf/config/scripts/authorize_internal_ca.sh
# and modified to suit. As it is sourced by `scf-connector.sh` as an
# environment script we do things to ensure that we have an acceptable
# environment.

# And snarfed from
# https://github.com/SUSE/fissile/blob/master/scripts/dockerfiles/run.sh#L32
# a supporting function to determine the OS we are running on.
# Notes:
# - Removed `$chroot` reference. Unclear where it came from in the
#   original script, and not needed here.
# - Added /etc/os-release and detection of debian, that is the base OS
#   for the sidecar images, at the moment.

function get_os_type {
    centos_file=/etc/centos-release
    rhel_file=/etc/redhat-release
    ubuntu_file=/etc/lsb-release
    photonos_file=/etc/photon-release
    opensuse_file=/etc/SuSE-release
    general_file=/etc/os-release

    os_type=''
    if [ -f $photonos_file ]
    then
	os_type='photonos'
    elif [ -f $ubuntu_file ]
    then
	os_type='ubuntu'
    elif [ -f $centos_file ]
    then
	os_type='centos'
    elif [ -f $rhel_file ]
    then
	os_type='rhel'
    elif [ -f $opensuse_file ]
    then
	os_type='opensuse'
    elif [ -f $general_file ]
    then
	if grep -qsi 'debian' $general_file
	then
	    os_type='debian'
	fi
    fi

    echo $os_type
}

os_type=$(get_os_type)
if [ "$os_type" == "ubuntu" ]; then
    ca_path=/usr/local/share/ca-certificates
elif [ "$os_type" == "debian" ]; then
    ca_path=/usr/local/share/ca-certificates
elif [ "$os_type" == "opensuse" ]; then
    ca_path=/etc/pki/trust/anchors
else
    printf "Error: unknown operating system '${os_type}'"
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
