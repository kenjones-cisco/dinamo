#!/bin/bash

set -euo pipefail

if [[ "$#" != 1 ]]; then
    echo >&2 "usage: $0 snapshot|vX.Y.Z"
    exit 1
fi

sha=$(git rev-parse HEAD)
version=$(jq -er '.version' dist/metadata.json)
[[ "$(jq -er '.commit' dist/metadata.json)" == "$sha" ]]
if [[ "$1" == snapshot ]]; then
    [[ "$version" == "0.0.0-dev-$(git rev-parse --short HEAD)" ]]
    display="$version"
else
    [[ "$1" =~ ^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$ ]]
    [[ "$version" == "${1#v}" && "$(jq -er '.tag' dist/metadata.json)" == "$1" ]]
    display="$1"
fi

archives=(dist/*.zip)
[[ "${#archives[@]}" == 6 ]]
[[ "$(wc -l < "dist/dinamo_${version}_SHA256SUMS")" == 6 ]]
(cd dist && sha256sum --check "dinamo_${version}_SHA256SUMS")

for os in linux darwin windows; do
    for arch in amd64 arm64; do
        binary=dinamo
        if [[ "$os" == windows ]]; then binary=dinamo.exe; fi
        archive="dist/dinamo_${version}_${os}_${arch}.zip"
        expected=$(printf '%s\n' CHANGELOG.md LICENSE README.md "$binary" | LC_ALL=C sort)
        [[ "$(unzip -Z1 "$archive" | LC_ALL=C sort)" == "$expected" ]]
        binaries=(dist/dinamo_"${os}"_"${arch}"*/"$binary")
        [[ "${#binaries[@]}" == 1 ]]
        unzip -p "$archive" "$binary" | cmp - "${binaries[0]}"
        metadata=$(go version -m "${binaries[0]}")
        grep -Fq "GOOS=$os" <<< "$metadata"
        grep -Fq "GOARCH=$arch" <<< "$metadata"
        grep -Fq "vcs.revision=$sha" <<< "$metadata"
        if [[ "$1" != snapshot ]]; then
            grep -Fq 'vcs.modified=false' <<< "$metadata"
        fi
    done
done

output=$(dist/dinamo_linux_amd64_v1/dinamo --version 2>&1)
grep -Fxq $' version\t'"$display" <<< "$output"
grep -Fxq $' Git commit\t'"$sha" <<< "$output"
printf 'Verified six ZIPs, archive contents, SHA-256 checksums, and binary metadata (%s, %s)\n' "$display" "$sha"
