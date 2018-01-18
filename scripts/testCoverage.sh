#!/bin/sh

set -o errexit
set -o nounset

TOPDIR=$(cd "$(dirname "$0")/.." && pwd)

workdir=.cover
profile="${workdir}/cover.out"
mode=count

# clean up the artifacts for coverage.
rm -rf "${workdir}"

# Manipulate the GOPATH to include both go-swagger and go-openapi from
# the sub-moduled swagger. Unfortunately, go(1) gets _really_ confused when
# symlinks are involved; copying the tree is more likely to do the correct thing
swaggerdir="$(mktemp -d)"
trap "rm -rf '${swaggerdir}'" EXIT
cp -r "${TOPDIR}/src/github.com/go-swagger/go-swagger/vendor" "${swaggerdir}/src"
export GOPATH="${swaggerdir}:${GOPATH}"

generate_cover_data() {
    mkdir -p "${workdir}"

    for pkg in "$@"; do
        f="${workdir}/$(echo "${pkg}" | tr / -).cover"
        printf "%s => %s\n" "${pkg}" "${f}"
        go test -covermode="${mode}" -coverprofile="${f}" "${pkg}"
    done

    echo "mode: ${mode}" >"${profile}"
    grep -h -v "^mode:" "${workdir}"/*.cover >>"${profile}"
}

generate_cover_data $(go list ./cmd/... ./src/... | grep -v github.com/go-swagger)
go tool cover -func="${profile}"
go tool cover -html="${profile}"
