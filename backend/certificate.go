package backend

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type SelfSignedCertificateResult struct {
	CertificatePEM string `json:"certificatePem"`
	PrivateKeyPEM  string `json:"privateKeyPem"`
}

// MTLSCertificateResult は mTLS 用の証明書セットを保持する構造体
type MTLSCertificateResult struct {
	CACertificatePEM     string `json:"caCertificatePem"`
	CAPrivateKeyPEM      string `json:"caPrivateKeyPem"`
	ServerCertificatePEM string `json:"serverCertificatePem"`
	ServerPrivateKeyPEM  string `json:"serverPrivateKeyPem"`
	ClientCertificatePEM string `json:"clientCertificatePem"`
	ClientPrivateKeyPEM  string `json:"clientPrivateKeyPem"`
}

// ClientCertificate は個別のクライアント証明書情報を保持する構造体
type ClientCertificate struct {
	CommonName     string `json:"commonName"`
	CertificatePEM string `json:"certificatePem"`
	PrivateKeyPEM  string `json:"privateKeyPem"`
}

// MTLSCertificatesMultiClientResult は複数のクライアント証明書を含む mTLS 証明書セット
type MTLSCertificatesMultiClientResult struct {
	CACertificatePEM     string              `json:"caCertificatePem"`
	CAPrivateKeyPEM      string              `json:"caPrivateKeyPem"`
	CRLPEM               string              `json:"crlPem"`
	ServerCertificatePEM string              `json:"serverCertificatePem"`
	ServerPrivateKeyPEM  string              `json:"serverPrivateKeyPem"`
	ClientCertificates   []ClientCertificate `json:"clientCertificates"`
}

// CRLUpdateResult は CRL 更新結果を保持する構造体
type CRLUpdateResult struct {
	CRLPEM               string   `json:"crlPem"`
	AddedCount           int      `json:"addedCount"`
	TotalRevokedCount    int      `json:"totalRevokedCount"`
	RevokedSerialNumbers []string `json:"revokedSerialNumbers"`
}

// VerifyKeyPair は秘密鍵と公開鍵のペアが一致しているかを検証する
// 秘密鍵から公開鍵を抽出して、与えられた公開鍵と比較することで検証を行う
func VerifyKeyPair(privateKey string, publicKey string, algorithm string) (bool, error) {
	privateKey = strings.TrimSpace(privateKey)
	publicKey = strings.TrimSpace(publicKey)

	if privateKey == "" {
		return false, fmt.Errorf("private key cannot be empty")
	}
	if publicKey == "" {
		return false, fmt.Errorf("public key cannot be empty")
	}

	// サポートされているアルゴリズムかチェック
	if !isAsymmetricAlgorithm(algorithm) {
		return false, fmt.Errorf("algorithm %s does not support key pair verification (only RSA, ECDSA, and EdDSA are supported)", algorithm)
	}

	// 秘密鍵から公開鍵を抽出
	extractedPublicKey, err := ExtractPublicKey(privateKey, algorithm)
	if err != nil {
		return false, fmt.Errorf("failed to extract public key from private key: %w", err)
	}

	// 抽出した公開鍵と与えられた公開鍵をDER形式で比較
	// PEM形式でのホワイトスペースの差異を無視するため、DER形式で比較
	extractedBlock, _ := pem.Decode([]byte(extractedPublicKey))
	providedBlock, _ := pem.Decode([]byte(publicKey))

	if extractedBlock == nil || providedBlock == nil {
		return false, fmt.Errorf("failed to parse one or both PEM blocks")
	}

	// DER形式で比較
	if string(extractedBlock.Bytes) != string(providedBlock.Bytes) {
		return false, nil
	}

	return true, nil
}

// GeneratePrivateKey は指定されたアルゴリズムに応じて秘密鍵を生成する
func GeneratePrivateKey(algorithm string) (string, error) {
	// サポートされているアルゴリズムかチェック
	if !isAsymmetricAlgorithm(algorithm) {
		return "", fmt.Errorf("algorithm %s does not support key generation (only RSA, ECDSA, and EdDSA are supported)", algorithm)
	}

	var privateKeyDER []byte

	// アルゴリズムに応じて秘密鍵を生成
	switch {
	case strings.HasPrefix(algorithm, "RS") || strings.HasPrefix(algorithm, "PS"):
		// RSA鍵を生成（2048ビット）
		rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return "", fmt.Errorf("failed to generate RSA key: %w", err)
		}
		// PKCS8形式でエンコード
		privateKeyDER, err = x509.MarshalPKCS8PrivateKey(rsaKey)
		if err != nil {
			return "", fmt.Errorf("failed to marshal RSA private key: %w", err)
		}

	case strings.HasPrefix(algorithm, "ES"):
		// ECDSA鍵を生成
		var curve elliptic.Curve
		switch algorithm {
		case "ES256":
			curve = elliptic.P256()
		case "ES384":
			curve = elliptic.P384()
		case "ES512":
			curve = elliptic.P521() // P-521 for ES512
		default:
			return "", fmt.Errorf("unsupported ECDSA algorithm: %s", algorithm)
		}

		ecdsaKey, err := ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			return "", fmt.Errorf("failed to generate ECDSA key: %w", err)
		}
		// PKCS8形式でエンコード
		privateKeyDER, err = x509.MarshalPKCS8PrivateKey(ecdsaKey)
		if err != nil {
			return "", fmt.Errorf("failed to marshal ECDSA private key: %w", err)
		}

	case algorithm == "EdDSA":
		// Ed25519鍵を生成
		_, ed25519Key, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return "", fmt.Errorf("failed to generate Ed25519 key: %w", err)
		}
		// PKCS8形式でエンコード
		privateKeyDER, err = x509.MarshalPKCS8PrivateKey(ed25519Key)
		if err != nil {
			return "", fmt.Errorf("failed to marshal Ed25519 private key: %w", err)
		}

	default:
		return "", fmt.Errorf("unexpected algorithm: %s", algorithm)
	}

	// PEM形式にエンコード
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyDER,
	})

	return string(privateKeyPEM), nil
}

// GenerateSelfSignedCertificate は自己署名証明書と秘密鍵をPEM形式で生成する
func GenerateSelfSignedCertificate(commonName string, organization string, sanEntries []string, validDays int, algorithm string) (SelfSignedCertificateResult, error) {
	commonName = strings.TrimSpace(commonName)
	organization = strings.TrimSpace(organization)
	algorithm = strings.TrimSpace(algorithm)

	if commonName == "" {
		return SelfSignedCertificateResult{}, fmt.Errorf("common name cannot be empty")
	}
	if validDays <= 0 {
		return SelfSignedCertificateResult{}, fmt.Errorf("valid days must be greater than zero")
	}
	if algorithm == "" {
		algorithm = "RS256"
	}

	var privateKey interface{}
	switch {
	case strings.HasPrefix(algorithm, "RS") || strings.HasPrefix(algorithm, "PS"):
		rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			return SelfSignedCertificateResult{}, fmt.Errorf("failed to generate RSA key: %w", err)
		}
		privateKey = rsaKey
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
			return SelfSignedCertificateResult{}, fmt.Errorf("unsupported ECDSA algorithm: %s", algorithm)
		}

		ecdsaKey, err := ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			return SelfSignedCertificateResult{}, fmt.Errorf("failed to generate ECDSA key: %w", err)
		}
		privateKey = ecdsaKey
	case algorithm == "EdDSA":
		_, ed25519Key, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return SelfSignedCertificateResult{}, fmt.Errorf("failed to generate Ed25519 key: %w", err)
		}
		privateKey = ed25519Key
	default:
		return SelfSignedCertificateResult{}, fmt.Errorf("unsupported algorithm: %s", algorithm)
	}

	serialLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialLimit)
	if err != nil {
		return SelfSignedCertificateResult{}, fmt.Errorf("failed to generate serial number: %w", err)
	}

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
	if len(dnsNames) == 0 && len(ipAddresses) == 0 {
		dnsNames = append(dnsNames, commonName)
	}

	notBefore := time.Now().Add(-5 * time.Minute)
	notAfter := notBefore.Add(time.Duration(validDays) * 24 * time.Hour)

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: commonName,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		DNSNames:              dnsNames,
		IPAddresses:           ipAddresses,
	}
	if organization != "" {
		template.Subject.Organization = []string{organization}
	}

	var publicKey interface{}
	switch key := privateKey.(type) {
	case *rsa.PrivateKey:
		publicKey = &key.PublicKey
	case *ecdsa.PrivateKey:
		publicKey = &key.PublicKey
	case ed25519.PrivateKey:
		publicKey = key.Public()
	default:
		return SelfSignedCertificateResult{}, fmt.Errorf("unsupported private key type: %T", privateKey)
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, publicKey, privateKey)
	if err != nil {
		return SelfSignedCertificateResult{}, fmt.Errorf("failed to create certificate: %w", err)
	}

	privateKeyDER, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return SelfSignedCertificateResult{}, fmt.Errorf("failed to marshal private key: %w", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privateKeyDER})

	return SelfSignedCertificateResult{
		CertificatePEM: string(certPEM),
		PrivateKeyPEM:  string(privateKeyPEM),
	}, nil
}

// GenerateMTLSCertificates は mTLS 用の CA、サーバー、クライアント証明書を生成する
func GenerateMTLSCertificates(
	caCommonName string,
	serverCommonName string,
	clientCommonName string,
	organization string,
	serverSanEntries []string,
	clientSanEntries []string,
	validDays int,
	algorithm string,
) (MTLSCertificateResult, error) {
	caCommonName = strings.TrimSpace(caCommonName)
	serverCommonName = strings.TrimSpace(serverCommonName)
	clientCommonName = strings.TrimSpace(clientCommonName)
	organization = strings.TrimSpace(organization)
	algorithm = strings.TrimSpace(algorithm)

	if caCommonName == "" {
		return MTLSCertificateResult{}, fmt.Errorf("CA common name cannot be empty")
	}
	if serverCommonName == "" {
		return MTLSCertificateResult{}, fmt.Errorf("server common name cannot be empty")
	}
	if clientCommonName == "" {
		return MTLSCertificateResult{}, fmt.Errorf("client common name cannot be empty")
	}
	if validDays <= 0 {
		return MTLSCertificateResult{}, fmt.Errorf("valid days must be greater than zero")
	}
	if algorithm == "" {
		algorithm = "RS256"
	}

	// 1. CA 証明書と秘密鍵を生成
	caPrivateKey, err := generatePrivateKeyForAlgorithm(algorithm)
	if err != nil {
		return MTLSCertificateResult{}, fmt.Errorf("failed to generate CA private key: %w", err)
	}

	caSerialLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	caSerialNumber, err := rand.Int(rand.Reader, caSerialLimit)
	if err != nil {
		return MTLSCertificateResult{}, fmt.Errorf("failed to generate CA serial number: %w", err)
	}

	notBefore := time.Now().Add(-5 * time.Minute)
	notAfter := notBefore.Add(time.Duration(validDays) * 24 * time.Hour)

	caTemplate := &x509.Certificate{
		SerialNumber: caSerialNumber,
		Subject: pkix.Name{
			CommonName: caCommonName,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		BasicConstraintsValid: true,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		MaxPathLen:            0,
	}
	if organization != "" {
		caTemplate.Subject.Organization = []string{organization}
	}

	caPublicKey := getPublicKeyFromPrivate(caPrivateKey)
	caCertDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, caPublicKey, caPrivateKey)
	if err != nil {
		return MTLSCertificateResult{}, fmt.Errorf("failed to create CA certificate: %w", err)
	}

	caPrivateKeyDER, err := x509.MarshalPKCS8PrivateKey(caPrivateKey)
	if err != nil {
		return MTLSCertificateResult{}, fmt.Errorf("failed to marshal CA private key: %w", err)
	}

	caCertPEM := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caCertDER}))
	caPrivateKeyPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: caPrivateKeyDER}))

	// CA 証明書をパースして署名に使用
	caCert, err := x509.ParseCertificate(caCertDER)
	if err != nil {
		return MTLSCertificateResult{}, fmt.Errorf("failed to parse CA certificate: %w", err)
	}

	// 2. サーバー証明書と秘密鍵を生成
	serverPrivateKey, err := generatePrivateKeyForAlgorithm(algorithm)
	if err != nil {
		return MTLSCertificateResult{}, fmt.Errorf("failed to generate server private key: %w", err)
	}

	serverSerialNumber, err := rand.Int(rand.Reader, caSerialLimit)
	if err != nil {
		return MTLSCertificateResult{}, fmt.Errorf("failed to generate server serial number: %w", err)
	}

	serverDNSNames, serverIPAddresses := parseSANEntries(serverSanEntries, serverCommonName)

	serverTemplate := &x509.Certificate{
		SerialNumber: serverSerialNumber,
		Subject: pkix.Name{
			CommonName: serverCommonName,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:              serverDNSNames,
		IPAddresses:           serverIPAddresses,
	}
	if organization != "" {
		serverTemplate.Subject.Organization = []string{organization}
	}

	serverPublicKey := getPublicKeyFromPrivate(serverPrivateKey)
	serverCertDER, err := x509.CreateCertificate(rand.Reader, serverTemplate, caCert, serverPublicKey, caPrivateKey)
	if err != nil {
		return MTLSCertificateResult{}, fmt.Errorf("failed to create server certificate: %w", err)
	}

	serverPrivateKeyDER, err := x509.MarshalPKCS8PrivateKey(serverPrivateKey)
	if err != nil {
		return MTLSCertificateResult{}, fmt.Errorf("failed to marshal server private key: %w", err)
	}

	serverCertPEM := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: serverCertDER}))
	serverPrivateKeyPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: serverPrivateKeyDER}))

	// 3. クライアント証明書と秘密鍵を生成
	clientPrivateKey, err := generatePrivateKeyForAlgorithm(algorithm)
	if err != nil {
		return MTLSCertificateResult{}, fmt.Errorf("failed to generate client private key: %w", err)
	}

	clientSerialNumber, err := rand.Int(rand.Reader, caSerialLimit)
	if err != nil {
		return MTLSCertificateResult{}, fmt.Errorf("failed to generate client serial number: %w", err)
	}

	clientDNSNames, clientIPAddresses := parseSANEntries(clientSanEntries, clientCommonName)

	clientTemplate := &x509.Certificate{
		SerialNumber: clientSerialNumber,
		Subject: pkix.Name{
			CommonName: clientCommonName,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		DNSNames:              clientDNSNames,
		IPAddresses:           clientIPAddresses,
	}
	if organization != "" {
		clientTemplate.Subject.Organization = []string{organization}
	}

	clientPublicKey := getPublicKeyFromPrivate(clientPrivateKey)
	clientCertDER, err := x509.CreateCertificate(rand.Reader, clientTemplate, caCert, clientPublicKey, caPrivateKey)
	if err != nil {
		return MTLSCertificateResult{}, fmt.Errorf("failed to create client certificate: %w", err)
	}

	clientPrivateKeyDER, err := x509.MarshalPKCS8PrivateKey(clientPrivateKey)
	if err != nil {
		return MTLSCertificateResult{}, fmt.Errorf("failed to marshal client private key: %w", err)
	}

	clientCertPEM := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: clientCertDER}))
	clientPrivateKeyPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: clientPrivateKeyDER}))

	return MTLSCertificateResult{
		CACertificatePEM:     caCertPEM,
		CAPrivateKeyPEM:      caPrivateKeyPEM,
		ServerCertificatePEM: serverCertPEM,
		ServerPrivateKeyPEM:  serverPrivateKeyPEM,
		ClientCertificatePEM: clientCertPEM,
		ClientPrivateKeyPEM:  clientPrivateKeyPEM,
	}, nil
}

// GenerateMTLSCertificatesMultiClient は複数のクライアント証明書を生成する
func GenerateMTLSCertificatesMultiClient(
	caCommonName string,
	serverCommonName string,
	clientCommonNamePrefix string,
	clientCount int,
	organization string,
	serverSanEntries []string,
	clientSanEntries []string,
	validDays int,
	algorithm string,
) (MTLSCertificatesMultiClientResult, error) {
	caCommonName = strings.TrimSpace(caCommonName)
	serverCommonName = strings.TrimSpace(serverCommonName)
	clientCommonNamePrefix = strings.TrimSpace(clientCommonNamePrefix)
	organization = strings.TrimSpace(organization)
	algorithm = strings.TrimSpace(algorithm)

	if caCommonName == "" {
		return MTLSCertificatesMultiClientResult{}, fmt.Errorf("CA common name cannot be empty")
	}
	if serverCommonName == "" {
		return MTLSCertificatesMultiClientResult{}, fmt.Errorf("server common name cannot be empty")
	}
	if clientCommonNamePrefix == "" {
		return MTLSCertificatesMultiClientResult{}, fmt.Errorf("client common name prefix cannot be empty")
	}
	if clientCount <= 0 || clientCount > 100 {
		return MTLSCertificatesMultiClientResult{}, fmt.Errorf("client count must be between 1 and 100")
	}
	if validDays <= 0 {
		return MTLSCertificatesMultiClientResult{}, fmt.Errorf("valid days must be greater than zero")
	}
	if algorithm == "" {
		algorithm = "RS256"
	}

	// 1. CA 証明書と秘密鍵を生成
	caPrivateKey, err := generatePrivateKeyForAlgorithm(algorithm)
	if err != nil {
		return MTLSCertificatesMultiClientResult{}, fmt.Errorf("failed to generate CA private key: %w", err)
	}

	caSerialLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	caSerialNumber, err := rand.Int(rand.Reader, caSerialLimit)
	if err != nil {
		return MTLSCertificatesMultiClientResult{}, fmt.Errorf("failed to generate CA serial number: %w", err)
	}

	notBefore := time.Now().Add(-5 * time.Minute)
	notAfter := notBefore.Add(time.Duration(validDays) * 24 * time.Hour)

	caTemplate := &x509.Certificate{
		SerialNumber: caSerialNumber,
		Subject: pkix.Name{
			CommonName: caCommonName,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		BasicConstraintsValid: true,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		MaxPathLen:            0,
	}
	if organization != "" {
		caTemplate.Subject.Organization = []string{organization}
	}

	caPublicKey := getPublicKeyFromPrivate(caPrivateKey)
	caCertDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, caPublicKey, caPrivateKey)
	if err != nil {
		return MTLSCertificatesMultiClientResult{}, fmt.Errorf("failed to create CA certificate: %w", err)
	}

	caPrivateKeyDER, err := x509.MarshalPKCS8PrivateKey(caPrivateKey)
	if err != nil {
		return MTLSCertificatesMultiClientResult{}, fmt.Errorf("failed to marshal CA private key: %w", err)
	}

	caCertPEM := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caCertDER}))
	caPrivateKeyPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: caPrivateKeyDER}))

	// CA 証明書をパースして署名に使用
	caCert, err := x509.ParseCertificate(caCertDER)
	if err != nil {
		return MTLSCertificatesMultiClientResult{}, fmt.Errorf("failed to parse CA certificate: %w", err)
	}

	// 1.5 CRL を生成
	caSigner, err := getSignerFromPrivateKey(caPrivateKey)
	if err != nil {
		return MTLSCertificatesMultiClientResult{}, fmt.Errorf("failed to build CRL signer: %w", err)
	}

	crlDER, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{
		Number:     big.NewInt(1),
		ThisUpdate: notBefore,
		NextUpdate: notAfter,
	}, caCert, caSigner)
	if err != nil {
		return MTLSCertificatesMultiClientResult{}, fmt.Errorf("failed to create CRL: %w", err)
	}
	crlPEM := string(pem.EncodeToMemory(&pem.Block{Type: "X509 CRL", Bytes: crlDER}))

	// 2. サーバー証明書と秘密鍵を生成
	serverPrivateKey, err := generatePrivateKeyForAlgorithm(algorithm)
	if err != nil {
		return MTLSCertificatesMultiClientResult{}, fmt.Errorf("failed to generate server private key: %w", err)
	}

	serverSerialNumber, err := rand.Int(rand.Reader, caSerialLimit)
	if err != nil {
		return MTLSCertificatesMultiClientResult{}, fmt.Errorf("failed to generate server serial number: %w", err)
	}

	serverDNSNames, serverIPAddresses := parseSANEntries(serverSanEntries, serverCommonName)

	serverTemplate := &x509.Certificate{
		SerialNumber: serverSerialNumber,
		Subject: pkix.Name{
			CommonName: serverCommonName,
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:              serverDNSNames,
		IPAddresses:           serverIPAddresses,
	}
	if organization != "" {
		serverTemplate.Subject.Organization = []string{organization}
	}

	serverPublicKey := getPublicKeyFromPrivate(serverPrivateKey)
	serverCertDER, err := x509.CreateCertificate(rand.Reader, serverTemplate, caCert, serverPublicKey, caPrivateKey)
	if err != nil {
		return MTLSCertificatesMultiClientResult{}, fmt.Errorf("failed to create server certificate: %w", err)
	}

	serverPrivateKeyDER, err := x509.MarshalPKCS8PrivateKey(serverPrivateKey)
	if err != nil {
		return MTLSCertificatesMultiClientResult{}, fmt.Errorf("failed to marshal server private key: %w", err)
	}

	serverCertPEM := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: serverCertDER}))
	serverPrivateKeyPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: serverPrivateKeyDER}))

	// 3. 複数のクライアント証明書と秘密鍵を生成
	clientCerts := make([]ClientCertificate, 0, clientCount)

	for i := 0; i < clientCount; i++ {
		clientCommonName := fmt.Sprintf("%s-%d", clientCommonNamePrefix, i+1)

		clientPrivateKey, err := generatePrivateKeyForAlgorithm(algorithm)
		if err != nil {
			return MTLSCertificatesMultiClientResult{}, fmt.Errorf("failed to generate client %d private key: %w", i+1, err)
		}

		clientSerialNumber, err := rand.Int(rand.Reader, caSerialLimit)
		if err != nil {
			return MTLSCertificatesMultiClientResult{}, fmt.Errorf("failed to generate client %d serial number: %w", i+1, err)
		}

		clientDNSNames, clientIPAddresses := parseSANEntries(clientSanEntries, clientCommonName)

		clientTemplate := &x509.Certificate{
			SerialNumber: clientSerialNumber,
			Subject: pkix.Name{
				CommonName: clientCommonName,
			},
			NotBefore:             notBefore,
			NotAfter:              notAfter,
			BasicConstraintsValid: true,
			KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			DNSNames:              clientDNSNames,
			IPAddresses:           clientIPAddresses,
		}
		if organization != "" {
			clientTemplate.Subject.Organization = []string{organization}
		}

		clientPublicKey := getPublicKeyFromPrivate(clientPrivateKey)
		clientCertDER, err := x509.CreateCertificate(rand.Reader, clientTemplate, caCert, clientPublicKey, caPrivateKey)
		if err != nil {
			return MTLSCertificatesMultiClientResult{}, fmt.Errorf("failed to create client %d certificate: %w", i+1, err)
		}

		clientPrivateKeyDER, err := x509.MarshalPKCS8PrivateKey(clientPrivateKey)
		if err != nil {
			return MTLSCertificatesMultiClientResult{}, fmt.Errorf("failed to marshal client %d private key: %w", i+1, err)
		}

		clientCertPEM := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: clientCertDER}))
		clientPrivateKeyPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: clientPrivateKeyDER}))

		clientCerts = append(clientCerts, ClientCertificate{
			CommonName:     clientCommonName,
			CertificatePEM: clientCertPEM,
			PrivateKeyPEM:  clientPrivateKeyPEM,
		})
	}

	return MTLSCertificatesMultiClientResult{
		CACertificatePEM:     caCertPEM,
		CAPrivateKeyPEM:      caPrivateKeyPEM,
		CRLPEM:               crlPEM,
		ServerCertificatePEM: serverCertPEM,
		ServerPrivateKeyPEM:  serverPrivateKeyPEM,
		ClientCertificates:   clientCerts,
	}, nil
}

// AddCertificatesToCRL は既存 CRL に失効対象証明書を追加して再発行する
func AddCertificatesToCRL(
	caCertificatePEM string,
	caPrivateKeyPEM string,
	existingCRLPEM string,
	revokedCertificatePEMs []string,
	nextUpdateDays int,
) (CRLUpdateResult, error) {
	caCertificatePEM = strings.TrimSpace(caCertificatePEM)
	caPrivateKeyPEM = strings.TrimSpace(caPrivateKeyPEM)
	existingCRLPEM = strings.TrimSpace(existingCRLPEM)

	if caCertificatePEM == "" {
		return CRLUpdateResult{}, fmt.Errorf("CA certificate cannot be empty")
	}
	if caPrivateKeyPEM == "" {
		return CRLUpdateResult{}, fmt.Errorf("CA private key cannot be empty")
	}
	if len(revokedCertificatePEMs) == 0 {
		return CRLUpdateResult{}, fmt.Errorf("at least one certificate must be provided")
	}
	if nextUpdateDays <= 0 {
		nextUpdateDays = 30
	}

	caCert, err := parseCertificateFromPEM(caCertificatePEM)
	if err != nil {
		return CRLUpdateResult{}, fmt.Errorf("failed to parse CA certificate: %w", err)
	}

	caPrivateKey, err := parsePrivateKeyFromPEM(caPrivateKeyPEM)
	if err != nil {
		return CRLUpdateResult{}, fmt.Errorf("failed to parse CA private key: %w", err)
	}

	caSigner, err := getSignerFromPrivateKey(caPrivateKey)
	if err != nil {
		return CRLUpdateResult{}, fmt.Errorf("failed to build CRL signer: %w", err)
	}

	now := time.Now().UTC()
	entries := make([]x509.RevocationListEntry, 0)
	existingSerials := make(map[string]struct{})
	crlNumber := big.NewInt(1)

	if existingCRLPEM != "" {
		existingCRL, err := parseCRLFromPEM(existingCRLPEM)
		if err != nil {
			return CRLUpdateResult{}, fmt.Errorf("failed to parse existing CRL: %w", err)
		}

		for _, entry := range existingCRL.RevokedCertificateEntries {
			if entry.SerialNumber == nil {
				continue
			}
			serialHex := entry.SerialNumber.Text(16)
			if _, exists := existingSerials[serialHex]; exists {
				continue
			}
			existingSerials[serialHex] = struct{}{}
			entries = append(entries, x509.RevocationListEntry{
				SerialNumber:   new(big.Int).Set(entry.SerialNumber),
				RevocationTime: entry.RevocationTime,
			})
		}

		if existingCRL.Number != nil {
			crlNumber = new(big.Int).Add(existingCRL.Number, big.NewInt(1))
		}
	}

	addedCount := 0
	for _, certPEM := range revokedCertificatePEMs {
		trimmed := strings.TrimSpace(certPEM)
		if trimmed == "" {
			continue
		}

		revokedCert, err := parseCertificateFromPEM(trimmed)
		if err != nil {
			return CRLUpdateResult{}, fmt.Errorf("failed to parse revoked certificate: %w", err)
		}

		serialHex := revokedCert.SerialNumber.Text(16)
		if _, exists := existingSerials[serialHex]; exists {
			continue
		}

		existingSerials[serialHex] = struct{}{}
		entries = append(entries, x509.RevocationListEntry{
			SerialNumber:   new(big.Int).Set(revokedCert.SerialNumber),
			RevocationTime: now,
		})
		addedCount++
	}

	if len(entries) == 0 {
		return CRLUpdateResult{}, fmt.Errorf("no revoked certificate entries found")
	}

	crlDER, err := x509.CreateRevocationList(rand.Reader, &x509.RevocationList{
		Number:                    crlNumber,
		ThisUpdate:                now,
		NextUpdate:                now.Add(time.Duration(nextUpdateDays) * 24 * time.Hour),
		RevokedCertificateEntries: entries,
	}, caCert, caSigner)
	if err != nil {
		return CRLUpdateResult{}, fmt.Errorf("failed to create CRL: %w", err)
	}

	serialNumbers := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.SerialNumber == nil {
			continue
		}
		serialNumbers = append(serialNumbers, strings.ToUpper(entry.SerialNumber.Text(16)))
	}
	sort.Strings(serialNumbers)

	return CRLUpdateResult{
		CRLPEM:               string(pem.EncodeToMemory(&pem.Block{Type: "X509 CRL", Bytes: crlDER})),
		AddedCount:           addedCount,
		TotalRevokedCount:    len(entries),
		RevokedSerialNumbers: serialNumbers,
	}, nil
}

// SaveTextFile は指定内容を Downloads/oreno-tools-certs 配下に保存する
func SaveTextFile(content string, filename string) (string, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return "", fmt.Errorf("file content cannot be empty")
	}

	safeFilename := sanitizeFilename(filename)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to resolve home directory: %w", err)
	}

	targetDir := filepath.Join(homeDir, "Downloads", "oreno-tools-certs")
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	targetPath := filepath.Join(targetDir, safeFilename)
	if err := os.WriteFile(targetPath, []byte(content+"\n"), 0o600); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return targetPath, nil
}

func parsePrivateKeyFromPEM(privateKeyPEM string) (interface{}, error) {
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, fmt.Errorf("invalid PEM data")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err == nil {
		return privateKey, nil
	}

	if rsaPrivateKey, rsaErr := x509.ParsePKCS1PrivateKey(block.Bytes); rsaErr == nil {
		return rsaPrivateKey, nil
	}

	if ecdsaPrivateKey, ecErr := x509.ParseECPrivateKey(block.Bytes); ecErr == nil {
		return ecdsaPrivateKey, nil
	}

	return nil, fmt.Errorf("unsupported private key format")
}

func parseCRLFromPEM(crlPEM string) (*x509.RevocationList, error) {
	block, _ := pem.Decode([]byte(crlPEM))
	if block == nil {
		return nil, fmt.Errorf("invalid PEM data")
	}
	if block.Type != "X509 CRL" && block.Type != "CERTIFICATE REVOCATION LIST" {
		return nil, fmt.Errorf("invalid PEM block type: %s", block.Type)
	}

	crl, err := x509.ParseRevocationList(block.Bytes)
	if err != nil {
		return nil, err
	}

	return crl, nil
}

func sanitizeFilename(filename string) string {
	trimmed := strings.TrimSpace(filename)
	if trimmed == "" {
		return "certificate.pem"
	}

	base := filepath.Base(trimmed)
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
	)
	safe := replacer.Replace(base)
	if safe == "" || safe == "." || safe == ".." {
		return "certificate.pem"
	}

	return safe
}
