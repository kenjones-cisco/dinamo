MAKEFLAGS += -r --warn-undefined-variables
SHELL := /bin/bash
.SHELLFLAGS := -o pipefail -euc
.DEFAULT_GOAL := help

include Makefile.variables
include Makefile.local

.PHONY: help clean veryclean build vendor dep-* format check test cover docs adhoc xcompile bump package upload release

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
	@echo '    vendor          Install dependencies using glide if glide.yaml changed.'
	@echo '    dep-update      Update dependencies using glide.'
	@echo '    dep-add         Add new dependencies to glide and install.'
	@echo '    format          Run code formatter.'
	@echo '    check           Run static code analysis (lint).'
	@echo '    test            Run tests on project.'
	@echo '    cover           Run tests and capture code coverage metrics on project.'
	@echo '    clean           Clean the directory tree of produced artifacts.'
	@echo '    veryclean       Same as clean but also removes cached dependencies.'
	@echo
	@echo '  ## Release Commands'
	@echo '    bump            Increment version, generate changelog, tag project and push to Github.'
	@echo '    package         Create archives for each binary and checksum.'
	@echo '    upload          Uploaded archives and create release on Github.'
	@echo '    release         Orchestrates bump, xcompile, package and upload tasks.'
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
	@${DOCKERRUN} bash -c 'rm -rf bin build release cover *.out *.xml'

## Same as clean but also removes cached dependencies.
veryclean: clean
	@${DOCKERRUN} bash -c 'rm -rf tmp .glide vendor'

## builds the dev container
prepare: tmp/dev_image_id
tmp/dev_image_id: Dockerfile.dev
	@mkdir -p tmp
	@docker rmi -f ${DEV_IMAGE} > /dev/null 2>&1 || true
	@echo "## Building dev container"
	@docker build --quiet -t ${DEV_IMAGE} --build-arg DEVELOPER="${DEVELOPER}" -f Dockerfile.dev .
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
	@rm -rf build/
	@mkdir -p build
	${DOCKERRUN} bash ./scripts/xcompile.sh
	@find build -type d -exec chmod 755 {} \; || :
	@find build -type f -exec chmod 755 {} \; || :

# ----------------------------------------------
# dependencies

## Install dependencies using glide if glide.yaml changed.
vendor: tmp/glide-installed
tmp/glide-installed: tmp/dev_image_id glide.yaml
	@mkdir -p vendor
	${DOCKERRUN} glide install --skip-test --strip-vendor
	@date > tmp/glide-installed
	@chmod 644 glide.lock || :

## Update dependencies using glide.
dep-update: prepare
	${DOCKERRUN} glide up --skip-test --strip-vendor
	@chmod 644 glide.lock || :

# usage DEP=github.com/owner/package make dep-add
## Add new dependencies to glide and install.
dep-add: prepare
ifeq ($(strip $(DEP)),)
	$(error "No dependency provided. Expected: DEP=<go import path>")
endif
	${DOCKERNOVENDOR} glide get --skip-test --strip-vendor ${DEP}
	@chmod 644 glide.lock || :

# ----------------------------------------------
# develop and test

## print environment info about this dev environment
debug:
	@echo IMPORT_PATH="$(IMPORT_PATH)"
	@echo ROOT="$(ROOT)"
	@echo RELEASE_TYPE="$(RELEASE_TYPE)"
	@echo
	@echo docker commands run as:
	@echo "$(DOCKERRUN)"

## Run code formatter.
format: tmp/glide-installed
	${DOCKERNOVENDOR} bash ./scripts/format.sh
	@if [[ -n "$$(git -c core.fileMode=false status --porcelain)" ]]; then \
    	echo "goimports modified code; requires attention!" ; \
    	if [[ "${CI_ENABLED}" == "1" ]]; then \
        	exit 1 ; \
    	fi ; \
	fi

## Run static code analysis (lint).
check: format
ifeq ($(CI_ENABLED),1)
	${DOCKERNOVENDOR} bash ./scripts/check.sh --ci
else
	${DOCKERNOVENDOR} bash ./scripts/check.sh
endif

## Run tests on project.
test: check
	${DOCKERRUN} bash ./scripts/test.sh

## Run tests and capture code coverage metrics on project.
cover: check
	@rm -rf cover/
	@mkdir -p cover
ifeq ($(CI_ENABLED),1)
	${DOCKERRUN} bash ./scripts/cover.sh --ci
else
	${DOCKERRUN} bash ./scripts/cover.sh
	@chmod 644 cover/coverage.html
endif

docs: prepare
	@rm -rf docs/
	@mkdir -p docs/usage
	${DOCKERNOVENDOR} bash ./scripts/docs.sh
	@chmod 755 docs
	@chmod 755 docs/usage
	@chmod 644 docs/usage/*.md

# usage: make adhoc RUNTHIS='command to run inside of dev container'
# example: make adhoc RUNTHIS='which glide'
adhoc: prepare
	@${DOCKERRUN} ${RUNTHIS}

# ----------------------------------------------
# release

bump: prepare
	${DOCKERNOVENDOR} bash ./scripts/bump.sh

package: xcompile
	@rm -rf release/
	@mkdir -p release/
	${DOCKERNOVENDOR} bash ./scripts/package.sh

upload: prepare
	${DOCKERNOVENDOR} bash ./scripts/upload.sh

release: bump package upload
