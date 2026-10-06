# Release Runbook

## Approved Design

- Release Please owns version decisions, release PRs, and changelog updates.
- Merging a reviewed release PR prepares a draft; a human publishes it.
- Bootstrap from historical release `0.3.0`; the first new release is `v0.4.0`.
- Future tags use `vX.Y.Z`. Historical tags and releases remain unchanged.
- Before 1.0, fixes increment patch; features and breaking changes increment minor.
- Build Linux, macOS (darwin), and Windows for both amd64 and arm64.
- Use ZIP for all six builds. Preserve `dinamo_<version>_<os>_<arch>.zip`
  and `dinamo_<version>_SHA256SUMS`; filename versions omit the tag's `v`.
- Local builds identify as development builds. Release metadata comes from the
  exact release tag and full commit SHA, not a committed source version.
- Use a repository-scoped GitHub App for release PRs so their events trigger CI.
  A separate publisher uses only the required GitHub token permissions.
- Approval consists of release PR review and manual draft publication, without
  an environment approval gate.
- Keep this public repository independent of cicd-common.

## Implementation Plan

1. Verify pinned tool capabilities, bootstrap tracking at the `0.3.0` commit,
   and preserve historical changelog content.
2. Orchestrate draft packaging explicitly from Release Please's tag and SHA
   outputs. Ensure the tag exists at that exact SHA before packaging.
3. Validate the exact release commit without formatting or dependency rewrites.
   Use GoReleaser OSS for binaries, ZIP archives, checksums, and uploads.
4. Restrict retries to matching unpublished drafts. Never rewrite a tag or
   modify published releases or their assets.
5. Replace legacy release tools after tracing all consumers; retain useful
   non-publishing local build and package previews.
6. Document App setup, repository protections, first release, and recovery.
7. Run lint, tests, race tests, the 80% coverage gate, workflow/config validation,
   snapshot builds, archive/checksum/metadata checks, and offline retry tests.
   Report human-only setup and hosted checks requiring a real workflow run.

## Ownership And Timing

The manifest is Release Please's version ledger, initialized to `0.3.0` at
`2706b8335015043ea36b354b61d77a24cab9dfd5`. It is not a binary version source.
Release Please prepends new changelog entries; historical content stays intact.
The Go strategy updates only the changelog and manifest, not `version/info.go`.

On a push to `master`, the repository-scoped App maintains the release PR. After
that PR is merged, Release Please creates the exact tag and a draft at the merge
SHA. `force-tag-creation` is enabled: GitHub otherwise defers creating a draft's
tag until publication. The pinned implementation creates the ref before the
draft. The workflow uses `release_created`, `tag_name`, and `sha` outputs to
explicitly start packaging; it does not depend on a tag or release event.

App-created PRs trigger CI, unlike PRs created with `GITHUB_TOKEN`. The App token
stays in the Release Please job and is revoked afterward. The packaging job uses
only `contents: write` on `GITHUB_TOKEN`; it does not need Issues, Pull requests,
Actions, Workflows, Administration, or any cicd-common access.

The publisher checks that the checkout, existing tag, manifest, and draft agree
on the version/full SHA, that the source is clean, and that the commit belongs
to `origin/master`. It runs validation without fixes or `go mod tidy`, rechecks
the draft before upload, and uses the exact tag via `GORELEASER_CURRENT_TAG`.
Release binary display is `vX.Y.Z` plus the full SHA. Archives use `X.Y.Z`.

GoReleaser OSS builds, archives, computes SHA-256, and uploads to the existing
draft while retaining Release Please's notes. It never automatically publishes.
The workflow verifies six archives, contents, checksums, and binary metadata and
retains build/coverage evidence for 14 days. OSS has no pre-publish artifact hook:
verification finishes after upload, while the release is still an unpublished
draft. Do not publish unless the entire packaging job succeeds.

Pinned tools: Release Please Action v5.0.0 (Release Please 17.6.0), GoReleaser
Action v7.2.3, GoReleaser OSS v2.18.2, Go 1.26.7, golangci-lint v2.10.1,
gotestsum v1.13.0, actionlint 1.7.12, and ShellCheck v0.11.0. Actions use full
commit SHAs. The development image and workflows share the GoReleaser version.

## Manual Setup

These steps require the repository owner's GitHub access. Never put the App
private key in chat, a commit, a command argument, or a workflow log.

1. Open <https://github.com/settings/apps/new> while signed in as
  `kenjones-cisco`, the verified personal account owning this public repository.
  Choose an available name such as `dinamo-release-please`. Homepage:
  `https://github.com/kenjones-cisco/dinamo`. Disable webhook **Active**; leave
  callback/setup URLs blank, OAuth authorization and device flow disabled.
  Select **Only on this account** when registering under `kenjones-cisco`.
2. Set repository permissions to **Contents: Read and write**, **Pull requests:
  Read and write**, and **Issues: Read and write**. Metadata read is automatic.
  Leave all other repository, organization, and account permissions at no
  access. Contents covers release-PR commits and release/tag creation; Pull
  requests covers PR creation/updates; Issues covers labels and PR comments.
  No webhooks or user OAuth tokens are needed.
3. Create the App, record its **Client ID** and numeric **App ID**, then use **Install App** to
  install on `kenjones-cisco`, selecting **Only select repositories: dinamo**.
  Check the installation permissions and repository selection afterward.
4. On the App's General page, **Generate a private key**. At
  <https://github.com/kenjones-cisco/dinamo/settings/secrets/actions>, create
  repository secret **RELEASE_APP_PRIVATE_KEY** with the complete PEM contents,
  including its header/footer and newlines. Enter the key directly into GitHub.
  On the **Variables** tab, create **RELEASE_APP_CLIENT_ID** with the App's
  Client ID (typically beginning with `Iv`), not the numeric App ID,
  installation ID, or client secret. The numeric App ID is not used by the
  workflow: the pinned token action recommends Client ID instead of its
  deprecated App ID input. Securely retain the key outside
  the repository; rotate it by replacing the secret before deleting the old
  App key. No PAT, client secret, or GoReleaser license is required.
5. At <https://github.com/kenjones-cisco/dinamo/settings/actions>, ensure Actions
  is enabled and the policy allows the pinned `actions/*`,
  `googleapis/release-please-action`, `goreleaser/goreleaser-action`, and existing
  Coveralls action. Docker must be usable on GitHub-hosted Ubuntu runners.
  Keep default **Workflow permissions: Read repository contents and packages**;
  the publisher requests write only in its own job. Release Please uses the
  App, so **Allow GitHub Actions to create and approve pull requests** is not
  required for its token. The App does not approve or merge its own PRs.
6. At <https://github.com/kenjones-cisco/dinamo/settings/rules>, review `master`
  protection: require PR review and successful **Check And Coverage** and
  **Release Preview** from GitHub Actions; wait for their first runs if they
  are not yet selectable. Keep any stronger existing requirements. Do not
  grant the App a branch-review/CI bypass or apply a creation restriction that
  prevents its `release-please--*` PR branches. Squash merges must be available.
7. Review tag rules for both historical tags and future `v*` tags. Block updates
  and deletion. Allow only the release App (and deliberate owner recovery) to
  create future release tags. If using a creation-restriction ruleset with an
  App bypass, keep no-update/no-delete rules in a separate ruleset without
  that bypass. Do not require signed tags: Release Please creates lightweight
  refs, not signed annotated tags. Never remove existing historical tags.
8. No release environment or required environment reviewers are configured by
  this design. Keep publication manual and limited to trusted maintainers with
  release-write access. Consider GitHub's immutable releases setting to protect
  artifacts after publication; enabling it is a separate owner decision, not
  part of this migration.

Completion checks: installation shows only `dinamo`; variable and secret names
match the workflow; App-token creation succeeds; the App opens a release PR and
both CI jobs run on it. The implementation session's token could not read Actions
administration settings (`403`), so the settings above require owner verification.
Nothing in this setup authorizes this assistant to change remote settings.

## First Release Checklist

1. Complete App setup and review the migration diff. Merge the migration using
  a Conventional Commit squash message with a blank-line-separated footer:

  ```text
  Release-As: 0.4.0
  ```

  This is a one-time commit instruction, not a permanently pinned config
  version. Do not omit it: maintenance-only commits otherwise propose 0.3.1.
  The pinned action's `release-as` input is ignored in manifest mode.
2. Confirm Release Please opens a PR for **0.4.0**, changes the manifest to
  `0.4.0`, and prepends release notes without dropping historical entries.
  Review the actual changes since bare tag `0.3.0`; older non-Conventional
  Commit messages may not produce automatic notes. Add necessary user-facing
  notes to the release PR and re-review after any automation update.
3. Wait for **Check And Coverage** and **Release Preview** and required human
  approval. Squash-merge the release PR; retain its generated title/body and
  `autorelease: pending` label. Do not rebase-merge a multi-commit release PR.
4. In <https://github.com/kenjones-cisco/dinamo/actions/workflows/release.yml>,
  confirm all release jobs succeed, the `v0.4.0` tag points at the release PR's
  merge SHA, and the release remains a draft. Never publish during an active run.
5. At <https://github.com/kenjones-cisco/dinamo/releases>, inspect the draft:
  Release Please notes, six `dinamo_0.4.0_<os>_<arch>.zip` files, and
  `dinamo_0.4.0_SHA256SUMS`. Each ZIP contains `dinamo` (`dinamo.exe` on Windows),
  `LICENSE`, `README.md`, and `CHANGELOG.md`. Download all seven assets to an
  empty directory and run `sha256sum --check dinamo_0.4.0_SHA256SUMS`.
  Run the matching native binary with `--version`: version must be `v0.4.0`
  and Git commit must match the full tag SHA. Inspect retained workflow evidence.
6. Only after successful verification, choose **Edit** on that draft and
  **Publish release**. Pre-1.0 is not automatically marked as a prerelease.
  Verify public assets/downloads and leave tags `0.1.0`, `0.2.0`, `0.3.0` alone.

Subsequent releases follow steps 2-6 with the proposed version; omit the
`Release-As` footer unless deliberately overriding a specific future release.

## Recovery

All workflow runs are serialized without cancelling a running release. Never
publish a draft while a packaging/recovery run is active: GitHub does not provide
an atomic lock between a draft-state check and asset upload. The workflow checks
state immediately before packaging, but human publication during upload would
violate that guarantee. Published or historical releases are never retry targets.

- **PR/App failure:** fix the installation or variable/secret directly in GitHub.
  Run **Release** on `master` with both recovery fields empty. This maintains
  release PRs and creates any still-pending release. Do not edit the manifest to
  trick the automation or remove existing tags/releases.
- **Tag exists but draft creation failed:** rerun Release Please with empty
  fields. The pinned implementation tolerates an existing tag; the publisher
  still verifies the exact tag SHA. Stop if the tag points anywhere else.
- **Draft exists but packaging failed or outputs were lost:** do not rerun the
  whole workflow expecting `release_created` again. Use **Run workflow** on
  `master` with `tag=v0.4.0` (or the affected new tag) and `sha=<full release PR
  merge SHA>`. This skips Release Please, checks out the exact commit, reruns
  validation/builds, and replaces matching assets only in that same unpublished
  draft. No tags, draft objects, notes, or published assets are deleted. An
  existing draft can contain partial uploads; do not publish until recovery is
  fully green and all checksums verify. CLI equivalent for a maintainer:

  ```bash
  gh workflow run release.yml --repo kenjones-cisco/dinamo --ref master \
   -f tag=v0.4.0 -f sha=<full-release-commit-sha>
  ```

- **Release source validation fails:** fix source through a reviewed PR and make
  a new version. Retrying the same SHA cannot fix its code. Do not move the tag
  to the repair commit. A new draft may supersede the failed draft; leave the
  latter unpublished and document the reason.
- **Already published:** retries fail closed. Ship corrections in a new release;
  never delete, reupload, retag, or overwrite published artifacts.

Recovery inputs must name an existing draft, exact `vX.Y.Z` tag and full SHA with
a matching manifest on `master`. Bare historical tags, missing drafts, mismatched
SHAs, dirty source, non-master commits, and API errors are rejected. The offline
tests invoke mocks only: they never create tags/releases or upload artifacts.

## Local Verification

```bash
CI=1 make check
CI=1 make test
CI=1 make test-race
CI=1 make cover
CI=1 make package
bash scripts/release_test.sh
bash scripts/verify-artifacts.sh snapshot
docker run --rm -v "$PWD:/work" -w /work rhysd/actionlint:1.7.12 -color
docker run --rm -v "$PWD:/work" -w /work koalaman/shellcheck:v0.11.0 \
  scripts/release.sh scripts/release_test.sh scripts/verify-artifacts.sh
```

Local snapshots explicitly identify as development builds, may use a dirty tree,
and never publish. Run artifact verification on Linux amd64; cross-platform
binaries are inspected using Go's build metadata, while the Linux amd64 binary
is also executed. A real hosted run is still needed to verify installation-token
access, App-triggered PR CI, protections, draft/tag creation, and asset uploads.

## Verified Migration Results

Verified locally on 2026-10-06 without committing, pushing, creating tags or
releases, publishing drafts, or changing remote settings:

- Pinned Go 1.26.7 lint: zero issues; tests: 17 passed; race tests: 17 passed.
- Statement coverage: 93.69%, enforcing the minimum of 80%.
- actionlint, ShellCheck, script syntax, and GoReleaser configuration: passed.
- Six-platform snapshot ZIPs, checksums, contents, Go build metadata, and Linux
  amd64 runtime version/full-SHA checks: passed.
- Non-publishing release-mode builds display `v0.4.0` and the expected full SHA;
  the verifier rejects their uncommitted source tree as intended. The test used
  an environment-only tag override, not a newly created tag.
- Offline Release Please 17.6.0 schema and real-history/mock-GitHub tests:
  bare-tag bootstrap, first-version footer, historical changelog preservation,
  manifest-only version tracking, and subsequent version increments passed.
- Draft retry and failed-upload recovery passed; 14 unsafe/error cases rejected.
- Retained `build`, `xcompile`, and `package` Make targets passed. Native build
  executed successfully despite an existing ignored container-ownership warning
  from Make's host-side `chmod`; this migration does not change that behavior.
- Historical tags, changelog, dependency files, and HEAD remain unchanged.

Pending human setup: register/install the App, enter its Client ID/private key
directly into GitHub, and review Actions/branch/tag protections. Hosted token
authentication, App-triggered PR CI, draft/tag creation, actual artifact uploads,
and a complete clean-release workflow require a real run after authorized merge.
Cross-platform binaries were built and inspected, not executed on native macOS,
Windows, or arm64 hosts. Publication must wait for a fully successful hosted run.

## References

- [Release Please Action and outputs](https://github.com/googleapis/release-please-action/tree/v5.0.0)
- [Manifest configuration](https://github.com/googleapis/release-please/blob/v17.6.0/docs/manifest-releaser.md)
- [Pinned draft/tag implementation](https://github.com/googleapis/release-please/blob/v17.6.0/src/github-api.ts)
- [GoReleaser releases and draft reuse](https://goreleaser.com/customization/publish/scm/)
- [GoReleaser archives](https://goreleaser.com/customization/package/archives/)
- [GitHub App registration](https://docs.github.com/en/apps/creating-github-apps/registering-a-github-app/registering-a-github-app)
- [GitHub token event restrictions](https://docs.github.com/en/actions/concepts/security/github_token)
