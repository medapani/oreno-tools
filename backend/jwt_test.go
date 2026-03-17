package backend

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"strings"
	"testing"
)

func TestJWTEncode(t *testing.T) {
	tests := []struct {
		name      string
		payload   string
		secret    string
		algorithm string
		wantErr   bool
	}{
		{
			name:      "HS256 正常",
			payload:   `{"sub":"1234567890","name":"John Doe","iat":1516239022}`,
			secret:    "mysecret",
			algorithm: "HS256",
			wantErr:   false,
		},
		{
			name:      "HS384 正常",
			payload:   `{"sub":"test"}`,
			secret:    "mysecret",
			algorithm: "HS384",
			wantErr:   false,
		},
		{
			name:      "HS512 正常",
			payload:   `{"sub":"test"}`,
			secret:    "mysecret",
			algorithm: "HS512",
			wantErr:   false,
		},
		{
			name:      "空のsecret",
			payload:   `{"sub":"test"}`,
			secret:    "",
			algorithm: "HS256",
			wantErr:   true,
		},
		{
			name:      "不正なペイロード JSON",
			payload:   `{not valid json}`,
			secret:    "mysecret",
			algorithm: "HS256",
			wantErr:   true,
		},
		{
			name:      "不正なアルゴリズム",
			payload:   `{"sub":"test"}`,
			secret:    "mysecret",
			algorithm: "INVALID",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JWTEncode(tt.payload, tt.secret, tt.algorithm)
			if (err != nil) != tt.wantErr {
				t.Errorf("JWTEncode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// JWTは3つのパートに分かれているはず
				parts := strings.Split(got, ".")
				if len(parts) != 3 {
					t.Errorf("JWTEncode() returned invalid JWT format: %q", got)
				}
			}
		})
	}
}

func TestJWTDecode(t *testing.T) {
	// テスト用トークンを先に生成
	payload := `{"sub":"1234567890","name":"John Doe","iat":1516239022}`
	token, err := JWTEncode(payload, "mysecret", "HS256")
	if err != nil {
		t.Fatalf("failed to create test token: %v", err)
	}

	tests := []struct {
		name        string
		token       string
		wantSubject string
		wantErr     bool
	}{
		{
			name:        "正常なトークン",
			token:       token,
			wantSubject: "1234567890",
			wantErr:     false,
		},
		{
			name:    "空のトークン",
			token:   "",
			wantErr: true,
		},
		{
			name:    "不正なフォーマット (2パーツ)",
			token:   "part1.part2",
			wantErr: true,
		},
		{
			name:    "不正なbase64ヘッダー",
			token:   "!!!.payload.sig",
			wantErr: true,
		},
		{
			name:    "前後の空白を含むトークン",
			token:   "  " + token + "  ",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JWTDecode(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("JWTDecode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			// ペイロードをJSONとしてパース
			var claims map[string]interface{}
			if err := json.Unmarshal([]byte(got.Payload), &claims); err != nil {
				t.Errorf("JWTDecode() Payload is not valid JSON: %v", err)
				return
			}

			if tt.wantSubject != "" {
				sub, ok := claims["sub"].(string)
				if !ok || sub != tt.wantSubject {
					t.Errorf("Payload sub = %q, want %q", sub, tt.wantSubject)
				}
			}
		})
	}
}

func TestJWTVerify(t *testing.T) {
	secret := "mysecret"
	payload := `{"sub":"1234567890","name":"Test User"}`

	validToken, err := JWTEncode(payload, secret, "HS256")
	if err != nil {
		t.Fatalf("failed to create test token: %v", err)
	}

	tests := []struct {
		name      string
		token     string
		secret    string
		wantValid bool
		wantErr   bool
	}{
		{
			name:      "正しいsecretで検証成功",
			token:     validToken,
			secret:    secret,
			wantValid: true,
			wantErr:   false,
		},
		{
			name:      "誤ったsecretで検証失敗",
			token:     validToken,
			secret:    "wrongsecret",
			wantValid: false,
			wantErr:   false,
		},
		{
			name:    "空のトークン",
			token:   "",
			secret:  secret,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JWTVerify(tt.token, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("JWTVerify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.Valid != tt.wantValid {
				t.Errorf("JWTVerify().Valid = %v, want %v (Error: %s)", got.Valid, tt.wantValid, got.Error)
			}
		})
	}
}

func TestJWTEncodeDecodeRoundTrip(t *testing.T) {
	original := map[string]interface{}{
		"sub":  "user123",
		"name": "テストユーザー",
		"role": "admin",
	}

	payloadBytes, _ := json.Marshal(original)
	// #nosec G101 -- テスト用の固定シークレット
	secret := "round-trip-secret"

	token, err := JWTEncode(string(payloadBytes), secret, "HS256")
	if err != nil {
		t.Fatalf("JWTEncode() error: %v", err)
	}

	decoded, err := JWTDecode(token)
	if err != nil {
		t.Fatalf("JWTDecode() error: %v", err)
	}

	var claims map[string]interface{}
	if err := json.Unmarshal([]byte(decoded.Payload), &claims); err != nil {
		t.Fatalf("Payload JSON parse error: %v", err)
	}

	if claims["sub"] != original["sub"] {
		t.Errorf("sub = %v, want %v", claims["sub"], original["sub"])
	}
	if claims["name"] != original["name"] {
		t.Errorf("name = %v, want %v", claims["name"], original["name"])
	}
}

func TestIsAsymmetricAlgorithm(t *testing.T) {
	tests := []struct {
		algorithm string
		want      bool
	}{
		{"RS256", true},
		{"RS384", true},
		{"RS512", true},
		{"ES256", true},
		{"ES384", true},
		{"ES512", true},
		{"PS256", true},
		{"PS384", true},
		{"PS512", true},
		{"EdDSA", true},
		{"HS256", false},
		{"HS384", false},
		{"HS512", false},
		{"INVALID", false},
	}

	for _, tt := range tests {
		t.Run(tt.algorithm, func(t *testing.T) {
			got := isAsymmetricAlgorithm(tt.algorithm)
			if got != tt.want {
				t.Errorf("isAsymmetricAlgorithm(%q) = %v, want %v", tt.algorithm, got, tt.want)
			}
		})
	}
}

// ---- getSigningMethod ----

func TestGetSigningMethod(t *testing.T) {
	algorithms := []string{
		"HS256", "HS384", "HS512",
		"RS256", "RS384", "RS512",
		"ES256", "ES384", "ES512",
		"PS256", "PS384", "PS512",
		"EdDSA",
	}
	for _, alg := range algorithms {
		t.Run(alg, func(t *testing.T) {
			got, err := getSigningMethod(alg)
			if err != nil {
				t.Errorf("getSigningMethod(%q) error = %v", alg, err)
				return
			}
			if got == nil {
				t.Errorf("getSigningMethod(%q) returned nil", alg)
			}
		})
	}

	t.Run("不正なアルゴリズム", func(t *testing.T) {
		_, err := getSigningMethod("INVALID")
		if err == nil {
			t.Error("expected error for invalid algorithm, got nil")
		}
	})
}

// ---- parsePrivateKey (security.go) ----

func TestParsePrivateKey(t *testing.T) {
	t.Run("PKCS8形式 (EdDSA)", func(t *testing.T) {
		keyPEM, err := GeneratePrivateKey("EdDSA")
		if err != nil {
			t.Fatalf("GeneratePrivateKey failed: %v", err)
		}
		got, err := parsePrivateKey(keyPEM)
		if err != nil {
			t.Errorf("parsePrivateKey() error = %v", err)
			return
		}
		if got == nil {
			t.Error("parsePrivateKey() returned nil")
		}
	})

	t.Run("PKCS1形式 (RSA)", func(t *testing.T) {
		rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatalf("rsa.GenerateKey failed: %v", err)
		}
		pkcs1PEM := string(pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(rsaKey),
		}))
		got, err := parsePrivateKey(pkcs1PEM)
		if err != nil {
			t.Errorf("parsePrivateKey() PKCS1 error = %v", err)
			return
		}
		if got == nil {
			t.Error("parsePrivateKey() PKCS1 returned nil")
		}
	})

	t.Run("EC形式 (ECDSA)", func(t *testing.T) {
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
		got, err := parsePrivateKey(ecPEM)
		if err != nil {
			t.Errorf("parsePrivateKey() EC error = %v", err)
			return
		}
		if got == nil {
			t.Error("parsePrivateKey() EC returned nil")
		}
	})

	t.Run("無効なPEMブロック", func(t *testing.T) {
		_, err := parsePrivateKey("not a pem")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("サポート外フォーマット", func(t *testing.T) {
		unsupportedPEM := string(pem.EncodeToMemory(&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: []byte("invalid der bytes"),
		}))
		_, err := parsePrivateKey(unsupportedPEM)
		if err == nil {
			t.Error("expected error for unsupported format, got nil")
		}
	})
}

// ---- parsePublicKey ----

func TestParsePublicKey(t *testing.T) {
	t.Run("PKIX形式 (EdDSA公開鍵)", func(t *testing.T) {
		keyPEM, err := GeneratePrivateKey("EdDSA")
		if err != nil {
			t.Fatalf("GeneratePrivateKey failed: %v", err)
		}
		pubPEM, err := ExtractPublicKey(keyPEM, "EdDSA")
		if err != nil {
			t.Fatalf("ExtractPublicKey failed: %v", err)
		}
		got, err := parsePublicKey(pubPEM)
		if err != nil {
			t.Errorf("parsePublicKey() error = %v", err)
			return
		}
		if got == nil {
			t.Error("parsePublicKey() returned nil")
		}
	})

	t.Run("PKCS1形式 (RSA公開鍵)", func(t *testing.T) {
		rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatalf("rsa.GenerateKey failed: %v", err)
		}
		pkcs1PEM := string(pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(&rsaKey.PublicKey),
		}))
		got, err := parsePublicKey(pkcs1PEM)
		if err != nil {
			t.Errorf("parsePublicKey() PKCS1 error = %v", err)
			return
		}
		if got == nil {
			t.Error("parsePublicKey() PKCS1 returned nil")
		}
	})

	t.Run("無効なPEMブロック", func(t *testing.T) {
		_, err := parsePublicKey("not a pem")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("サポート外フォーマット", func(t *testing.T) {
		unsupportedPEM := string(pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: []byte("invalid der bytes"),
		}))
		_, err := parsePublicKey(unsupportedPEM)
		if err == nil {
			t.Error("expected error for unsupported format, got nil")
		}
	})
}

// ---- generatePrivateKeyForAlgorithm ----

func TestGeneratePrivateKeyForAlgorithm(t *testing.T) {
	tests := []struct {
		name    string
		alg     string
		wantErr bool
	}{
		{"RS256 (RSA)", "RS256", false},
		{"ES256 (P-256)", "ES256", false},
		{"ES384 (P-384)", "ES384", false},
		{"ES512 (P-521)", "ES512", false},
		{"EdDSA", "EdDSA", false},
		{"不正なESアルゴリズム (ES128)", "ES128", true},
		{"完全に不正なアルゴリズム", "INVALID", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generatePrivateKeyForAlgorithm(tt.alg)
			if (err != nil) != tt.wantErr {
				t.Errorf("generatePrivateKeyForAlgorithm(%q) error = %v, wantErr %v", tt.alg, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("generatePrivateKeyForAlgorithm() returned nil without error")
			}
		})
	}
}

// ---- getPublicKeyFromPrivate ----

func TestGetPublicKeyFromPrivate(t *testing.T) {
	t.Run("RSA秘密鍵", func(t *testing.T) {
		rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
		if err != nil {
			t.Fatalf("rsa.GenerateKey failed: %v", err)
		}
		got := getPublicKeyFromPrivate(rsaKey)
		if got == nil {
			t.Error("getPublicKeyFromPrivate() RSA returned nil")
		}
	})

	t.Run("ECDSA秘密鍵", func(t *testing.T) {
		ecKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			t.Fatalf("ecdsa.GenerateKey failed: %v", err)
		}
		got := getPublicKeyFromPrivate(ecKey)
		if got == nil {
			t.Error("getPublicKeyFromPrivate() ECDSA returned nil")
		}
	})

	t.Run("Ed25519秘密鍵", func(t *testing.T) {
		_, edKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			t.Fatalf("ed25519.GenerateKey failed: %v", err)
		}
		got := getPublicKeyFromPrivate(edKey)
		if got == nil {
			t.Error("getPublicKeyFromPrivate() Ed25519 returned nil")
		}
	})

	t.Run("サポート外の型 -> nil", func(t *testing.T) {
		got := getPublicKeyFromPrivate("string is not a key")
		if got != nil {
			t.Errorf("getPublicKeyFromPrivate() unsupported type = %v, want nil", got)
		}
	})
}

// ---- getSignerFromPrivateKey ----

func TestGetSignerFromPrivateKey(t *testing.T) {
	t.Run("Ed25519秘密鍵 (Signerを実装)", func(t *testing.T) {
		_, edKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			t.Fatalf("ed25519.GenerateKey failed: %v", err)
		}
		got, err := getSignerFromPrivateKey(edKey)
		if err != nil {
			t.Errorf("getSignerFromPrivateKey() error = %v", err)
			return
		}
		if got == nil {
			t.Error("getSignerFromPrivateKey() returned nil")
		}
	})

	t.Run("Signerを実装しない型 -> エラー", func(t *testing.T) {
		_, err := getSignerFromPrivateKey("not a signer")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

// ---- JWTEncode (非対称鍵) ----

func TestJWTEncodeAsymmetric(t *testing.T) {
	payload := `{"sub":"test","name":"Alice"}`

	t.Run("EdDSAで署名", func(t *testing.T) {
		keyPEM, err := GeneratePrivateKey("EdDSA")
		if err != nil {
			t.Fatalf("GeneratePrivateKey failed: %v", err)
		}
		got, err := JWTEncode(payload, keyPEM, "EdDSA")
		if err != nil {
			t.Errorf("JWTEncode() EdDSA error = %v", err)
			return
		}
		if len(strings.Split(got, ".")) != 3 {
			t.Errorf("JWTEncode() invalid JWT format: %q", got)
		}
	})

	t.Run("ES256で署名", func(t *testing.T) {
		keyPEM, err := GeneratePrivateKey("ES256")
		if err != nil {
			t.Fatalf("GeneratePrivateKey failed: %v", err)
		}
		got, err := JWTEncode(payload, keyPEM, "ES256")
		if err != nil {
			t.Errorf("JWTEncode() ES256 error = %v", err)
			return
		}
		if len(strings.Split(got, ".")) != 3 {
			t.Errorf("JWTEncode() invalid JWT format: %q", got)
		}
	})

	t.Run("不正な秘密鍵PEMでエラー", func(t *testing.T) {
		_, err := JWTEncode(payload, "not a valid pem", "RS256")
		if err == nil {
			t.Error("expected error for invalid private key, got nil")
		}
	})
}

// ---- JWTVerify (非対称鍵) ----

func TestJWTVerifyAsymmetric(t *testing.T) {
	payload := `{"sub":"test"}`

	t.Run("EdDSA 署名検証成功", func(t *testing.T) {
		keyPEM, err := GeneratePrivateKey("EdDSA")
		if err != nil {
			t.Fatalf("GeneratePrivateKey failed: %v", err)
		}
		pubPEM, err := ExtractPublicKey(keyPEM, "EdDSA")
		if err != nil {
			t.Fatalf("ExtractPublicKey failed: %v", err)
		}
		token, err := JWTEncode(payload, keyPEM, "EdDSA")
		if err != nil {
			t.Fatalf("JWTEncode failed: %v", err)
		}
		got, err := JWTVerify(token, pubPEM)
		if err != nil {
			t.Errorf("JWTVerify() error = %v", err)
			return
		}
		if !got.Valid {
			t.Errorf("JWTVerify() Valid = false, want true (Error: %s)", got.Error)
		}
	})

	t.Run("ES256 署名検証成功", func(t *testing.T) {
		keyPEM, err := GeneratePrivateKey("ES256")
		if err != nil {
			t.Fatalf("GeneratePrivateKey failed: %v", err)
		}
		pubPEM, err := ExtractPublicKey(keyPEM, "ES256")
		if err != nil {
			t.Fatalf("ExtractPublicKey failed: %v", err)
		}
		token, err := JWTEncode(payload, keyPEM, "ES256")
		if err != nil {
			t.Fatalf("JWTEncode failed: %v", err)
		}
		got, err := JWTVerify(token, pubPEM)
		if err != nil {
			t.Errorf("JWTVerify() error = %v", err)
			return
		}
		if !got.Valid {
			t.Errorf("JWTVerify() Valid = false, want true (Error: %s)", got.Error)
		}
	})

	t.Run("不正な公開鍵PEMで検証失敗 (Valid=false)", func(t *testing.T) {
		keyPEM, err := GeneratePrivateKey("EdDSA")
		if err != nil {
			t.Fatalf("GeneratePrivateKey failed: %v", err)
		}
		token, err := JWTEncode(payload, keyPEM, "EdDSA")
		if err != nil {
			t.Fatalf("JWTEncode failed: %v", err)
		}
		got, err := JWTVerify(token, "invalid pem")
		if err != nil {
			t.Errorf("JWTVerify() unexpected error = %v", err)
			return
		}
		if got.Valid {
			t.Error("JWTVerify() Valid = true, want false for invalid public key")
		}
	})
}

// ---- ExtractPublicKey (追加パス) ----

func TestExtractPublicKeyAdditional(t *testing.T) {
	t.Run("PS256 (RSA-PSS) から公開鍵抽出", func(t *testing.T) {
		keyPEM, err := GeneratePrivateKey("PS256")
		if err != nil {
			t.Fatalf("GeneratePrivateKey failed: %v", err)
		}
		pubPEM, err := ExtractPublicKey(keyPEM, "PS256")
		if err != nil {
			t.Errorf("ExtractPublicKey() PS256 error = %v", err)
			return
		}
		if pubPEM == "" {
			t.Error("ExtractPublicKey() PS256 returned empty string")
		}
	})
}
