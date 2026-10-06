#!/bin/bash

set -euo pipefail

verify_release() {
    local tag="$1" sha="$2" release checkout_sha release_tag_sha source_status source_version

    if [[ ! "$tag" =~ ^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$ || ! "$sha" =~ ^[0-9a-f]{40}$ ]]; then
        echo >&2 "error: an exact vX.Y.Z tag and full commit SHA are required"
        return 1
    fi
    checkout_sha=$(git rev-parse HEAD) || return 1
    release_tag_sha=$(git rev-parse "refs/tags/$tag^{commit}") || return 1
    if [[ "$checkout_sha" != "$sha" || "$release_tag_sha" != "$sha" ]]; then
        echo >&2 "error: checkout and existing release tag must match the expected SHA"
        return 1
    fi
    if ! git merge-base --is-ancestor "$sha" refs/remotes/origin/master; then
        echo >&2 "error: release commit must belong to origin/master"
        return 1
    fi
    source_status=$(git status --porcelain) || return 1
    if [[ -n "$source_status" ]]; then
        echo >&2 "error: release source must be clean"
        return 1
    fi
    source_version=$(jq -er '.["."]' .release-please-manifest.json) || return 1
    if [[ "$source_version" != "${tag#v}" ]]; then
        echo >&2 "error: release manifest must match the tag"
        return 1
    fi
    release=$(gh api "repos/kenjones-cisco/dinamo/releases/tags/$tag") || return 1
    if ! jq -e --arg tag "$tag" --arg sha "$sha" \
        '.draft == true and .tag_name == $tag and .target_commitish == $sha' <<< "$release" >/dev/null; then
        echo >&2 "error: only an existing draft targeting the exact SHA may be packaged"
        return 1
    fi
}

release_main() {
    if [[ "$#" != 3 || ( "$1" != check && "$1" != publish ) ]]; then
        echo >&2 "usage: $0 check|publish vX.Y.Z FULL_SHA"
        return 1
    fi
    verify_release "$2" "$3" || return 1
    if [[ "$1" == publish ]]; then
        GORELEASER_CURRENT_TAG="$2" goreleaser release --clean
    fi
}

if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
    release_main "$@"
fi
