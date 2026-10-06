MAKEFLAGS += -r --warn-undefined-variables
SHELL := /bin/bash
.SHELLFLAGS := -o pipefail -euc
.DEFAULT_GOAL := help

include Makefile.variables
include Makefile.local

.PHONY: help clean veryclean build vendor dep-* format check test test-race cover docs adhoc xcompile package

## display this help message
help:
	@echo 'Management commands for dinamo:'
	@echo
	@echo 'Usage:'
	@echo '  ## Build Commands'
	@echo '    build           Compile the project.'
	@echo '    xcompile        Compile the project for multiple OS and Architectures.'
	@echo
	@echo '  ## Develop / Test Commands'
	@echo '    vendor          Install dependencies using dep if Gopkg.toml changed.'
	@echo '    dep-update      Update dependencies using dep.'
	@echo '    dep-add         Add new dependencies to dep and install.'
	@echo '    format          Alias for check (format and lint).'
	@echo '    check           Format and lint code; validate only in CI.'
	@echo '    test            Run tests on project.'
	@echo '    test-race       Run tests with the race detector.'
	@echo '    cover           Run tests and capture code coverage metrics on project.'
	@echo '    clean           Clean the directory tree of produced artifacts.'
	@echo '    veryclean       Same as clean but also removes cached dependencies.'
	@echo
	@echo '  ## Release Commands'
	@echo '    package         Preview ZIP archives and checksums in dist/; never publishes.'
	@echo '    Releases        Review release PRs and manually publish drafts; see docs/releasing.md.'
	@echo
	@echo '  ## Local Commands'
	@echo '    setup           Configures Minishfit/Docker directory mounts.'
	@echo '    drma            Removes all stopped containers.'
	@echo '    drmia           Removes all unlabelled images.'
	@echo '    drmvu           Removes all unused container volumes.'
	@echo

.ci-clean:
ifeq ($(CI_ENABLED),1)
	@rm -f tmp/dev_image_id || :
endif

## Clean the directory tree of produced artifacts.
clean: .ci-clean prepare
	@${DOCKERRUN} bash -c 'rm -rf bin build dist release cover *.out *.xml'

## Same as clean but also removes cached dependencies.
veryclean: clean
	@${DOCKERRUN} bash -c 'rm -rf tmp .mod'

## builds the dev container
prepare: tmp/dev_image_id
tmp/dev_image_id: Dockerfile.dev
	@mkdir -p tmp
	@docker rmi -f ${DEV_IMAGE} > /dev/null 2>&1 || true
	@echo "## Building dev container"
	@docker build --quiet -t ${DEV_IMAGE} -f Dockerfile.dev .
	@docker inspect -f "{{ .ID }}" ${DEV_IMAGE} > tmp/dev_image_id

# ----------------------------------------------
# build

## Compile the project.
build: build/dev

build/dev: check */*.go
	@rm -rf bin/
	@mkdir -p bin
	${DOCKERRUN} bash ./scripts/build.sh
	@chmod 755 bin/* || :

## Compile the project for multiple OS and Architectures.
xcompile: check
	${DOCKERRUN} goreleaser build --snapshot --clean

# ----------------------------------------------
# dependencies

## Install dependencies using go mod if go.mod changed.
vendor: tmp/vendor-installed
tmp/vendor-installed: tmp/dev_image_id go.mod
	@mkdir -p .mod
	${DOCKERRUN} go mod tidy
	@date > tmp/vendor-installed
	@chmod 644 go.sum || :

# ----------------------------------------------
# develop and test

## print environment info about this dev environment
debug:
	@echo IMPORT_PATH="$(IMPORT_PATH)"
	@echo ROOT="$(ROOT)"
	@echo
	@echo docker commands run as:
	@echo "$(DOCKERRUN)"

## Alias for check (format and lint).
format: check

## Format and lint code; validate only in CI.
check: tmp/vendor-installed
ifeq ($(CI_ENABLED),1)
	${DOCKERNOVENDOR} bash ./scripts/check.sh --check
else
	${DOCKERNOVENDOR} bash ./scripts/check.sh
endif

## Run tests on project.
test: check
	${DOCKERRUN} bash ./scripts/test.sh

## Run tests with the race detector.
test-race: check
	${DOCKERRUN} bash ./scripts/test.sh --race

## Run tests and capture code coverage metrics on project.
cover: check
	@rm -rf cover/
	@mkdir -p cover
	${DOCKERRUN} bash ./scripts/cover.sh

docs: prepare
	@rm -f docs/dinamo*.md
	@mkdir -p docs
	${DOCKERNOVENDOR} go run gendocs.go
	@chmod 755 docs
	@chmod 644 docs/*.md

# usage: make adhoc RUNTHIS='command to run inside of dev container'
# example: make adhoc RUNTHIS='which jq'
adhoc: prepare
	@${DOCKERRUN} ${RUNTHIS}

# ----------------------------------------------
# release

package: check
	${DOCKERRUN} goreleaser check
	${DOCKERRUN} goreleaser release --snapshot --clean
