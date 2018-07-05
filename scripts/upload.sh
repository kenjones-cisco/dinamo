#!/bin/bash

set -e

PROJECT_FILE="${PROJECT_FILE:-project.yml}"
GITHUB_TOKEN="${GITHUB_TOKEN:?missing required input \'GITHUB_TOKEN\'}"

repo=$(yaml read "${PROJECT_FILE}" metadata.namespace)
project=$(yaml read "${PROJECT_FILE}" metadata.name)
version=$(gobump show -r version)

echo "==> Uploading..."
# upload the new release as a draft
ghr -u "$repo" -r "$project" -draft "$version" release
