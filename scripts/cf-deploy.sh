#! /usr/bin/env bash

# This script deploys the given service on the default kubernetes context

set -o errexit -o nounset
service="$(tr '[:upper:]' '[:lower:]' <<<"${1:-mysql}")"
host="${service^^}"
extra=""
case "${service}" in
    mysql)
        ;;
    postgres)
        host="POSTGRESQL"
        extra="--set env.SERVICE_POSTGRESQL_SSLMODE=disable"
        ;;
    *)
        printf "Unknown service %s\n" "${service}" >&2
        exit 1
        ;;
esac

deployment_name="$(helm list | grep 'DEPLOYED' | awk '$NF == "cf" { print $1 }' | tail -n1)"

get_value() {
    helm get values "${deployment_name}" | y2j | jq -r "${1}"
}

get_secret() {
    kubectl get secret -n cf secret -o jsonpath="{@${1}}" | base64 -d
}

helm list --all | \
    awk "\$NF == \"dev-${service}\" { print \$1 }" | \
    xargs --no-run-if-empty helm delete --purge

helm install \
    "$(dirname "${0}")/../csm-extensions/services/dev-${service}/output/helm" \
    --name "dev-${service}" \
    --namespace "dev-${service}" \
    --wait \
    --timeout 300 \
    --set env.CF_ADMIN_PASSWORD="$(get_secret .data.cluster-admin-password)" \
    --set env.CF_ADMIN_USER=admin \
    --set env.CF_CA_CERT="$(get_secret .data.internal-ca-cert)" \
    --set env.CF_DOMAIN="$(get_value .env.DOMAIN)" \
    --set env.SERVICE_LOCATION="http://cf-usb-sidecar-${service}.dev-${service}.svc.cluster.local:8081" \
    --set env.UAA_CA_CERT="$(get_secret .data.uaa-ca-cert)" \
    --set env.SERVICE_${host}_HOST=AUTO \
    --set kube.registry.hostname="${DOCKER_REPOSITORY:-docker.io}" \
    --set kube.organization="${DOCKER_ORGANIZATION:-splatform}" \
    ${extra}
