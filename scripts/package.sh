#!/bin/bash

ROOT_DIR=$(pwd)
VERSION=$(gobump show -r version)

PROJECT_FILE="${PROJECT_FILE:-project.yml}"
PROJECT=$(yaml read "${PROJECT_FILE}" metadata.name)

# Zip and copy to the release dir
echo "==> Packaging..."
mapfile -t bplatforms < <(find ./build -mindepth 1 -maxdepth 1 -type d)
for PLATFORM in "${bplatforms[@]}"; do
    OSARCH=$(basename "${PLATFORM}")
    echo "--> ${OSARCH}"

    pushd "${PLATFORM}" >/dev/null 2>&1
    zip "${ROOT_DIR}/release/${OSARCH}.zip" ./*
    popd >/dev/null 2>&1
done

pushd ./release >/dev/null 2>&1
shasum -a256 ./* > ./"${PROJECT}_${VERSION}_SHA256SUMS"
popd >/dev/null 2>&1
