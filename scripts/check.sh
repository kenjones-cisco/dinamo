#!/bin/bash

set -o errexit
set -o pipefail

RUN_CMD=(gometalinter)
RESULTS=/dev/stdout

set_opts() {

read -r COMMON_OPTS << EOM
    --exclude=vendor \
    --tests \
    --vendor
EOM

read -r SIMPLE_OPTS << EOM
    ${COMMON_OPTS} \
    --disable=maligned \
    --disable=gotype \
    --disable=gotypex \
    --disable=structcheck \
    --disable=varcheck \
    --disable=interfacer \
    --disable=unconvert \
    --disable=dupl \
    --cyclo-over=15 \
    --deadline=60s \
    ./...
EOM

# The following checks get into lower level detailed analysis of the code, but they
# all are common in that they scan the project differently then the others above that
# accept a simple base path recursion.
read -r DETAILED_OPTS << EOM
    ${COMMON_OPTS} \
    --deadline=60s \
    --disable-all \
    --enable=structcheck \
    --enable=varcheck \
    --enable=interfacer \
    --enable=unconvert
EOM

}

# Excludes:
#   - when using defer there is no way to check to returned value so ignore
#   - some generated code has output parameters named as err that result in vetshadow issue so ignore

# The --exclude statements get passed directly to avoid bash interpolation or escaping of the single quotes
# that results in gometalinter ignore the exclude lines.
check() {
    echo "==> Simple Check..."

    "${RUN_CMD[@]}" \
        --exclude='error return value not checked.*(Close|Log|Print|Shutdown|Unsetenv).*\(errcheck\)$' \
        ${SIMPLE_OPTS} > "${RESULTS}"

    echo "==> Detailed Check..."

    "${RUN_CMD[@]}" ${DETAILED_OPTS} > "${RESULTS}"
}

case "$1" in
    --ci)
        RUN_CMD=(gometalinter --checkstyle)
        # if simple fails then detailed never runs meaning file not overwritten
        # but if simple passes then file is essentially empty so if detailed
        # fails then real content is written to the file else it will be an
        # essentially empty file indicating all checks passed.
        RESULTS=report.xml
        rm -f "${RESULTS}"
        ;;
    *)
        ;;
esac

set_opts
check
