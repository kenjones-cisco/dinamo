#!/bin/bash

set -euo pipefail

if [[ "$#" -ne 0 ]]; then
    echo >&2 "usage: $0"
    exit 1
fi

mkdir -p cover
rm -f cover/cover.out cover/coverage.txt cover/coverage.html

go test -count=1 -covermode=atomic -coverpkg=./... -coverprofile=cover/cover.out ./...
go tool cover -func=cover/cover.out | tee cover/coverage.txt
go tool cover -html=cover/cover.out -o=cover/coverage.html

awk '
NR > 1 {
    blockStatements[$1] = $2
    blockExecutions[$1] += $3
}
END {
    for (block in blockStatements) {
        statements += blockStatements[block]
        if (blockExecutions[block] > 0) {
            covered += blockStatements[block]
        }
    }
    if (statements == 0) {
        print "error: coverage profile contains no statements" > "/dev/stderr"
        exit 1
    }
    coverage = 100 * covered / statements
    printf "Statement coverage: %.2f%% (minimum: 80%%)\n", coverage
    if (coverage < 80) {
        print "error: statement coverage is below 80%" > "/dev/stderr"
        exit 1
    }
}' cover/cover.out
