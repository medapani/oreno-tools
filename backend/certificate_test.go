package backend

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"net"
	"strings"
	"testing"
)

func TestGeneratePrivateKey(t *testing.T) {
	tests := []struct {
		name      string
		algorithm string
		wantErr   bool
	}{
		{"EdDSA", "EdDSA", false},
		{"ES256", "ES256", false},
		{"ES384", "ES384", false},
		{"ES512", "ES512", false},
		{"RS256", "RS256", false},
		{"PS256", "PS256", false},
		{"不正なアルゴリズム (HMAC)", "HS256", true},
		{"不正なアルゴリズム (unknown)", "UNKNOWN", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GeneratePrivateKey(tt.algorithm)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePrivateKey(%q) error = %v, wantErr %v", tt.algorithm, err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if !strings.Contains(got, "-----BEGIN PRIVATE KEY-----") {
				t.Errorf("GeneratePrivateKey(%q) did not return PEM: %q", tt.algorithm, got)
			}
			if !strings.Contains(got, "-----END PRIVATE KEY-----") {
				t.Errorf("GeneratePrivateKey(%q) PEM not terminated: %q", tt.algorithm, got)
			}
		})
	}
}

func TestExtractPublicKey(t *testing.T) {
	algorithms := []string{"EdDSA", "ES256", "RS256"}

	for _, alg := range algorithms {
		t.Run(alg, func(t *testing.T) {
			privateKeyPEM, err := GeneratePrivateKey(alg)
			if err != nil {
				t.Fatalf("GeneratePrivateKey(%q) failed: %v", alg, err)
			}

			publicKeyPEM, err := ExtractPublicKey(privateKeyPEM, alg)
			if err != nil {
				t.Errorf("ExtractPublicKey(%q) error = %v", alg, err)
				return
			}

			if !strings.Contains(publicKeyPEM, "-----BEGIN PUBLIC KEY-----") {
				t.Errorf("ExtractPublicKey(%q) did not return public key PEM", alg)
			}
		})
	}
}

func TestExtractPublicKeyErrors(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		algorithm string
		wantErr   bool
	}{
		{
			name:      "空の秘密鍵",
			key:       "",
			algorithm: "RS256",
			wantErr:   true,
		},
		{
			name:      "HMACアルゴリズム (非対称鍵なし)",
			key:       "somesecret",
			algorithm: "HS256",
			wantErr:   true,
		},
		{
			name:      "不正なPEM",
			key:       "not a pem",
			algorithm: "RS256",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ExtractPublicKey(tt.key, tt.algorithm)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractPublicKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVerifyKeyPair(t *testing.T) {
	// EdDSA は生成が速いのでテスト用に使用
	alg := "EdDSA"
	privateKeyPEM, err := GeneratePrivateKey(alg)
	if err != nil {
		t.Fatalf("GeneratePrivateKey failed: %v", err)
	}
	publicKeyPEM, err := ExtractPublicKey(privateKeyPEM, alg)
	if err != nil {
		t.Fatalf("ExtractPublicKey failed: %v", err)
	}

	// 別のキーペアを生成（不一致チェック用）
	otherPrivateKeyPEM, err := GeneratePrivateKey(alg)
	if err != nil {
		t.Fatalf("GeneratePrivateKey (other) failed: %v", err)
	}
	otherPublicKeyPEM, err := ExtractPublicKey(otherPrivateKeyPEM, alg)
	if err != nil {
		t.Fatalf("ExtractPublicKey (other) failed: %v", err)
	}

	tests := []struct {
		name       string
		privateKey string
		publicKey  string
		algorithm  string
		wantMatch  bool
		wantErr    bool
	}{
		{
			name:       "一致するキーペア",
			privateKey: privateKeyPEM,
			publicKey:  publicKeyPEM,
			algorithm:  alg,
			wantMatch:  true,
			wantErr:    false,
		},
		{
			name:       "不一致のキーペア",
			privateKey: privateKeyPEM,
			publicKey:  otherPublicKeyPEM,
			algorithm:  alg,
			wantMatch:  false,
			wantErr:    false,
		},
		{
			name:       "空の秘密鍵",
			privateKey: "",
			publicKey:  publicKeyPEM,
			algorithm:  alg,
			wantErr:    true,
		},
		{
			name:       "空の公開鍵",
			privateKey: privateKeyPEM,
			publicKey:  "",
			algorithm:  alg,
			wantErr:    true,
		},
		{
			name:       "HMACアルゴリズム",
			privateKey: privateKeyPEM,
			publicKey:  publicKeyPEM,
			algorithm:  "HS256",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VerifyKeyPair(tt.privateKey, tt.publicKey, tt.algorithm)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyKeyPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got != tt.wantMatch {
				t.Errorf("VerifyKeyPair() = %v, want %v", got, tt.wantMatch)
			}
		})
	}
}

func TestGenerateSelfSignedCertificate(t *testing.T) {
	tests := []struct {
		name         string
		commonName   string
		organization string
		sanEntries   []string
		validDays    int
		algorithm    string
		wantErr      bool
	}{
		{
			name:         "EdDSA 基本",
			commonName:   "test.example.com",
			organization: "Test Org",
			sanEntries:   []string{"test.example.com"},
			validDays:    365,
			algorithm:    "EdDSA",
			wantErr:      false,
		},
		{
			name:         "ES256 SAN複数",
			commonName:   "example.com",
			organization: "Example Inc",
			sanEntries:   []string{"example.com", "www.example.com", "192.168.1.1"},
			validDays:    90,
			algorithm:    "ES256",
			wantErr:      false,
		},
		{
			name:         "空のCommonName はエラー",
			commonName:   "",
			organization: "Test Org",
			sanEntries:   nil,
			validDays:    365,
			algorithm:    "RS256",
			wantErr:      true,
		},
		{
			name:         "validDays = 0 はエラー",
			commonName:   "example.com",
			organization: "Test Org",
			sanEntries:   nil,
			validDays:    0,
			algorithm:    "EdDSA",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateSelfSignedCertificate(tt.commonName, tt.organization, tt.sanEntries, tt.validDays, tt.algorithm)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateSelfSignedCertificate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if !strings.Contains(got.CertificatePEM, "-----BEGIN CERTIFICATE-----") {
				t.Errorf("CertificatePEM does not contain certificate header")
			}
			if !strings.Contains(got.PrivateKeyPEM, "-----BEGIN PRIVATE KEY-----") {
				t.Errorf("PrivateKeyPEM does not contain private key header")
			}
		})
	}
}

func TestParseSANEntries(t *testing.T) {
	tests := []struct {
		name           string
		sanEntries     []string
		defaultDNS     string
		wantDNSCount   int
		wantIPCount    int
		wantDNSContain string
	}{
		{
			name:         "DNSのみ",
			sanEntries:   []string{"example.com", "www.example.com"},
			defaultDNS:   "",
			wantDNSCount: 2,
			wantIPCount:  0,
		},
		{
			name:         "IPのみ",
			sanEntries:   []string{"192.168.1.1", "10.0.0.1"},
			defaultDNS:   "",
			wantDNSCount: 0,
			wantIPCount:  2,
		},
		{
			name:         "DNSとIPの混在",
			sanEntries:   []string{"example.com", "192.168.1.1"},
			defaultDNS:   "",
			wantDNSCount: 1,
			wantIPCount:  1,
		},
		{
			name:           "空の場合はデフォルトDNSを使用",
			sanEntries:     []string{},
			defaultDNS:     "default.example.com",
			wantDNSCount:   1,
			wantIPCount:    0,
			wantDNSContain: "default.example.com",
		},
		{
			name:         "空エントリを無視",
			sanEntries:   []string{"example.com", "", "  "},
			defaultDNS:   "",
			wantDNSCount: 1,
			wantIPCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dnsNames, ipAddresses := parseSANEntries(tt.sanEntries, tt.defaultDNS)
			if len(dnsNames) != tt.wantDNSCount {
				t.Errorf("DNS count = %d, want %d", len(dnsNames), tt.wantDNSCount)
			}
			if len(ipAddresses) != tt.wantIPCount {
				t.Errorf("IP count = %d, want %d", len(ipAddresses), tt.wantIPCount)
			}
			if tt.wantDNSContain != "" {
				found := false
				for _, d := range dnsNames {
					if d == tt.wantDNSContain {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("DNS names %v did not contain %q", dnsNames, tt.wantDNSContain)
				}
			}
		})
	}
}

// IPアドレスとして正しくパースされることを確認
func TestParseSANEntriesIPv6(t *testing.T) {
	sanEntries := []string{"::1", "2001:db8::1"}
	dnsNames, ipAddresses := parseSANEntries(sanEntries, "")
	if len(dnsNames) != 0 {
		t.Errorf("expected 0 DNS names, got %d", len(dnsNames))
	}
	if len(ipAddresses) != 2 {
		t.Errorf("expected 2 IP addresses, got %d", len(ipAddresses))
	}
	// IPv6アドレスが正しくパースされているか確認
	for _, ip := range ipAddresses {
		if ip.To16() == nil {
			t.Errorf("expected valid IP, got %v", ip)
		}
	}
}

// IPアドレスかどうかの検証（net.ParseIP を使用）
func TestNetParseIP(t *testing.T) {
	if ip := net.ParseIP("192.168.1.1"); ip == nil {
		t.Error("192.168.1.1 should be a valid IP")
	}
	if ip := net.ParseIP("not-an-ip"); ip != nil {
		t.Error("not-an-ip should not be a valid IP")
	}
}

// ---- sanitizeFilename ----

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"通常のファイル名", "server.pem", "server.pem"},
		{"スラッシュを置換", "path/to/file.pem", "file.pem"},
		{"バックスラッシュを置換", `path\file.pem`, "path_file.pem"},
		{"コロンを置換", "file:name.pem", "file_name.pem"},
		{"アスタリスクを置換", "file*.pem", "file_.pem"},
		{"空文字列はデフォルト名", "", "certificate.pem"},
		{"スペースのみはデフォルト名", "   ", "certificate.pem"},
		{"ドットのみはデフォルト名", ".", "certificate.pem"},
		{"ダブルドットはデフォルト名", "..", "certificate.pem"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeFilename(tt.input)
			if got != tt.want {
				t.Errorf("sanitizeFilename(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// ---- parsePrivateKeyFromPEM ----

func TestParsePrivateKeyFromPEM(t *testing.T) {
	tests := []struct {
		name      string
		algorithm string
		wantErr   bool
	}{
		{"EdDSA", "EdDSA", false},
		{"ES256", "ES256", false},
		{"RS256", "RS256", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keyPEM, err := GeneratePrivateKey(tt.algorithm)
			if err != nil {
				t.Fatalf("GeneratePrivateKey failed: %v", err)
			}
			got, err := parsePrivateKeyFromPEM(keyPEM)
			if (err != nil) != tt.wantErr {
				t.Errorf("parsePrivateKeyFromPEM() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil {
				t.Error("parsePrivateKeyFromPEM() returned nil")
			}
		})
	}

	t.Run("不正なPEM", func(t *testing.T) {
		_, err := parsePrivateKeyFromPEM("not a pem")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("PKCS1形式RSA鍵", func(t *testing.T) {
		// GeneratePrivateKey は常に PKCS8 でエンコードするため、
		// PKCS1 パスは手動で PEM を作成してテストする
		rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatalf("rsa.GenerateKey failed: %v", err)
		}
		pkcs1DER := x509.MarshalPKCS1PrivateKey(rsaKey)
		pkcs1PEM := string(pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: pkcs1DER,
		}))
		got, err := parsePrivateKeyFromPEM(pkcs1PEM)
		if err != nil {
			t.Errorf("parsePrivateKeyFromPEM() PKCS1 error = %v", err)
			return
		}
		if got == nil {
			t.Error("parsePrivateKeyFromPEM() PKCS1 returned nil")
		}
	})

	t.Run("EC形式ECDSA鍵", func(t *testing.T) {
		// EC PRIVATE KEY 形式は ParsePKCS8PrivateKey / ParsePKCS1PrivateKey では
		// 解析できないため ParseECPrivateKey パスに到達する
		ecKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			t.Fatalf("ecdsa.GenerateKey failed: %v", err)
		}
		ecDER, err := x509.MarshalECPrivateKey(ecKey)
		if err != nil {
			t.Fatalf("MarshalECPrivateKey failed: %v", err)
		}
		ecPEM := string(pem.EncodeToMemory(&pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: ecDER,
		}))
		got, err := parsePrivateKeyFromPEM(ecPEM)
		if err != nil {
			t.Errorf("parsePrivateKeyFromPEM() EC error = %v", err)
			return
		}
		if got == nil {
			t.Error("parsePrivateKeyFromPEM() EC returned nil")
		}
	})

	t.Run("サポート外フォーマット", func(t *testing.T) {
		// 全パーサーが失敗する不正な DER バイト列を渡す
		unsupportedPEM := string(pem.EncodeToMemory(&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: []byte("invalid der bytes"),
		}))
		_, err := parsePrivateKeyFromPEM(unsupportedPEM)
		if err == nil {
			t.Error("expected error for unsupported format, got nil")
		}
		if !strings.Contains(err.Error(), "unsupported private key format") {
			t.Errorf("unexpected error message: %v", err)
		}
	})
}

// ---- parseCRLFromPEM ----

func TestParseCRLFromPEM(t *testing.T) {
	// EdDSA CA を生成して CRL を作成する
	result, err := GenerateMTLSCertificatesMultiClient(
		"Test CA", "server.example.com", "client", 1,
		"Test Org", nil, nil, 30, "EdDSA",
	)
	if err != nil {
		t.Fatalf("GenerateMTLSCertificatesMultiClient failed: %v", err)
	}

	t.Run("正常なCRL", func(t *testing.T) {
		crl, err := parseCRLFromPEM(result.CRLPEM)
		if err != nil {
			t.Errorf("parseCRLFromPEM() error = %v", err)
			return
		}
		if crl == nil {
			t.Error("parseCRLFromPEM() returned nil")
		}
	})

	t.Run("不正なPEM", func(t *testing.T) {
		_, err := parseCRLFromPEM("not a pem")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("証明書PEMを渡すとエラー", func(t *testing.T) {
		_, err := parseCRLFromPEM(result.CACertificatePEM)
		if err == nil {
			t.Error("expected error for non-CRL PEM, got nil")
		}
	})
}

// ---- GenerateMTLSCertificates ----

func TestGenerateMTLSCertificates(t *testing.T) {
	tests := []struct {
		name             string
		caCommonName     string
		serverCommonName string
		clientCommonName string
		organization     string
		serverSANs       []string
		clientSANs       []string
		validDays        int
		algorithm        string
		wantErr          bool
	}{
		{
			name:             "EdDSA 基本",
			caCommonName:     "Test CA",
			serverCommonName: "server.example.com",
			clientCommonName: "client.example.com",
			organization:     "Test Org",
			serverSANs:       []string{"server.example.com"},
			clientSANs:       []string{"client.example.com"},
			validDays:        30,
			algorithm:        "EdDSA",
			wantErr:          false,
		},
		{
			name:             "ES256 基本",
			caCommonName:     "Test CA",
			serverCommonName: "server.example.com",
			clientCommonName: "client.example.com",
			organization:     "Test Org",
			validDays:        30,
			algorithm:        "ES256",
			wantErr:          false,
		},
		{
			name:             "caCommonName 空はエラー",
			caCommonName:     "",
			serverCommonName: "server.example.com",
			clientCommonName: "client.example.com",
			validDays:        30,
			algorithm:        "EdDSA",
			wantErr:          true,
		},
		{
			name:             "serverCommonName 空はエラー",
			caCommonName:     "Test CA",
			serverCommonName: "",
			clientCommonName: "client.example.com",
			validDays:        30,
			algorithm:        "EdDSA",
			wantErr:          true,
		},
		{
			name:             "clientCommonName 空はエラー",
			caCommonName:     "Test CA",
			serverCommonName: "server.example.com",
			clientCommonName: "",
			validDays:        30,
			algorithm:        "EdDSA",
			wantErr:          true,
		},
		{
			name:             "validDays = 0 はエラー",
			caCommonName:     "Test CA",
			serverCommonName: "server.example.com",
			clientCommonName: "client.example.com",
			validDays:        0,
			algorithm:        "EdDSA",
			wantErr:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateMTLSCertificates(
				tt.caCommonName, tt.serverCommonName, tt.clientCommonName,
				tt.organization, tt.serverSANs, tt.clientSANs, tt.validDays, tt.algorithm,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateMTLSCertificates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			for _, pem := range []string{
				got.CACertificatePEM, got.CAPrivateKeyPEM,
				got.ServerCertificatePEM, got.ServerPrivateKeyPEM,
				got.ClientCertificatePEM, got.ClientPrivateKeyPEM,
			} {
				if pem == "" {
					t.Error("expected non-empty PEM field")
				}
			}
		})
	}
}

// CA で署名されたサーバー証明書が CA プールで検証できることを確認
func TestGenerateMTLSCertificatesVerify(t *testing.T) {
	result, err := GenerateMTLSCertificates(
		"Test CA", "server.example.com", "client.example.com",
		"Test Org", []string{"server.example.com"}, nil, 30, "EdDSA",
	)
	if err != nil {
		t.Fatalf("GenerateMTLSCertificates failed: %v", err)
	}

	caBlock, _ := pem.Decode([]byte(result.CACertificatePEM))
	caCert, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to parse CA cert: %v", err)
	}
	if !caCert.IsCA {
		t.Error("CA certificate IsCA should be true")
	}

	serverBlock, _ := pem.Decode([]byte(result.ServerCertificatePEM))
	serverCert, err := x509.ParseCertificate(serverBlock.Bytes)
	if err != nil {
		t.Fatalf("failed to parse server cert: %v", err)
	}

	pool := x509.NewCertPool()
	pool.AddCert(caCert)
	_, err = serverCert.Verify(x509.VerifyOptions{Roots: pool})
	if err != nil {
		t.Errorf("server cert verification failed: %v", err)
	}
}

// ---- GenerateMTLSCertificatesMultiClient ----

func TestGenerateMTLSCertificatesMultiClient(t *testing.T) {
	tests := []struct {
		name        string
		clientCount int
		algorithm   string
		wantErr     bool
	}{
		{"EdDSA 1クライアント", 1, "EdDSA", false},
		{"EdDSA 3クライアント", 3, "EdDSA", false},
		{"ES256 2クライアント", 2, "ES256", false},
		{"0クライアントはエラー", 0, "EdDSA", true},
		{"101クライアントはエラー", 101, "EdDSA", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateMTLSCertificatesMultiClient(
				"Test CA", "server.example.com", "client",
				tt.clientCount, "Test Org", nil, nil, 30, tt.algorithm,
			)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateMTLSCertificatesMultiClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if len(got.ClientCertificates) != tt.clientCount {
				t.Errorf("ClientCertificates count = %d, want %d", len(got.ClientCertificates), tt.clientCount)
			}
			if got.CRLPEM == "" {
				t.Error("CRLPEM should not be empty")
			}
			for i, cc := range got.ClientCertificates {
				expectedName := "client-" + strings.TrimPrefix(strings.TrimPrefix("000", ""), "")
				_ = expectedName
				if cc.CommonName == "" {
					t.Errorf("ClientCertificates[%d].CommonName is empty", i)
				}
				if cc.CertificatePEM == "" {
					t.Errorf("ClientCertificates[%d].CertificatePEM is empty", i)
				}
				if cc.PrivateKeyPEM == "" {
					t.Errorf("ClientCertificates[%d].PrivateKeyPEM is empty", i)
				}
			}
		})
	}
}

func TestGenerateMTLSCertificatesMultiClientValidation(t *testing.T) {
	base := func() (string, string, string, int, string, []string, []string, int, string) {
		return "CA", "server.example.com", "client", 1, "Org", nil, nil, 30, "EdDSA"
	}

	t.Run("caCommonName 空はエラー", func(t *testing.T) {
		_, err := GenerateMTLSCertificatesMultiClient("", "server.example.com", "client", 1, "", nil, nil, 30, "EdDSA")
		if err == nil {
			t.Error("expected error")
		}
	})
	t.Run("serverCommonName 空はエラー", func(t *testing.T) {
		ca, _, prefix, count, org, ss, cs, days, alg := base()
		_, err := GenerateMTLSCertificatesMultiClient(ca, "", prefix, count, org, ss, cs, days, alg)
		if err == nil {
			t.Error("expected error")
		}
	})
	t.Run("clientPrefix 空はエラー", func(t *testing.T) {
		ca, server, _, count, org, ss, cs, days, alg := base()
		_, err := GenerateMTLSCertificatesMultiClient(ca, server, "", count, org, ss, cs, days, alg)
		if err == nil {
			t.Error("expected error")
		}
	})
	t.Run("validDays = 0 はエラー", func(t *testing.T) {
		ca, server, prefix, count, org, ss, cs, _, alg := base()
		_, err := GenerateMTLSCertificatesMultiClient(ca, server, prefix, count, org, ss, cs, 0, alg)
		if err == nil {
			t.Error("expected error")
		}
	})
}

// ---- AddCertificatesToCRL ----

func TestAddCertificatesToCRL(t *testing.T) {
	// テスト用の mTLS セット（EdDSA）を生成
	mtls, err := GenerateMTLSCertificatesMultiClient(
		"Test CA", "server.example.com", "client",
		2, "Test Org", nil, nil, 30, "EdDSA",
	)
	if err != nil {
		t.Fatalf("GenerateMTLSCertificatesMultiClient failed: %v", err)
	}

	t.Run("クライアント証明書を1件失効", func(t *testing.T) {
		result, err := AddCertificatesToCRL(
			mtls.CACertificatePEM,
			mtls.CAPrivateKeyPEM,
			mtls.CRLPEM,
			[]string{mtls.ClientCertificates[0].CertificatePEM},
			7,
		)
		if err != nil {
			t.Fatalf("AddCertificatesToCRL() error = %v", err)
		}
		if result.AddedCount != 1 {
			t.Errorf("AddedCount = %d, want 1", result.AddedCount)
		}
		if result.TotalRevokedCount != 1 {
			t.Errorf("TotalRevokedCount = %d, want 1", result.TotalRevokedCount)
		}
		if result.CRLPEM == "" {
			t.Error("CRLPEM should not be empty")
		}
	})

	t.Run("同じ証明書を2回追加しても重複しない", func(t *testing.T) {
		first, err := AddCertificatesToCRL(
			mtls.CACertificatePEM, mtls.CAPrivateKeyPEM, mtls.CRLPEM,
			[]string{mtls.ClientCertificates[0].CertificatePEM}, 7,
		)
		if err != nil {
			t.Fatalf("1st AddCertificatesToCRL failed: %v", err)
		}
		second, err := AddCertificatesToCRL(
			mtls.CACertificatePEM, mtls.CAPrivateKeyPEM, first.CRLPEM,
			[]string{mtls.ClientCertificates[0].CertificatePEM}, 7,
		)
		if err != nil {
			t.Fatalf("2nd AddCertificatesToCRL failed: %v", err)
		}
		if second.TotalRevokedCount != 1 {
			t.Errorf("重複後の TotalRevokedCount = %d, want 1", second.TotalRevokedCount)
		}
		if second.AddedCount != 0 {
			t.Errorf("重複の AddedCount = %d, want 0", second.AddedCount)
		}
	})

	t.Run("2件失効", func(t *testing.T) {
		result, err := AddCertificatesToCRL(
			mtls.CACertificatePEM, mtls.CAPrivateKeyPEM, mtls.CRLPEM,
			[]string{
				mtls.ClientCertificates[0].CertificatePEM,
				mtls.ClientCertificates[1].CertificatePEM,
			}, 7,
		)
		if err != nil {
			t.Fatalf("AddCertificatesToCRL() error = %v", err)
		}
		if result.TotalRevokedCount != 2 {
			t.Errorf("TotalRevokedCount = %d, want 2", result.TotalRevokedCount)
		}
		if len(result.RevokedSerialNumbers) != 2 {
			t.Errorf("RevokedSerialNumbers count = %d, want 2", len(result.RevokedSerialNumbers))
		}
	})

	t.Run("caCertificate 空はエラー", func(t *testing.T) {
		_, err := AddCertificatesToCRL("", mtls.CAPrivateKeyPEM, mtls.CRLPEM,
			[]string{mtls.ClientCertificates[0].CertificatePEM}, 7)
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("caPrivateKey 空はエラー", func(t *testing.T) {
		_, err := AddCertificatesToCRL(mtls.CACertificatePEM, "", mtls.CRLPEM,
			[]string{mtls.ClientCertificates[0].CertificatePEM}, 7)
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("revokedCerts 空はエラー", func(t *testing.T) {
		_, err := AddCertificatesToCRL(mtls.CACertificatePEM, mtls.CAPrivateKeyPEM, mtls.CRLPEM,
			[]string{}, 7)
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("既存CRL なし（空文字）でも動作する", func(t *testing.T) {
		result, err := AddCertificatesToCRL(
			mtls.CACertificatePEM, mtls.CAPrivateKeyPEM, "",
			[]string{mtls.ClientCertificates[0].CertificatePEM}, 7,
		)
		if err != nil {
			t.Fatalf("AddCertificatesToCRL() error = %v", err)
		}
		if result.AddedCount != 1 {
			t.Errorf("AddedCount = %d, want 1", result.AddedCount)
		}
	})
}

// ---- parseCertificateFromPEM (内部関数) ----

func TestParseCertificateFromPEM(t *testing.T) {
	result, err := GenerateSelfSignedCertificate("example.com", "Test", nil, 30, "EdDSA")
	if err != nil {
		t.Fatalf("GenerateSelfSignedCertificate failed: %v", err)
	}

	t.Run("正常な証明書", func(t *testing.T) {
		cert, err := parseCertificateFromPEM(result.CertificatePEM)
		if err != nil {
			t.Errorf("parseCertificateFromPEM() error = %v", err)
			return
		}
		if cert.Subject.CommonName != "example.com" {
			t.Errorf("CommonName = %q, want %q", cert.Subject.CommonName, "example.com")
		}
	})

	t.Run("不正なPEM", func(t *testing.T) {
		_, err := parseCertificateFromPEM("not a pem")
		if err == nil {
			t.Error("expected error")
		}
	})

	t.Run("PRIVATE KEY PEMを渡すとエラー", func(t *testing.T) {
		_, err := parseCertificateFromPEM(result.PrivateKeyPEM)
		if err == nil {
			t.Error("expected error for non-certificate PEM")
		}
	})
}

// ---- GenerateSelfSignedCertificate 追加ケース ----

func TestGenerateSelfSignedCertificateRS256(t *testing.T) {
	got, err := GenerateSelfSignedCertificate("test.example.com", "Test Org", nil, 365, "RS256")
	if err != nil {
		t.Fatalf("GenerateSelfSignedCertificate(RS256) error = %v", err)
	}
	if !strings.Contains(got.CertificatePEM, "-----BEGIN CERTIFICATE-----") {
		t.Error("CertificatePEM does not contain certificate header")
	}
}

func TestGenerateSelfSignedCertificateDefaultAlgorithm(t *testing.T) {
	// algorithm 空文字 → RS256 にフォールバック
	got, err := GenerateSelfSignedCertificate("test.example.com", "", nil, 365, "")
	if err != nil {
		t.Fatalf("GenerateSelfSignedCertificate('') error = %v", err)
	}
	if !strings.Contains(got.CertificatePEM, "-----BEGIN CERTIFICATE-----") {
		t.Error("CertificatePEM does not contain certificate header")
	}
}

func TestGenerateSelfSignedCertificateIPOnly(t *testing.T) {
	// SAN が IP アドレスのみの場合 commonName が DNS に追加されない
	got, err := GenerateSelfSignedCertificate("example.com", "Org", []string{"192.168.1.1"}, 30, "EdDSA")
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	block, _ := pem.Decode([]byte(got.CertificatePEM))
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatalf("parse cert error = %v", err)
	}
	if len(cert.IPAddresses) == 0 {
		t.Error("expected at least 1 IP address in SAN")
	}
}

// ---- シリアル番号がソートされて返ることを確認 ----

func TestAddCertificatesToCRLSerialNumbersSorted(t *testing.T) {
	mtls, err := GenerateMTLSCertificatesMultiClient(
		"Test CA", "server.example.com", "client",
		3, "Test Org", nil, nil, 30, "EdDSA",
	)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	result, err := AddCertificatesToCRL(
		mtls.CACertificatePEM, mtls.CAPrivateKeyPEM, "",
		[]string{
			mtls.ClientCertificates[0].CertificatePEM,
			mtls.ClientCertificates[1].CertificatePEM,
			mtls.ClientCertificates[2].CertificatePEM,
		}, 7,
	)
	if err != nil {
		t.Fatalf("AddCertificatesToCRL() error = %v", err)
	}

	for i := 1; i < len(result.RevokedSerialNumbers); i++ {
		prev := result.RevokedSerialNumbers[i-1]
		curr := result.RevokedSerialNumbers[i]
		if prev > curr {
			t.Errorf("RevokedSerialNumbers not sorted at index %d: %s > %s",
				i, prev, curr)
		}
	}
}
