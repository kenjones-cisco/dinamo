#!/bin/bash
# Generate test coverage statistics for Go packages.
#
# Works around the fact that `go test -coverprofile` currently does not work
# with multiple packages, see https://github.com/golang/go/issues/6909
#

set -e

workdir=cover
profile="$workdir/cover.out"
mode=count
results=test.out


generate_cover_data() {
    for pkg in $(glide nv);
    do
        for subpkg in $(go list "${pkg}");
        do
            f="$workdir/$(echo "$subpkg" | tr / -).cover"
            go test -v -covermode="$mode" -coverprofile="$f" "$subpkg" >> "$results"
        done
    done

    set -- "$workdir"/*.cover
    if [ ! -f "$1" ]; then
        rm -f "$results" || :
        echo "No Test Cases"; exit 0
    fi
    echo "mode: $mode" >"$profile"
    grep -h -v "^mode:" "$workdir"/*.cover >>"$profile"
}

show_html_report() {
    go tool cover -html="$profile" -o="$workdir"/coverage.html
}

show_ci_report() {
    goveralls -coverprofile="$profile" -service=travis-ci -package github.com/kenjones-cisco/dinamo
}

_done() {
    local error_code="$?"

    # display actual test results
    if [ -f "$results" ]; then
      cat "$results"
    fi

    return $error_code
}

trap "_done" EXIT

rm -f "$results"
generate_cover_data


case "$1" in
"")
    show_html_report ;;
--ci)
    show_ci_report ;;
*)
    echo >&2 "error: invalid option: $1"; exit 1 ;;
esac
