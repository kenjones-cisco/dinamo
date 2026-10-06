# Contributing to dinamo

This is a short guide on how to contribute to the project.

## Submitting a pull request

If you find a bug that you'd like to fix, or a new feature that you'd like to implement then please submit a pull request via Github.

You'll need a Go environment set up with GOPATH set. See [the Go getting started docs](https://golang.org/doc/install) for more info.

Now in your terminal

    git clone ssh://git@github.com/kenjones-cisco/dinamo.git
    cd dinamo

Make a branch to add your new feature

    git checkout -b my-new-feature master

And get hacking.

When ready - run the unit tests for the code you changed

    make test

Make sure you

  * Add documentation for a new feature
  * Add unit tests for a new feature
  * squash commits down to one per feature
  * rebase to master `git rebase master`

When you are done with that

    git push origin my-new-feature

Go to the Github website and click [Create pull request](https://help.github.com/articles/about-pull-requests/).

You patch will get reviewed and you might get asked to fix some stuff.

If so, then make the changes in the same branch, squash the commits, rebase it to master then push it to Github with `--force`.

## Test / Build

Tests are run using a testing framework, so at the top level you can run this to run all the tests.

**Assumes that you have Minishift / Docker Toolbox / Docker for Mac / Docker for Windows installed.**

`make check` applies formatting and available lint fixes using `.golangci.yml`.
`make format` is a compatibility alias for the same checks.
Use `CI=1 make check` to validate formatting and lint without modifying source files.
For `make local` commands, install golangci-lint v2.10.1 locally.

`make test` uses gotestsum v1.13.0 to summarize tests and run benchmarks.
Use `TEST_NAME` for a Go test name regular expression and `TEST_PKG` for a single
package or package pattern; the default is all packages (`./...`). For example:

```bash
make test TEST_PKG=./generator TEST_NAME='^TestGenerate_withSourceData$'
make test-race
```

`make test-race` runs the same tests and benchmarks with race detection and
enables CGO for that invocation. For local testing, install
`gotest.tools/gotestsum@v1.13.0`; race mode also requires a C compiler.

`make cover` runs fresh tests with cross-package coverage and requires at least
80% statement coverage. It writes the native profile to `cover/cover.out`, the
function summary to `cover/coverage.txt`, and an HTML report to
`cover/coverage.html`. Open the HTML report in a browser to inspect source coverage.
Reports are generated before enforcing the threshold, so they remain available
when coverage is below 80%.

GitHub Actions runs `CI=1 make cover` for pushes to `master` and pull requests.
Coverage appears in the job summary, and the `coverage` artifact retains the
profile, summary, and HTML for 14 days. Download the artifact to view its HTML;
artifacts are not hosted webpages. Successful trusted runs publish to Coveralls
using the built-in GitHub token. Fork pull requests run the same validation and
retain reports but do not publish to Coveralls. Configure any hosted coverage
alerts in Coveralls; the local 80% gate is independent of that service.

```bash
# runs all tests (includes formatting and linting)
make test
make local test
# run all tests and generates code coverage (includes formatting and linting)
make cover
make local cover
# builds the default binary (linux amd64); (includes formatting and linting)
make build
make local build
```
