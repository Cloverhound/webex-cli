#!/bin/sh
# Publishes the packages staged by npm/stage.mjs: platform packages first, so
# the wrapper never points at versions that are missing from the registry.
# Versions already on the registry are skipped, so a failed release can be rerun.
set -eu

OUT="${1:-npm/build}"
VERSION=$(node -p "require('./${OUT}/webex-cli/package.json').version")
TAG=latest
case "$VERSION" in *-*) TAG=next ;; esac

for dir in "$OUT"/webex-cli-* "$OUT/webex-cli"; do
  name=$(node -p "require('./${dir}/package.json').name")
  if npm view "${name}@${VERSION}" version >/dev/null 2>&1; then
    echo "Skipping ${name}@${VERSION}: already published"
    continue
  fi
  npm publish "$dir" --access public --provenance --tag "$TAG"
done
