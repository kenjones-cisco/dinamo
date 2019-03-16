#!/bin/bash

set -o errexit
set -o pipefail

golangci-lint -v run --print-resources-usage -c .golangci.yml
