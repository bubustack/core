package connector

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bubustack/core/contracts"
)

func TestHubCredentialsUsesServerNameOverride(t *testing.T) {
	dir := t.TempDir()
	certFile, keyFile := writeTestTLSMaterial(t, dir)

	env := mapEnv{
		contracts.GRPCClientCertFileEnv: certFile,
		contracts.GRPCClientKeyFileEnv:  keyFile,
		contracts.GRPCCAFileEnv:         certFile,
		contracts.GRPCHubServerNameEnv:  "hub.internal",
	}

	tlsConfig, err := hubTLSConfig(env)
	if err != nil {
		t.Fatalf("hubTLSConfig failed: %v", err)
	}
	if got := tlsConfig.ServerName; got != "hub.internal" {
		t.Fatalf("expected server name override, got %q", got)
	}
	if tlsConfig.MinVersion != tls.VersionTLS13 {
		t.Fatalf("expected TLS 1.3 minimum, got %d", tlsConfig.MinVersion)
	}
}

func TestHubTLSConfigRequiresCA(t *testing.T) {
	dir := t.TempDir()
	certFile, keyFile := writeTestTLSMaterial(t, dir)

	env := mapEnv{
		contracts.GRPCClientCertFileEnv: certFile,
		contracts.GRPCClientKeyFileEnv:  keyFile,
	}

	if _, err := hubTLSConfig(env); err == nil {
		t.Fatalf("expected missing CA to be rejected")
	}
}

func writeTestTLSMaterial(t *testing.T, dir string) (string, string) {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "hub.internal",
		},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	der, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("create cert: %v", err)
	}

	certFile := filepath.Join(dir, "client.crt")
	keyFile := filepath.Join(dir, "client.key")

	if err := os.WriteFile(certFile, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der}), 0o600); err != nil {
		t.Fatalf("write cert: %v", err)
	}
	keyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	if err := os.WriteFile(
		keyFile,
		pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: keyBytes}),
		0o600,
	); err != nil {
		t.Fatalf("write key: %v", err)
	}

	return certFile, keyFile
}
