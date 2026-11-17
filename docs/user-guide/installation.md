# Installing DARP

**Classification:** Public
**Owner:** DARP CLI maintainers
**Source:** Created for the DARP CLI project
**Audience:** DARP CLI users

The official release assets are published in the
[`darp-labs/darp-cli` GitHub releases](https://github.com/darp-labs/darp-cli/releases).
DARP requires no Go installation.

## Supported targets

| Platform | Architecture | Archive |
|---|---|---|
| Windows | AMD64 | `darp_<version>_windows_amd64.zip` |
| Linux | AMD64 | `darp_<version>_linux_amd64.tar.gz` |
| Linux | ARM64 | `darp_<version>_linux_arm64.tar.gz` |
| macOS | AMD64 | `darp_<version>_darwin_amd64.tar.gz` |
| macOS | ARM64 | `darp_<version>_darwin_arm64.tar.gz` |

The archive contains only `darp.exe` on Windows or `darp` on Linux and macOS.

## PowerShell 7.6 or newer

Download the installer and run it for the latest stable release:

```powershell
Invoke-WebRequest `
  -Uri "https://raw.githubusercontent.com/darp-labs/darp-cli/main/scripts/install.ps1" `
  -OutFile "$env:TEMP\install-darp.ps1"
& "$env:TEMP\install-darp.ps1"
```

Install an exact version by passing a stable or prerelease semantic version:

```powershell
& "$env:TEMP\install-darp.ps1" -Version "v1.2.3"
```

The `v` prefix is optional. An explicit version that is missing or does not match
the resolved release fails; the installer never falls back to a different
release.

Windows installs to
`$env:LOCALAPPDATA\Programs\darp` and adds only that directory to the current
user's `Path`. It does not modify system `Path` or require elevation.

## POSIX shell

Download the installer and run it for the latest stable release:

```sh
curl -fsSL https://raw.githubusercontent.com/darp-labs/darp-cli/main/scripts/install.sh \
  -o /tmp/install-darp.sh
sh /tmp/install-darp.sh
```

Install an exact version:

```sh
sh /tmp/install-darp.sh v1.2.3
```

The `v` prefix is optional. Linux and macOS install to `$HOME/.local/bin` and
add it to the current shell and `$HOME/.profile` idempotently. The installer
does not change system paths or require elevation.

## Verification

The installer verifies, in this order:

1. The pinned Cosign bootstrap or an exact local Cosign `v3.1.4`.
2. The Sigstore-signed `checksums.txt.sigstore.json` bundle and its expected
   DARP repository, release workflow, tag, and source revision identity.
3. The selected archive's SHA-256 entry in `checksums.txt`.
4. The archive contents before replacing the existing executable.

GitHub Artifact Attestations are available for release archives, the checksum
manifest, the Cosign bundle, and bootstrap helper assets. Verify an asset with
the GitHub CLI:

```sh
gh attestation verify darp_1.2.3_linux_amd64.tar.gz \
  --repo darp-labs/darp-cli \
  --signer-workflow darp-labs/darp-cli/.github/workflows/release.yml \
  --predicate-type https://slsa.dev/provenance/v1
```

Replace the filename with the asset being verified.

## Common failures

- **Unsupported platform:** use one of the five supported target pairs.
- **PowerShell version:** require PowerShell 7.6 or newer.
- **Missing or invalid evidence:** keep the release asset, checksum manifest,
  Cosign bundle, and attestation available; do not bypass verification.
- **Signature or identity mismatch:** confirm the repository, workflow, tag,
  and source revision match `darp-labs/darp-cli` and
  `.github/workflows/release.yml`.
- **Checksum mismatch:** the archive or download was altered or incomplete;
  download it again.
- **Explicit version not found:** confirm the exact tag exists; the installer
  does not fall back to a latest or prerelease release.
- **Network or HTTPS download failure:** check network, proxy, and certificate
  configuration, then retry.

Any verification or installation failure exits non-zero. Existing installations
are preserved unless the verified executable has already been installed
successfully.
