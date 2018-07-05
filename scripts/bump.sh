#!/bin/bash

set -e

RELEASE_TYPE="${RELEASE_TYPE:?missing required input \'RELEASE_TYPE\'}"

echo "current version: $(gobump show -r version)"
gobump "$RELEASE_TYPE" -w version
next_version=$(gobump show -r version)
echo "new version: $next_version"

git-chglog --output CHANGELOG.md --next-tag "$next_version" "$next_version"

git add version/info.go CHANGELOG.md
git commit -m "Release $next_version"
git tag -a "$next_version" -m "Version $next_version"
git push && git push --tags
