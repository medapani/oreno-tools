package backend

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// JWTDecodeResult はJWTデコード結果を保持する構造体
type JWTDecodeResult struct {
	Header  string `json:"header"`
	Payload string `json:"payload"`
	Valid   bool   `json:"valid"`
	Error   string `json:"error"`
}

// JWTEncode はペイロード（JSON文字列）、秘密鍵、アルゴリズムからJWTトークンを生成する
func JWTEncode(payload string, secret string, algorithm string) (string, error) {
	if secret == "" {
		return "", fmt.Errorf("secret key cannot be empty")
	}

	// ペイロードをJSONとしてパース
	var claims jwt.MapClaims
	if err := json.Unmarshal([]byte(payload), &claims); err != nil {
		return "", fmt.Errorf("invalid JSON payload: %w", err)
	}

	// 署名メソッドを取得
	signingMethod, err := getSigningMethod(algorithm)
	if err != nil {
		return "", err
	}

	// JWTトークンを生成
	token := jwt.NewWithClaims(signingMethod, claims)

	var tokenString string
	// 非対称鍵暗号方式の場合はPEM形式の秘密鍵をパース
	if isAsymmetricAlgorithm(algorithm) {
		privateKey, err := parsePrivateKey(secret)
		if err != nil {
			return "", fmt.Errorf("failed to parse private key: %w", err)
		}
		tokenString, err = token.SignedString(privateKey)
		if err != nil {
			return "", fmt.Errorf("failed to sign JWT with asymmetric key: %w", err)
		}
	} else {
		// HMAC方式の場合は秘密鍵をバイト列として使用
		tokenString, err = token.SignedString([]byte(secret))
	}

	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return tokenString, nil
}

// JWTDecode はJWTトークンをパースしてヘッダーとペイロードを返す（署名検証なし）
func JWTDecode(token string) (JWTDecodeResult, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return JWTDecodeResult{}, fmt.Errorf("token cannot be empty")
	}

	// トークンを3つの部分に分割
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return JWTDecodeResult{}, fmt.Errorf("invalid JWT format: expected 3 parts, got %d", len(parts))
	}

	// ヘッダーをデコード
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return JWTDecodeResult{}, fmt.Errorf("failed to decode header: %w", err)
	}

	// ペイロードをデコード
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return JWTDecodeResult{}, fmt.Errorf("failed to decode payload: %w", err)
	}

	// JSONを整形
	var headerJSON, payloadJSON interface{}
	if err := json.Unmarshal(headerBytes, &headerJSON); err != nil {
		return JWTDecodeResult{}, fmt.Errorf("invalid header JSON: %w", err)
	}
	if err := json.Unmarshal(payloadBytes, &payloadJSON); err != nil {
		return JWTDecodeResult{}, fmt.Errorf("invalid payload JSON: %w", err)
	}

	// 整形されたJSONを文字列に変換
	headerFormatted, _ := json.MarshalIndent(headerJSON, "", "  ")
	payloadFormatted, _ := json.MarshalIndent(payloadJSON, "", "  ")

	return JWTDecodeResult{
		Header:  string(headerFormatted),
		Payload: string(payloadFormatted),
		Valid:   false,
		Error:   "",
	}, nil
}

// JWTVerify はJWTトークンの署名を検証し、ヘッダーとペイロードを返す
func JWTVerify(tokenString string, secret string) (JWTDecodeResult, error) {
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return JWTDecodeResult{}, fmt.Errorf("token cannot be empty")
	}

	// トークンをパースして検証
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// アルゴリズムに応じて適切な鍵を返す
		switch token.Method.(type) {
		case *jwt.SigningMethodHMAC:
			return []byte(secret), nil
		case *jwt.SigningMethodRSA, *jwt.SigningMethodRSAPSS:
			// RSA公開鍵をパース
			publicKey, err := parsePublicKey(secret)
			if err != nil {
				return nil, fmt.Errorf("failed to parse RSA public key: %w", err)
			}
			if _, ok := publicKey.(*rsa.PublicKey); !ok {
				return nil, fmt.Errorf("expected RSA public key, got %T", publicKey)
			}
			return publicKey, nil
		case *jwt.SigningMethodECDSA:
			// ECDSA公開鍵をパース
			publicKey, err := parsePublicKey(secret)
			if err != nil {
				return nil, fmt.Errorf("failed to parse ECDSA public key: %w", err)
			}
			if _, ok := publicKey.(*ecdsa.PublicKey); !ok {
				return nil, fmt.Errorf("expected ECDSA public key, got %T", publicKey)
			}
			return publicKey, nil
		case *jwt.SigningMethodEd25519:
			// EdDSA公開鍵をパース
			publicKey, err := parsePublicKey(secret)
			if err != nil {
				return nil, fmt.Errorf("failed to parse EdDSA public key: %w", err)
			}
			if _, ok := publicKey.(ed25519.PublicKey); !ok {
				return nil, fmt.Errorf("expected EdDSA public key, got %T", publicKey)
			}
			return publicKey, nil
		default:
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
	})

	// ヘッダーとペイロードを整形
	var headerJSON, payloadJSON interface{}
	parts := strings.Split(tokenString, ".")
	if len(parts) >= 2 {
		headerBytes, _ := base64.RawURLEncoding.DecodeString(parts[0])
		payloadBytes, _ := base64.RawURLEncoding.DecodeString(parts[1])
		err := json.Unmarshal(headerBytes, &headerJSON)
		if err != nil {
			return JWTDecodeResult{}, fmt.Errorf("invalid header JSON: %w", err)
		}
		err = json.Unmarshal(payloadBytes, &payloadJSON)
		if err != nil {
			return JWTDecodeResult{}, fmt.Errorf("invalid payload JSON: %w", err)
		}
	}

	headerFormatted, _ := json.MarshalIndent(headerJSON, "", "  ")
	payloadFormatted, _ := json.MarshalIndent(payloadJSON, "", "  ")

	result := JWTDecodeResult{
		Header:  string(headerFormatted),
		Payload: string(payloadFormatted),
	}

	if err != nil {
		result.Valid = false
		result.Error = err.Error()
		return result, nil
	}

	if !token.Valid {
		result.Valid = false
		result.Error = "invalid token"
		return result, nil
	}

	result.Valid = true
	result.Error = ""
	return result, nil
}

// ExtractPublicKey は秘密鍵から公開鍵を抽出する
// 非対称鍵アルゴリズム（RSA、ECDSA、EdDSA）にのみ対応
func ExtractPublicKey(privateKey string, algorithm string) (string, error) {
	privateKey = strings.TrimSpace(privateKey)
	if privateKey == "" {
		return "", fmt.Errorf("private key cannot be empty")
	}

	// サポートされているアルゴリズムかチェック
	if !isAsymmetricAlgorithm(algorithm) {
		return "", fmt.Errorf("algorithm %s does not support public key extraction (only RSA, ECDSA, and EdDSA are supported)", algorithm)
	}

	// 秘密鍵をパース
	parsedKey, err := parsePrivateKey(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}

	var publicKey interface{}
	var publicKeyDER []byte

	// アルゴリズムに応じて公開鍵を抽出
	switch algorithm[:2] {
	case "RS":
		// RSA秘密鍵から公開鍵を抽出
		rsaPrivateKey, ok := parsedKey.(*rsa.PrivateKey)
		if !ok {
			return "", fmt.Errorf("expected RSA private key for algorithm %s", algorithm)
		}
		publicKey = &rsaPrivateKey.PublicKey
		var err error
		publicKeyDER, err = x509.MarshalPKIXPublicKey(publicKey)
		if err != nil {
			return "", fmt.Errorf("failed to marshal RSA public key: %w", err)
		}

	case "ES":
		// ECDSA秘密鍵から公開鍵を抽出
		ecdsaPrivateKey, ok := parsedKey.(*ecdsa.PrivateKey)
		if !ok {
			return "", fmt.Errorf("expected ECDSA private key for algorithm %s", algorithm)
		}
		publicKey = &ecdsaPrivateKey.PublicKey
		var err error
		publicKeyDER, err = x509.MarshalPKIXPublicKey(publicKey)
		if err != nil {
			return "", fmt.Errorf("failed to marshal ECDSA public key: %w", err)
		}

	case "PS":
		// RSA-PSS秘密鍵から公開鍵を抽出（RSAと同じ）
		rsaPrivateKey, ok := parsedKey.(*rsa.PrivateKey)
		if !ok {
			return "", fmt.Errorf("expected RSA private key for algorithm %s", algorithm)
		}
		publicKey = &rsaPrivateKey.PublicKey
		var err error
		publicKeyDER, err = x509.MarshalPKIXPublicKey(publicKey)
		if err != nil {
			return "", fmt.Errorf("failed to marshal RSA public key: %w", err)
		}

	case "Ed":
		// EdDSA（Ed25519）秘密鍵から公開鍵を抽出
		ed25519PrivateKey, ok := parsedKey.(ed25519.PrivateKey)
		if !ok {
			return "", fmt.Errorf("expected Ed25519 private key for algorithm %s", algorithm)
		}
		publicKey = ed25519PrivateKey.Public()
		var err error
		publicKeyDER, err = x509.MarshalPKIXPublicKey(publicKey)
		if err != nil {
			return "", fmt.Errorf("failed to marshal Ed25519 public key: %w", err)
		}

	default:
		return "", fmt.Errorf("unexpected algorithm: %s", algorithm)
	}

	// PEM形式にエンコード
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyDER,
	})

	return string(publicKeyPEM), nil
}

// getSigningMethod はアルゴリズム文字列から署名メソッドを返す
func getSigningMethod(algorithm string) (jwt.SigningMethod, error) {
	switch algorithm {
	case "HS256":
		return jwt.SigningMethodHS256, nil
	case "HS384":
		return jwt.SigningMethodHS384, nil
	case "HS512":
		return jwt.SigningMethodHS512, nil
	case "RS256":
		return jwt.SigningMethodRS256, nil
	case "RS384":
		return jwt.SigningMethodRS384, nil
	case "RS512":
		return jwt.SigningMethodRS512, nil
	case "ES256":
		return jwt.SigningMethodES256, nil
	case "ES384":
		return jwt.SigningMethodES384, nil
	case "ES512":
		return jwt.SigningMethodES512, nil
	case "PS256":
		return jwt.SigningMethodPS256, nil
	case "PS384":
		return jwt.SigningMethodPS384, nil
	case "PS512":
		return jwt.SigningMethodPS512, nil
	case "EdDSA":
		return jwt.SigningMethodEdDSA, nil
	default:
		return nil, fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
}
