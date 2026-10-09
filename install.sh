#!/bin/sh
set -e

# Webex CLI installer
# Usage: curl -fsSL https://raw.githubusercontent.com/Cloverhound/webex-cli/main/install.sh | sh
#   Pin a version:   ... | sh -s -- -v 0.16.0   (or WEBEX_CLI_VERSION=0.16.0)
#   Install dir:     ... | INSTALL_DIR=/usr/local/bin sh

REPO="Cloverhound/webex-cli"
BINARY="webex"
INSTALL_DIR="${INSTALL_DIR:-${HOME}/.local/bin}"
RELEASES="https://github.com/${REPO}/releases"
GOPROXY_LATEST="https://proxy.golang.org/github.com/!cloverhound/webex-cli/@latest"
API_LATEST="https://api.github.com/repos/${REPO}/releases/latest"

VERSION="${WEBEX_CLI_VERSION:-}"
while [ $# -gt 0 ]; do
  case "$1" in
    -v|--version)
      [ $# -ge 2 ] || { echo "Error: $1 needs a version, e.g. $1 0.16.0"; exit 1; }
      VERSION="$2"; shift 2 ;;
    *) echo "Unknown option: $1"; exit 1 ;;
  esac
done
VERSION="${VERSION#v}"

# Detect OS
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case "$OS" in
  darwin) OS="darwin" ;;
  linux)  OS="linux" ;;
  *)      echo "Unsupported OS: $OS"; exit 1 ;;
esac

# Detect architecture
ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64)  ARCH="amd64" ;;
  arm64|aarch64)  ARCH="arm64" ;;
  *)              echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

TMPDIR=$(mktemp -d)
trap 'rm -rf "$TMPDIR"' EXIT

# The lookup avoids api.github.com where possible: sandboxes and shared CI
# hosts often block it or exhaust its unauthenticated rate limit.
version_from_checksums() {
  curl -fsSL "${RELEASES}/latest/download/checksums.txt" 2>/dev/null |
    sed -n "s/.*webex-cli_\([^ ]*\)_${OS}_${ARCH}\.tar\.gz\$/\1/p" | head -n 1
}

version_from_goproxy() {
  curl -fsSL "$GOPROXY_LATEST" 2>/dev/null |
    sed -n 's/.*"Version":"v\{0,1\}\([^"]*\)".*/\1/p' | head -n 1
}

version_from_api() {
  if [ -n "${GITHUB_TOKEN:-}" ]; then
    curl -fsSL -H "Authorization: Bearer ${GITHUB_TOKEN}" "$API_LATEST" 2>/dev/null
  else
    curl -fsSL "$API_LATEST" 2>/dev/null
  fi | grep '"tag_name"' | sed -E 's/.*"v?([^"]+)".*/\1/' | head -n 1
}

if [ -n "$VERSION" ]; then
  echo "Using pinned version: v${VERSION}"
else
  echo "Fetching latest release..."
  VERSION=$(version_from_checksums || true)
  [ -n "$VERSION" ] || VERSION=$(version_from_goproxy || true)
  [ -n "$VERSION" ] || VERSION=$(version_from_api || true)
  if [ -z "$VERSION" ]; then
    echo "Error: could not determine the latest version. Tried:"
    echo "  ${RELEASES}/latest/download/checksums.txt"
    echo "  ${GOPROXY_LATEST}"
    echo "  ${API_LATEST}"
    echo "Pin a version instead:"
    echo "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | WEBEX_CLI_VERSION=x.y.z sh"
    exit 1
  fi
  echo "Latest version: v${VERSION}"
fi

# Download
TARBALL="${BINARY}-cli_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="${RELEASES}/download/v${VERSION}/${TARBALL}"

echo "Downloading ${URL}..."
curl -fsSL "$URL" -o "${TMPDIR}/${TARBALL}"
curl -fsSL "${RELEASES}/download/v${VERSION}/checksums.txt" -o "${TMPDIR}/checksums.txt"

# Verify
EXPECTED=$(awk -v f="$TARBALL" '$2 == f || $2 == "*" f { print $1; exit }' "${TMPDIR}/checksums.txt")
if [ -z "$EXPECTED" ]; then
  echo "Error: ${TARBALL} is not listed in checksums.txt; refusing to install"
  exit 1
fi
if command -v sha256sum >/dev/null 2>&1; then
  ACTUAL=$(sha256sum "${TMPDIR}/${TARBALL}" | awk '{ print $1 }')
elif command -v shasum >/dev/null 2>&1; then
  ACTUAL=$(shasum -a 256 "${TMPDIR}/${TARBALL}" | awk '{ print $1 }')
else
  echo "Error: neither sha256sum nor shasum is available to verify the download"
  exit 1
fi
if [ "$ACTUAL" != "$EXPECTED" ]; then
  echo "Error: checksum mismatch for ${TARBALL}; refusing to install"
  echo "  expected: ${EXPECTED}"
  echo "  actual:   ${ACTUAL}"
  exit 1
fi
echo "Checksum verified."

# Extract
tar -xzf "${TMPDIR}/${TARBALL}" -C "$TMPDIR"

# Install
mkdir -p "$INSTALL_DIR"
mv "${TMPDIR}/${BINARY}" "${INSTALL_DIR}/${BINARY}"
chmod +x "${INSTALL_DIR}/${BINARY}"

echo "Installed ${BINARY} v${VERSION} to ${INSTALL_DIR}/${BINARY}"
echo ""

# Post-install prompts need a terminal; agents and CI run this without one,
# where post-install prints the PATH line and installs skills for detected
# agents. A post-install failure must not fail an install that succeeded.
if [ -z "${WEBEX_CLI_NONINTERACTIVE:-}" ] && [ -t 1 ] && [ -r /dev/tty ]; then
  "${INSTALL_DIR}/${BINARY}" post-install < /dev/tty || true
else
  WEBEX_CLI_NONINTERACTIVE=1 "${INSTALL_DIR}/${BINARY}" post-install < /dev/null || true
fi

echo ""
echo "Get started:"
echo "  webex config set client-id <your-client-id>      # if not using built-in defaults"
echo "  webex config set client-secret <your-client-secret>"
echo "  webex login"
