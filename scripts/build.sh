#!/bin/bash

PROJECT_FILE="${PROJECT_FILE:-project.yml}"


get_targets() {
    local data

    data=$(yaml -j read "${PROJECT_FILE}" metadata.build[*] | jq -r '.[] | [ .["target"] ] | join(" ")')
    echo "$data"
}

build() {
    local project
    local import_path
    local ldflags
    local binary

    project=$(yaml read "${PROJECT_FILE}" metadata.name)
    import_path=$(yaml read "${PROJECT_FILE}" metadata.import)

    ldflags="-X ${import_path}/version.GitCommit=$(git rev-parse --short HEAD)"
    ldflags="${ldflags} -X ${import_path}/version.GitDescribe=$(git describe --tags --always)"

    mapfile -t targets < <(get_targets)
    for target in "${targets[@]}"; do
        if [[ "$target" == "." ]]; then
          binary="$project"
        else
          binary=$(basename "$target")
        fi
        echo "building: $target ==> bin/$binary"
        go build -ldflags "${ldflags}" -o "bin/$binary" "${import_path}/$target"
    done
}

build
