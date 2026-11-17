package main

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/theupdateframework/go-tuf/v2/metadata"
	"github.com/theupdateframework/go-tuf/v2/metadata/trustedmetadata"
)

func TestVerifySignatureAcceptsAuthenticCosignBytes(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	message := []byte("verified Cosign release bytes")
	digest := sha256.Sum256(message)
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, digest[:])
	if err != nil {
		t.Fatal(err)
	}
	publicDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Fatal(err)
	}
	publicPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: publicDER})
	if err := verifySignature(publicPEM, message, signature); err != nil {
		t.Fatalf("valid signature rejected: %v", err)
	}
	if err := verifySignature(publicPEM, []byte("modified release bytes"), signature); err == nil {
		t.Fatal("modified binary accepted")
	}
}

func TestCosignAssetNameAllowsOnlySupportedTargets(t *testing.T) {
	cases := map[string]string{
		"linux/amd64":   "cosign-linux-amd64",
		"linux/arm64":   "cosign-linux-arm64",
		"darwin/amd64":  "cosign-darwin-amd64",
		"darwin/arm64":  "cosign-darwin-arm64",
		"windows/amd64": "cosign-windows-amd64.exe",
	}
	for target, expected := range cases {
		name, err := cosignAssetName(splitTarget(target))
		if err != nil || name != expected {
			t.Errorf("cosignAssetName(%q) = %q, %v; want %q", target, name, err, expected)
		}
		if got, want := cosignSignatureAssetName(name), expected+"-kms.sigstore.json"; got != want {
			t.Errorf("signature asset for %q = %q; want %q", expected, got, want)
		}
	}
	if _, err := cosignAssetName("windows", "arm64"); err == nil {
		t.Fatal("unsupported target was accepted")
	}
}

func splitTarget(target string) (string, string) {
	for i := range target {
		if target[i] == '/' {
			return target[:i], target[i+1:]
		}
	}
	return target, ""
}

func TestStableCosignVersionValidation(t *testing.T) {
	for _, version := range []string{"v3.1.4", "v0.1.0"} {
		if !stableCosignVersion.MatchString(version) {
			t.Errorf("stable version %q rejected", version)
		}
	}
	for _, version := range []string{"", "3.1.4", "v03.1.4", "v3.1.4-rc.1", "v3.1"} {
		if stableCosignVersion.MatchString(version) {
			t.Errorf("invalid version %q accepted", version)
		}
	}
}

func TestBundleParsingRequiresSignature(t *testing.T) {
	valid := []byte(`{"messageSignature":{"signature":"` + base64.StdEncoding.EncodeToString([]byte("sig")) + `"}}`)
	signature, err := parseBundleSignature(valid)
	if err != nil || string(signature) != "sig" {
		t.Fatalf("valid signature not decoded: %q, %v", signature, err)
	}
	for _, invalid := range [][]byte{
		[]byte("not json"),
		[]byte(`{"messageSignature":{}}`),
		[]byte(`{"messageSignature":{"signature":"%%%"}}`),
	} {
		if _, err := parseBundleSignature(invalid); err == nil {
			t.Errorf("invalid bundle was accepted: %s", invalid)
		}
	}
}

func TestDownloadRejectsHTTPURLs(t *testing.T) {
	if err := download("http://example.invalid/cosign", "unused", 10, 0600); err == nil {
		t.Fatal("insecure bootstrap URL was accepted")
	}
}

func TestVerifyCosignRejectsMissingFilesBeforeTrustLookup(t *testing.T) {
	if err := verifyCosign("v3.1.4", filepath.Join(t.TempDir(), "missing-cosign"), filepath.Join(t.TempDir(), "missing-bundle")); err == nil {
		t.Fatal("missing binary was accepted")
	}
	binaryPath := filepath.Join(t.TempDir(), "cosign")
	if err := os.WriteFile(binaryPath, []byte("cosign"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := verifyCosign("v3.1.4", binaryPath, filepath.Join(t.TempDir(), "missing-bundle")); err == nil {
		t.Fatal("missing signature bundle was accepted")
	}
}

func TestVerifyCosignRejectsEmptyAndOversizedBinary(t *testing.T) {
	tempDir := t.TempDir()
	bundlePath := filepath.Join(tempDir, "bundle")
	if err := os.WriteFile(bundlePath, []byte(`{"messageSignature":{"signature":"c2ln"}}`), 0600); err != nil {
		t.Fatal(err)
	}
	for name, contents := range map[string][]byte{"empty": {}, "oversized": make([]byte, maxCosignSize+1)} {
		t.Run(name, func(t *testing.T) {
			binaryPath := filepath.Join(t.TempDir(), "cosign")
			if err := os.WriteFile(binaryPath, contents, 0600); err != nil {
				t.Fatal(err)
			}
			if err := verifyCosign("v3.1.4", binaryPath, bundlePath); err == nil {
				t.Fatal("invalid binary was accepted")
			}
		})
	}
}

func TestExpiredFinalTUFRootIsRejected(t *testing.T) {
	trusted, err := trustedmetadata.New(initialRoot)
	if err != nil {
		t.Fatalf("historical bootstrap root should be accepted as the starting anchor: %v", err)
	}
	_, err = trusted.UpdateTimestamp([]byte(`{}`))
	var expired *metadata.ErrExpiredMetadata
	if !errors.As(err, &expired) {
		t.Fatalf("expired final root must reject timestamp metadata, got %v", err)
	}
}

func TestTUFUpdaterRejectsInvalidRootRotationMetadata(t *testing.T) {
	requestedRoot := false
	client := &http.Client{Transport: responseTransport(func(req *http.Request) *http.Response {
		if req.URL.Path == "/11.root.json" {
			requestedRoot = true
			return httpResponse(req, http.StatusOK, []byte("not valid TUF metadata"))
		}
		return httpResponse(req, http.StatusNotFound, nil)
	})}
	_, err := fetchArtifactKeyWith(initialRoot, "https://tuf.invalid", "https://tuf.invalid/targets", client)
	if err == nil {
		t.Fatal("invalid TUF root rotation metadata was accepted")
	}
	if !requestedRoot {
		t.Fatalf("expected the TUF client to request the next root before rejecting it: %v", err)
	}
}

func TestTUFUpdaterRejectsUnauthenticatedRootRotation(t *testing.T) {
	unsignedRoot := bytes.Replace(initialRoot, []byte(`"version": 10`), []byte(`"version": 11`), 1)
	if bytes.Equal(unsignedRoot, initialRoot) {
		t.Fatal("failed to prepare changed root metadata")
	}
	client := &http.Client{Transport: responseTransport(func(req *http.Request) *http.Response {
		if req.URL.Path == "/11.root.json" {
			return httpResponse(req, http.StatusOK, unsignedRoot)
		}
		return httpResponse(req, http.StatusNotFound, nil)
	})}
	if _, err := fetchArtifactKeyWith(initialRoot, "https://tuf.invalid", "https://tuf.invalid/targets", client); err == nil {
		t.Fatal("root metadata with signatures over different contents was accepted")
	}
}

func TestTUFUpdaterRejectsMissingMetadata(t *testing.T) {
	client := &http.Client{Transport: responseTransport(func(req *http.Request) *http.Response {
		return httpResponse(req, http.StatusNotFound, nil)
	})}
	if _, err := fetchArtifactKeyWith(initialRoot, "https://tuf.invalid", "https://tuf.invalid/targets", client); err == nil {
		t.Fatal("missing TUF metadata was accepted")
	}
}

type responseTransport func(*http.Request) *http.Response

func (transport responseTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return transport(request), nil
}

func httpResponse(request *http.Request, status int, body []byte) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(string(body))),
		Request:    request,
	}
}
