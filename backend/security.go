package backend

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"strings"
)

// parsePrivateKey はPEM形式の秘密鍵をパースする
func parsePrivateKey(keyData string) (interface{}, error) {
	block, _ := pem.Decode([]byte(keyData))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}

	// PKCS8形式の秘密鍵をパース
	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	// PKCS1形式のRSA秘密鍵をパース
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	// EC秘密鍵をパース
	if key, err := x509.ParseECPrivateKey(block.Bytes); err == nil {
		return key, nil
	}

	return nil, fmt.Errorf("unsupported private key format")
}

// parsePublicKey はPEM形式の公開鍵をパースする
func parsePublicKey(keyData string) (interface{}, error) {
	block, _ := pem.Decode([]byte(keyData))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}

	// PKIX形式の公開鍵をパース
	if key, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
		return key, nil
	}

	// PKCS1形式のRSA公開鍵をパース
	if key, err := x509.ParsePKCS1PublicKey(block.Bytes); err == nil {
		return key, nil
	}

	return nil, fmt.Errorf("unsupported public key format")
}

// isAsymmetricAlgorithm はアルゴリズムが非対称鍵暗号方式かどうかを判定する
func isAsymmetricAlgorithm(algorithm string) bool {
	return strings.HasPrefix(algorithm, "RS") ||
		strings.HasPrefix(algorithm, "ES") ||
		strings.HasPrefix(algorithm, "PS") ||
		algorithm == "EdDSA"
}

// generatePrivateKeyForAlgorithm は指定されたアルゴリズムに基づいて秘密鍵を生成する
func generatePrivateKeyForAlgorithm(algorithm string) (interface{}, error) {
	switch {
	case strings.HasPrefix(algorithm, "RS") || strings.HasPrefix(algorithm, "PS"):
		rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return nil, fmt.Errorf("failed to generate RSA key: %w", err)
		}
		return rsaKey, nil
	case strings.HasPrefix(algorithm, "ES"):
		var curve elliptic.Curve
		switch algorithm {
		case "ES256":
			curve = elliptic.P256()
		case "ES384":
			curve = elliptic.P384()
		case "ES512":
			curve = elliptic.P521()
		default:
			return nil, fmt.Errorf("unsupported ECDSA algorithm: %s", algorithm)
		}

		ecdsaKey, err := ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			return nil, fmt.Errorf("failed to generate ECDSA key: %w", err)
		}
		return ecdsaKey, nil
	case algorithm == "EdDSA":
		_, ed25519Key, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, fmt.Errorf("failed to generate Ed25519 key: %w", err)
		}
		return ed25519Key, nil
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
}

// getPublicKeyFromPrivate は秘密鍵から公開鍵を取得する
func getPublicKeyFromPrivate(privateKey interface{}) interface{} {
	switch key := privateKey.(type) {
	case *rsa.PrivateKey:
		return &key.PublicKey
	case *ecdsa.PrivateKey:
		return &key.PublicKey
	case ed25519.PrivateKey:
		return key.Public()
	default:
		return nil
	}
}

// parseSANEntries は SAN エントリをパースして DNS 名と IP アドレスに分離する
func parseSANEntries(sanEntries []string, defaultDNS string) ([]string, []net.IP) {
	dnsNames := make([]string, 0, len(sanEntries)+1)
	ipAddresses := make([]net.IP, 0, len(sanEntries))

	for _, san := range sanEntries {
		trimmed := strings.TrimSpace(san)
		if trimmed == "" {
			continue
		}
		if ip := net.ParseIP(trimmed); ip != nil {
			ipAddresses = append(ipAddresses, ip)
			continue
		}
		dnsNames = append(dnsNames, trimmed)
	}

	if len(dnsNames) == 0 && len(ipAddresses) == 0 && defaultDNS != "" {
		dnsNames = append(dnsNames, defaultDNS)
	}

	return dnsNames, ipAddresses
}

func getSignerFromPrivateKey(privateKey interface{}) (crypto.Signer, error) {
	signer, ok := privateKey.(crypto.Signer)
	if !ok {
		return nil, fmt.Errorf("private key does not implement crypto.Signer: %T", privateKey)
	}

	return signer, nil
}

func parseCertificateFromPEM(certificatePEM string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(certificatePEM))
	if block == nil {
		return nil, fmt.Errorf("invalid PEM data")
	}
	if block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("invalid PEM block type: %s", block.Type)
	}

	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return certificate, nil
}
