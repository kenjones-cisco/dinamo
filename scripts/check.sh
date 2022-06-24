#!/bin/bash

set -o errexit
set -o pipefail

golangci-lint -v run --color=always --print-resources-usage -c .golangci.yml --fix
