#!/bin/sh

set -e

TOPDIR=$(cd "$(dirname "$0")/.." && pwd)

workdir=.cover
profile="$workdir/cover.out"
mode=count

# clean up the artifacts for coverage.
rm -rf $workdir

generate_cover_data() {
    mkdir -p "$workdir"

    for pkg in "$@"; do
		echo $pkg
        f="$workdir/$(echo $pkg | tr / -).cover"
		echo $f
        GOPATH="${TOPDIR}:${TOPDIR}/v:${GOPATH}" godep go test \
	    -covermode="$mode" \
	    -coverprofile="$f" "$pkg"
    done

    echo "mode: $mode" >"$profile"
    grep -h -v "^mode:" "$workdir"/*.cover >>"$profile"
}

generate_cover_data $(go list ./cmd/... ./src/... | grep -v github.com/go-swagger)
go tool cover -func="$profile"
go tool cover -html="$profile"
