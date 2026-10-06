#!/bin/bash

set -euo pipefail

source scripts/release.sh

sha=0123456789012345678901234567890123456789
head_sha="$sha"
tag_sha="$sha"
target_sha="$sha"
draft=true
clean=true
ancestor=true
manifest_version=0.4.0
api_ok=true
git_ok=true
tag_exists=true
status_ok=true
manifest_ok=true
upload_ok=true
uploads=0

git() {
    [[ "$git_ok" == true ]] || return 1
    case "$1" in
        rev-parse)
            if [[ "$2" == HEAD ]]; then
                printf '%s\n' "$head_sha"
            else
                [[ "$tag_exists" == true ]] || return 1
                printf '%s\n' "$tag_sha"
            fi
            ;;
        merge-base) [[ "$ancestor" == true ]] ;;
        status)
            [[ "$status_ok" == true ]] || return 1
            if [[ "$clean" != true ]]; then echo ' M version/info.go'; fi
            ;;
        *) echo >&2 "unexpected git operation: $*"; return 1 ;;
    esac
}

gh() {
    [[ "$api_ok" == true ]] || return 1
    [[ "$*" == 'api repos/kenjones-cisco/dinamo/releases/tags/v0.4.0' ]] || return 1
    printf '{"draft":%s,"tag_name":"v0.4.0","target_commitish":"%s"}\n' "$draft" "$target_sha"
}

jq() {
    if [[ "$*" == '-er .["."] .release-please-manifest.json' ]]; then
        [[ "$manifest_ok" == true ]] || return 1
        printf '%s\n' "$manifest_version"
    else
        command jq "$@"
    fi
}

goreleaser() {
    [[ "$*" == 'release --clean' && "$GORELEASER_CURRENT_TAG" == v0.4.0 ]] || return 1
    uploads=$((uploads + 1))
    [[ "$upload_ok" == true ]]
}

reject() {
    if release_main publish "$1" "$2"; then
        echo >&2 "error: unsafe release was accepted"
        exit 1
    fi
    [[ "$uploads" == 4 ]]
}

release_main check v0.4.0 "$sha"
[[ "$uploads" == 0 ]]
release_main publish v0.4.0 "$sha"
release_main publish v0.4.0 "$sha"
[[ "$uploads" == 2 ]]
upload_ok=false
if release_main publish v0.4.0 "$sha"; then
    echo >&2 'error: failed upload was reported as successful'
    exit 1
fi
[[ "$uploads" == 3 ]]
upload_ok=true
release_main publish v0.4.0 "$sha"
[[ "$uploads" == 4 ]]
reject 0.3.0 "$sha"
reject v0.4.0 master
head_sha=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
reject v0.4.0 "$sha"
head_sha="$sha"
tag_sha=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
reject v0.4.0 "$sha"
tag_sha="$sha"
ancestor=false
reject v0.4.0 "$sha"
ancestor=true
clean=false
reject v0.4.0 "$sha"
clean=true
manifest_version=0.3.0
reject v0.4.0 "$sha"
manifest_version=0.4.0
draft=false
reject v0.4.0 "$sha"
draft=true
target_sha=aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
reject v0.4.0 "$sha"
target_sha="$sha"
api_ok=false
reject v0.4.0 "$sha"
api_ok=true
git_ok=false
reject v0.4.0 "$sha"
git_ok=true
tag_exists=false
reject v0.4.0 "$sha"
tag_exists=true
status_ok=false
reject v0.4.0 "$sha"
status_ok=true
manifest_ok=false
reject v0.4.0 "$sha"
echo 'Release guards: validation, failed upload/retry, draft retry, and 14 rejection cases passed'
