#!/bin/bash

set -o errexit
set -o pipefail

golangci-lint run
