#!/bin/sh

. "${SIDECAR_ROOT}/scripts/colors.sh"

force_rebuild="${1:-}"

current_dir="$(pwd)"
script_dir="$(dirname "$0")"

if ! [ -e "${current_dir}/.git" ]
then
	printf "%b==> Please execute script from catalog-service-manager's project root directory%b\n" "${ERROR_COLOR}" "${NO_COLOR}"
	exit 1
fi

if [ "${script_dir}" != "scripts/docker" ]
then
	printf "%b==> Script directory is not correct, please run script from \${PROJECT_ROOT}/scripts/dockerv%b\n" "${ERROR_COLOR}" "${NO_COLOR}"
	exit 1	
fi

if docker images | grep "${SIDECAR_BUILD_BASE_IMAGE_NAME}" | grep --quiet "${SIDECAR_BUILD_BASE_IMAGE_TAG}" > /dev/null 2>&1
then
	if [ "$force_rebuild" != "rebuild-image" ]
	then
		printf "%b==> ${SIDECAR_BUILD_BASE_IMAGE_NAME}:${SIDECAR_BUILD_BASE_IMAGE_TAG} image already exists!%b\n" "${WARN_MAGENTA}" "${NO_COLOR}"
		exit 0
	fi
	
	if [ "$force_rebuild" = "rebuild-image" ]
	then
		printf "%b==> Removing old ${SIDECAR_BUILD_BASE_IMAGE_NAME}:${SIDECAR_BUILD_BASE_IMAGE_TAG} images ...%b\n" "${OK_GREEN_COLOR}" "${NO_COLOR}"
		docker images | grep "${SIDECAR_BUILD_BASE_IMAGE_NAME}" | grep "${SIDECAR_BUILD_BASE_IMAGE_TAG}" | awk '{print $3}' | xargs -n 1 docker rmi -f > /dev/null 2>&1
		sleep 5
	fi	
fi

printf "%b==> Building ${SIDECAR_BUILD_BASE_IMAGE_NAME}:${SIDECAR_BUILD_BASE_IMAGE_TAG} image ...%b\n" "${OK_GREEN_COLOR}" "${NO_COLOR}"
docker build \
    --tag "${SIDECAR_BUILD_BASE_IMAGE_NAME}:${SIDECAR_BUILD_BASE_IMAGE_TAG}" \
    --rm \
    --file "${script_dir}/Dockerfile-build" .
