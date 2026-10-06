#!/bin/bash

set -euo pipefail

if [[ "$#" -gt 1 ]]; then
	echo >&2 "usage: $0 [--fix|--check]"
	exit 1
fi

case "${1:-}" in
""|--fix|--check)
	;;
*)
	echo >&2 "error: invalid option: $1"
	exit 1
	;;
esac

opts=(--color=always --config=.golangci.yml)
if [[ -n "${DEBUG:-}" ]]; then
	opts+=(--verbose)
fi

if [[ "${1:-}" == "--check" ]]; then
	golangci-lint fmt "${opts[@]}" --diff
	golangci-lint run "${opts[@]}"
else
	golangci-lint fmt "${opts[@]}"
	golangci-lint run "${opts[@]}" --fix
fi
