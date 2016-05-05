#!/bin/sh

set -e

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
        godep go test -covermode="$mode" -coverprofile="$f" "$pkg"
    done

    echo "mode: $mode" >"$profile"
    grep -h -v "^mode:" "$workdir"/*.cover >>"$profile"
}

generate_cover_data $(go list ./... | 
grep -v github.com/hpcloud/catalog-service-manager/generated | 
grep -v github.com/hpcloud/catalog-service-manager/example | 
grep -v github.com/hpcloud/catalog-service-manager/csm-extensions |
grep -v github.com/hpcloud/catalog-service-manager/scripts | 
grep -v github.com/hpcloud/catalog-service-manager/cmd/catalog-service-manager/handlers | 
grep -v github.com/hpcloud/catalog-service-manager/src/api)
go tool cover -func="$profile"
go tool cover -html="$profile"
