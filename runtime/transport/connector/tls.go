package connector

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/bubustack/core/contracts"
	"google.golang.org/grpc/credentials"
)

// HubCredentials builds TLS credentials from the shared env vars.
func HubCredentials(env Env) (credentials.TransportCredentials, error) {
	tlsConfig, err := hubTLSConfig(env)
	if err != nil {
		return nil, err
	}
	return credentials.NewTLS(tlsConfig), nil
}

func hubTLSConfig(env Env) (*tls.Config, error) {
	env = ensureEnv(env)
	certFile := trimEnv(env, contracts.GRPCClientCertFileEnv)
	keyFile := trimEnv(env, contracts.GRPCClientKeyFileEnv)
	caFile := trimEnv(env, contracts.GRPCCAFileEnv)
	if certFile == "" || keyFile == "" || caFile == "" {
		return nil, fmt.Errorf("hub TLS assets not configured")
	}
	clientCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("load hub TLS cert/key: %w", err)
	}
	caBytes, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("read hub CA: %w", err)
	}
	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caBytes) {
		return nil, fmt.Errorf("parse hub CA")
	}
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS13,
		RootCAs:    caPool,
		Certificates: []tls.Certificate{
			clientCert,
		},
	}
	if serverName := trimEnv(env, contracts.GRPCHubServerNameEnv); serverName != "" {
		tlsConfig.ServerName = serverName
	}
	return tlsConfig, nil
}
