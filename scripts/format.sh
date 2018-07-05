#!/bin/bash

PROJECT_FILE="${PROJECT_FILE:-project.yml}"
local_import=$(yaml read "${PROJECT_FILE}" metadata.import)

find . \( -path ./vendor -o -path ./.glide \) -prune -o -name "*.go" -exec goimports -local "${local_import}" -w {} \;
