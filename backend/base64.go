package backend

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// Base64Encode は文字列を Base64 へエンコードする
// urlSafe が true の場合、URL-safe Base64 (RFC 4648) を使用
func Base64Encode(input string, urlSafe bool) string {
	if urlSafe {
		return base64.URLEncoding.EncodeToString([]byte(input))
	}
	return base64.StdEncoding.EncodeToString([]byte(input))
}

// Base64Decode は Base64 文字列をデコードする
// urlSafe が true の場合、URL-safe Base64 (RFC 4648) を使用
func Base64Decode(input string, urlSafe bool) (string, error) {
	inputTrimmed := strings.TrimSpace(input)
	var decoded []byte
	var err error

	if urlSafe {
		decoded, err = base64.URLEncoding.DecodeString(inputTrimmed)
		if err != nil {
			// パディングなしの場合も試す
			decoded, err = base64.RawURLEncoding.DecodeString(inputTrimmed)
		}
	} else {
		decoded, err = base64.StdEncoding.DecodeString(inputTrimmed)
	}

	if err != nil {
		return "", fmt.Errorf("invalid base64 input: %w", err)
	}

	return string(decoded), nil
}
