// darp-trust-bootstrap authenticates a pinned Cosign release before it is run.
// It is a temporary release bootstrap helper, not part of the DARP CLI.
package main

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	_ "embed"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"time"

	"github.com/theupdateframework/go-tuf/v2/metadata/config"
	"github.com/theupdateframework/go-tuf/v2/metadata/updater"
)

const (
	tufMetadataURL = "https://tuf-repo-cdn.sigstore.dev"
	tufTargetsURL  = "https://tuf-repo-cdn.sigstore.dev/targets"
	maxSignature   = 1 << 20
	maxCosignSize  = 256 << 20
)

var stableCosignVersion = regexp.MustCompile(`^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$`)

//go:embed 10.root.json
var initialRoot []byte

type cosignBundle struct {
	MessageSignature struct {
		Signature string `json:"signature"`
	} `json:"messageSignature"`
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "darp-trust-bootstrap:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet("verify-cosign", flag.ContinueOnError)
	version := flags.String("version", "", "pinned Cosign release (for example v3.1.4)")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 || !stableCosignVersion.MatchString(*version) {
		return errors.New("usage: darp-trust-bootstrap -version vX.Y.Z (stable Cosign version required)")
	}
	path, err := bootstrapCosign(*version, runtime.GOOS, runtime.GOARCH)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(os.Stdout, path)
	return err
}

func bootstrapCosign(version, goos, goarch string) (string, error) {
	asset, err := cosignAssetName(goos, goarch)
	if err != nil {
		return "", err
	}
	artifactKey, err := fetchArtifactKey()
	if err != nil {
		return "", fmt.Errorf("authenticate current Sigstore TUF metadata before downloading Cosign: %w", err)
	}
	workDir, err := os.MkdirTemp("", "darp-cosign-bootstrap-*")
	if err != nil {
		return "", fmt.Errorf("create temporary Cosign directory: %w", err)
	}
	// Keep the verified binary for the installer to execute. The unverified
	// download and signature bundle are removed before returning.
	binaryPath := filepath.Join(workDir, asset)
	bundleName := cosignSignatureAssetName(asset)
	bundlePath := filepath.Join(workDir, bundleName)
	baseURL := "https://github.com/sigstore/cosign/releases/download/" + url.PathEscape(version) + "/"
	if err := download(baseURL+url.PathEscape(asset), binaryPath, maxCosignSize, 0700); err != nil {
		os.RemoveAll(workDir)
		return "", fmt.Errorf("download pinned Cosign %s: %w", version, err)
	}
	if err := download(baseURL+url.PathEscape(bundleName), bundlePath, maxSignature, 0600); err != nil {
		os.RemoveAll(workDir)
		return "", fmt.Errorf("download pinned Cosign signature bundle: %w", err)
	}
	if err := verifyCosignWithKey(version, binaryPath, bundlePath, artifactKey); err != nil {
		os.RemoveAll(workDir)
		return "", err
	}
	if err := os.Remove(bundlePath); err != nil {
		os.RemoveAll(workDir)
		return "", fmt.Errorf("remove temporary signature bundle: %w", err)
	}
	return binaryPath, nil
}

func cosignAssetName(goos, goarch string) (string, error) {
	if goos == "windows" && goarch == "amd64" {
		return "cosign-windows-amd64.exe", nil
	}
	if (goos == "linux" || goos == "darwin") && (goarch == "amd64" || goarch == "arm64") {
		return "cosign-" + goos + "-" + goarch, nil
	}
	return "", fmt.Errorf("unsupported platform for Cosign bootstrap: %s/%s", goos, goarch)
}

func cosignSignatureAssetName(binaryAsset string) string {
	return binaryAsset + "-kms.sigstore.json"
}

func verifyCosign(version, binaryPath, bundlePath string) error {
	binary, signature, err := readCosignInputs(binaryPath, bundlePath)
	if err != nil {
		return err
	}
	key, err := fetchArtifactKey()
	if err != nil {
		return err
	}
	return verifyCosignBytes(version, binary, signature, key)
}

func verifyCosignWithKey(version, binaryPath, bundlePath string, key []byte) error {
	binary, signature, err := readCosignInputs(binaryPath, bundlePath)
	if err != nil {
		return err
	}
	return verifyCosignBytes(version, binary, signature, key)
}

func readCosignInputs(binaryPath, bundlePath string) ([]byte, []byte, error) {
	info, err := os.Stat(binaryPath)
	if err != nil {
		return nil, nil, fmt.Errorf("stat Cosign binary: %w", err)
	}
	if !info.Mode().IsRegular() || info.Size() == 0 || info.Size() > maxCosignSize {
		return nil, nil, errors.New("Cosign binary is not a regular file or exceeds the size limit")
	}
	binary, err := os.ReadFile(binaryPath)
	if err != nil {
		return nil, nil, fmt.Errorf("read Cosign binary: %w", err)
	}
	if len(binary) == 0 {
		return nil, nil, errors.New("Cosign binary is empty")
	}
	bundleBytes, err := os.ReadFile(bundlePath)
	if err != nil {
		return nil, nil, fmt.Errorf("read Cosign signature bundle: %w", err)
	}
	if len(bundleBytes) == 0 || len(bundleBytes) > maxSignature {
		return nil, nil, errors.New("Cosign signature bundle is empty or exceeds the size limit")
	}
	signature, err := parseBundleSignature(bundleBytes)
	if err != nil {
		return nil, nil, err
	}
	return binary, signature, nil
}

func verifyCosignBytes(version string, binary, signature, key []byte) error {
	if err := verifySignature(key, binary, signature); err != nil {
		return fmt.Errorf("Cosign %s failed Sigstore artifact-key verification: %w", version, err)
	}
	return nil
}

func parseBundleSignature(bundleBytes []byte) ([]byte, error) {
	var bundle cosignBundle
	if err := json.Unmarshal(bundleBytes, &bundle); err != nil {
		return nil, fmt.Errorf("parse Cosign signature bundle: %w", err)
	}
	signature, err := base64.StdEncoding.DecodeString(bundle.MessageSignature.Signature)
	if err != nil || len(signature) == 0 {
		return nil, errors.New("Cosign signature bundle has no valid message signature")
	}
	return signature, nil
}

func download(rawURL, destination string, maxBytes int64, mode os.FileMode) error {
	parsedURL, err := url.Parse(rawURL)
	if err != nil || parsedURL.Scheme != "https" || parsedURL.Host == "" {
		return errors.New("bootstrap downloads require an absolute HTTPS URL")
	}
	client := &http.Client{
		Timeout: 2 * time.Minute,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if req.URL.Scheme != "https" {
				return errors.New("refusing HTTPS downgrade during bootstrap download")
			}
			if len(via) >= 10 {
				return errors.New("too many bootstrap download redirects")
			}
			return nil
		},
	}
	request, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return err
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned HTTP %s", response.Status)
	}
	file, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, mode)
	if err != nil {
		return err
	}
	defer file.Close()
	written, err := io.Copy(file, io.LimitReader(response.Body, maxBytes+1))
	if err != nil {
		os.Remove(destination)
		return err
	}
	if written == 0 || written > maxBytes {
		os.Remove(destination)
		return fmt.Errorf("download size %d is empty or exceeds %d bytes", written, maxBytes)
	}
	if err := file.Sync(); err != nil {
		os.Remove(destination)
		return err
	}
	return file.Close()
}

func fetchArtifactKey() ([]byte, error) {
	return fetchArtifactKeyWith(initialRoot, tufMetadataURL, tufTargetsURL, &http.Client{
		Timeout: 45 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if req.URL.Scheme != "https" {
				return errors.New("refusing HTTPS downgrade while updating Sigstore TUF metadata")
			}
			if len(via) >= 10 {
				return errors.New("too many Sigstore TUF metadata redirects")
			}
			return nil
		},
	})
}

func fetchArtifactKeyWith(root []byte, metadataURL, targetsURL string, httpClient *http.Client) ([]byte, error) {
	workDir, err := os.MkdirTemp("", "darp-tuf-*")
	if err != nil {
		return nil, fmt.Errorf("create temporary TUF cache: %w", err)
	}
	defer os.RemoveAll(workDir)

	cfg, err := config.New(metadataURL, root)
	if err != nil {
		return nil, fmt.Errorf("configure Sigstore TUF client: %w", err)
	}
	cfg.LocalMetadataDir = filepath.Join(workDir, "metadata")
	cfg.LocalTargetsDir = filepath.Join(workDir, "targets")
	if httpClient == nil {
		return nil, errors.New("TUF metadata HTTP client is required")
	}
	if err := cfg.SetDefaultFetcherHTTPClient(httpClient); err != nil {
		return nil, fmt.Errorf("configure bounded TUF downloads: %w", err)
	}
	client, err := updater.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("initialize Sigstore TUF client from pinned root: %w", err)
	}
	if err := client.Refresh(); err != nil {
		return nil, fmt.Errorf("refresh Sigstore TUF metadata (root rotation, signatures, versions, and expiry must validate): %w", err)
	}
	target, err := client.GetTargetInfo("artifact.pub")
	if err != nil {
		return nil, fmt.Errorf("resolve artifact.pub from authenticated Sigstore targets metadata: %w", err)
	}
	_, contents, err := client.DownloadTarget(target, "", targetsURL)
	if err != nil {
		return nil, fmt.Errorf("download authenticated Sigstore artifact.pub: %w", err)
	}
	if len(contents) == 0 {
		return nil, errors.New("authenticated Sigstore artifact.pub is empty")
	}
	return contents, nil
}

func verifySignature(publicKeyPEM, message, signature []byte) error {
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		return errors.New("TUF artifact.pub is not PEM encoded")
	}
	var key any
	if parsed, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
		key = parsed
	} else if parsed, err := x509.ParsePKCS1PublicKey(block.Bytes); err == nil {
		key = parsed
	} else {
		return errors.New("TUF artifact.pub contains an unsupported public key")
	}
	digest := sha256.Sum256(message)
	switch publicKey := key.(type) {
	case *rsa.PublicKey:
		if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, digest[:], signature); err != nil {
			return errors.New("RSA signature verification failed")
		}
	case *ecdsa.PublicKey:
		var parsed struct{ R, S *big.Int }
		if _, err := asn1.Unmarshal(signature, &parsed); err != nil || parsed.R == nil || parsed.S == nil {
			return errors.New("invalid ECDSA signature encoding")
		}
		if !ecdsa.Verify(publicKey, digest[:], parsed.R, parsed.S) {
			return errors.New("ECDSA signature verification failed")
		}
	default:
		return errors.New("TUF artifact.pub uses an unsupported key type")
	}
	return nil
}
