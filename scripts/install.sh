#!/usr/bin/env sh

set -eu

REPOSITORY=${DARP_REPOSITORY:-darp-labs/darp-cli}
VERSION=${1:-}
COSIGN_PIN=v3.1.4

fail() {
    printf 'darp installer: %s\n' "$*" >&2
    exit 1
}

sha256_file() {
    if command -v sha256sum >/dev/null 2>&1; then
        sha256sum "$1" | awk '{print tolower($1)}'
    elif command -v shasum >/dev/null 2>&1; then
        shasum -a 256 "$1" | awk '{print tolower($1)}'
    else
        fail 'SHA-256 utility is required'
    fi
}

download() {
    url=$1
    output=$2
    case "$url" in
        https://*) ;;
        *) fail "refusing insecure download: $url" ;;
    esac
    if ! curl -fsSL --proto '=https' --proto-redir '=https' --tlsv1.2 "$url" -o "$output"; then
        fail "failed to download $url"
    fi
    if [ ! -s "$output" ]; then
        fail "download is missing or empty: $url"
    fi
}

valid_semver() {
    printf '%s\n' "$1" | grep -Eq '^v?(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$'
}

json_value() {
    key=$1
    file=$2
    sed -n "s/.*\"$key\"[[:space:]]*:[[:space:]]*\"\\([^\"]*\\)\".*/\\1/p" "$file" | head -n 1
}

command -v curl >/dev/null 2>&1 || fail 'curl is required'

UNAME_OS=$(uname -s)
UNAME_ARCH=$(uname -m)
case "$UNAME_OS/$UNAME_ARCH" in
    Linux/x86_64|Linux/amd64) OS=linux; ARCH=amd64 ;;
    Linux/aarch64|Linux/arm64) OS=linux; ARCH=arm64 ;;
    Darwin/x86_64|Darwin/amd64) OS=darwin; ARCH=amd64 ;;
    Darwin/arm64|Darwin/aarch64) OS=darwin; ARCH=arm64 ;;
    *) fail "unsupported platform: $UNAME_OS/$UNAME_ARCH" ;;
esac

API_BASE=https://api.github.com/repos/$REPOSITORY
RELEASE_JSON=${TMPDIR:-/tmp}/darp-release-$$.json
trap 'rm -f "$RELEASE_JSON"' EXIT HUP INT TERM

if [ -n "$VERSION" ]; then
    valid_semver "$VERSION" || fail "invalid explicit version '$VERSION'"
    case "$VERSION" in v*) TAG=$VERSION ;; *) TAG=v$VERSION ;; esac
    if ! curl -fsSL --proto '=https' --proto-redir '=https' --tlsv1.2 \
        "$API_BASE/releases/tags/$TAG" -o "$RELEASE_JSON"; then
        fail "explicit release '$VERSION' was not found; no fallback is used"
    fi
else
    if ! curl -fsSL --proto '=https' --proto-redir '=https' --tlsv1.2 \
        "$API_BASE/releases/latest" -o "$RELEASE_JSON"; then
        fail 'could not resolve the latest stable release'
    fi
fi

TAG=$(json_value tag_name "$RELEASE_JSON")
[ -n "$TAG" ] || fail 'release returned an invalid tag'
valid_semver "$TAG" || fail 'release returned an invalid semantic version'
NUMBER=${TAG#v}

if [ -n "$VERSION" ]; then
    REQUESTED=${VERSION#v}
    [ "$NUMBER" = "$REQUESTED" ] || fail "resolved release $TAG does not match explicit version $VERSION"
fi

COMMIT_JSON=${TMPDIR:-/tmp}/darp-commit-$$.json
if ! curl -fsSL --proto '=https' --proto-redir '=https' --tlsv1.2 \
    "$API_BASE/commits/$TAG" -o "$COMMIT_JSON"; then
    rm -f "$COMMIT_JSON"
    fail "could not resolve the source revision for release $TAG"
fi
SOURCE_SHA=$(json_value sha "$COMMIT_JSON")
rm -f "$COMMIT_JSON"
printf '%s' "$SOURCE_SHA" | grep -Eq '^[0-9a-fA-F]{40}$' || fail "release $TAG returned an invalid source revision"
SOURCE_SHA=$(printf '%s' "$SOURCE_SHA" | tr 'A-F' 'a-f')

ASSET=darp_${NUMBER}_${OS}_${ARCH}.tar.gz
DOWNLOAD_BASE=https://github.com/$REPOSITORY/releases/download/$TAG
TEMP_DIR=$(mktemp -d "${TMPDIR:-/tmp}/darp-install-XXXXXX") || fail 'could not create a temporary directory'
trap 'rm -rf "$TEMP_DIR"; rm -f "$RELEASE_JSON"' EXIT HUP INT TERM

download "$DOWNLOAD_BASE/$ASSET" "$TEMP_DIR/$ASSET"
download "$DOWNLOAD_BASE/checksums.txt" "$TEMP_DIR/checksums.txt"
download "$DOWNLOAD_BASE/checksums.txt.sigstore.json" "$TEMP_DIR/checksums.txt.sigstore.json"

COSIGN=$(command -v cosign 2>/dev/null || true)
if [ -n "$COSIGN" ]; then
    COSIGN_VERSION=$(cosign version 2>/dev/null | sed -n 's/.*[Vv]ersion:[[:space:]]*v\?\([0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*\).*/\1/p' | head -n 1)
    [ "$COSIGN_VERSION" = "${COSIGN_PIN#v}" ] || COSIGN=
fi

if [ -z "$COSIGN" ]; then
    HELPER=darp-trust-bootstrap_${OS}_${ARCH}
    case "$OS/$ARCH" in
        linux/amd64) HELPER_HASH=ca16fc2e9c49a036910bbad5fb72932c18092ba3e7af6f95553dc47806459348 ;;
        linux/arm64) HELPER_HASH=01da55925a3626e4c510496e292483e10a2b913ba1747f7a55ae558a01026211 ;;
        darwin/amd64) HELPER_HASH=fd58a03708004c6e62266254386fbb186e14c8f15fb1e7a2c60cc680b6c8085a ;;
        darwin/arm64) HELPER_HASH=ba5077a8f014e3d4a705d4c8beb330c5d96854a5eb269bbafb213e435d1943f2 ;;
        *) fail "unsupported bootstrap platform: $OS/$ARCH" ;;
    esac
    download "$DOWNLOAD_BASE/$HELPER" "$TEMP_DIR/$HELPER"
    chmod 0700 "$TEMP_DIR/$HELPER"
    [ "$(sha256_file "$TEMP_DIR/$HELPER")" = "$HELPER_HASH" ] || fail 'bootstrap helper SHA-256 mismatch'
    COSIGN=$("$TEMP_DIR/$HELPER" -version "$COSIGN_PIN") || fail 'authenticated Cosign bootstrap failed'
    [ -n "$COSIGN" ] && [ -x "$COSIGN" ] || fail 'authenticated Cosign bootstrap did not return an executable'
fi

"$COSIGN" verify-blob \
    --bundle "$TEMP_DIR/checksums.txt.sigstore.json" \
    --certificate-identity "https://github.com/$REPOSITORY/.github/workflows/release.yml@refs/tags/$TAG" \
    --certificate-oidc-issuer https://token.actions.githubusercontent.com \
    --certificate-github-workflow-repository "$REPOSITORY" \
    --certificate-github-workflow-ref "refs/tags/$TAG" \
    --certificate-github-workflow-sha "$SOURCE_SHA" \
    "$TEMP_DIR/checksums.txt" >/dev/null ||
    fail 'signature verification failed'

EXPECTED_HASH=$(awk -v asset="$ASSET" '$0 ~ ("[[:space:]]\\*?" asset "$") { print tolower($1); count++ } END { if (count == 1) exit 0; exit 1 }' "$TEMP_DIR/checksums.txt") ||
    fail "checksum manifest is missing a unique entry for $ASSET"
printf '%s' "$EXPECTED_HASH" | grep -Eq '^[0-9a-f]{64}$' || fail "checksum manifest contains an invalid SHA-256 for $ASSET"
[ "$(sha256_file "$TEMP_DIR/$ASSET")" = "$EXPECTED_HASH" ] ||
    fail "artifact SHA-256 mismatch for $ASSET"

mkdir -p "$TEMP_DIR/extract"
tar -xzf "$TEMP_DIR/$ASSET" -C "$TEMP_DIR/extract" ||
    fail 'failed to extract release archive'
CANDIDATE=$TEMP_DIR/extract/darp
[ -f "$CANDIDATE" ] && [ -s "$CANDIDATE" ] || fail 'archive does not contain a valid root-level darp'

INSTALL_DIR=$HOME/.local/bin
TARGET=$INSTALL_DIR/darp
mkdir -p "$INSTALL_DIR"
PROFILE=$HOME/.profile
BACKUP=
if [ -f "$TARGET" ]; then
    BACKUP=$TEMP_DIR/darp.backup
    cp -p "$TARGET" "$BACKUP" || fail 'could not preserve the existing installation'
fi
if ! cp "$CANDIDATE" "$TARGET.tmp"; then
    [ -z "$BACKUP" ] || cp -p "$BACKUP" "$TARGET"
    fail 'could not stage the new executable'
fi
chmod 0755 "$TARGET.tmp"
if ! mv -f "$TARGET.tmp" "$TARGET"; then
    [ -z "$BACKUP" ] || cp -p "$BACKUP" "$TARGET"
    fail 'could not replace the existing installation'
fi

# Prepare the profile change beside the profile itself so the final rename is
# atomic. Do this only after the executable is installed; if preparation or
# replacement fails, restore the previous executable and leave the profile as-is.
PROFILE_TEMP=
if [ ! -f "$PROFILE" ] || ! grep -Fqx 'export PATH="$HOME/.local/bin:$PATH"' "$PROFILE" 2>/dev/null; then
    PROFILE_TEMP=$(mktemp "$HOME/.profile.darp.XXXXXX") || {
        [ -z "$BACKUP" ] || cp -p "$BACKUP" "$TARGET"
        [ -z "$BACKUP" ] && rm -f "$TARGET"
        fail 'could not stage the current-user PATH update'
    }
    if [ -f "$PROFILE" ] && ! cp -p "$PROFILE" "$PROFILE_TEMP"; then
        rm -f "$PROFILE_TEMP"
        [ -z "$BACKUP" ] || cp -p "$BACKUP" "$TARGET"
        [ -z "$BACKUP" ] && rm -f "$TARGET"
        fail 'could not preserve the current-user profile'
    fi
    if ! printf '\nexport PATH="$HOME/.local/bin:$PATH"\n' >>"$PROFILE_TEMP"; then
        rm -f "$PROFILE_TEMP"
        [ -z "$BACKUP" ] || cp -p "$BACKUP" "$TARGET"
        [ -z "$BACKUP" ] && rm -f "$TARGET"
        fail 'could not stage the current-user PATH update'
    fi
    if ! mv -f "$PROFILE_TEMP" "$PROFILE"; then
        rm -f "$PROFILE_TEMP"
        [ -z "$BACKUP" ] || cp -p "$BACKUP" "$TARGET"
        [ -z "$BACKUP" ] && rm -f "$TARGET"
        fail 'could not update the current-user profile; the previous installation was restored'
    fi
fi

export PATH="$INSTALL_DIR:$PATH"

printf 'Installed darp %s to %s\n' "$TAG" "$TARGET"
