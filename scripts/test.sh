#!/bin/bash

set -euo pipefail

if [[ "$#" -gt 1 ]]; then
	echo >&2 "usage: $0 [--race]"
	exit 1
fi

options=(-bench .)
case "${1:-}" in
"")
	;;
--race)
	export CGO_ENABLED=1
	options+=(-race)
	;;
*)
	echo >&2 "error: invalid option: $1"
	exit 1
	;;
esac

if [[ -n "${TEST_NAME:-}" ]]; then
	options+=(-run "$TEST_NAME")
fi

packages=("${TEST_PKG:-./...}")
exec gotestsum --no-color=false -- "${packages[@]}" "${options[@]}"
