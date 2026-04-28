package backend

import (
	"encoding/base64"
	"fmt"
	"strings"

	qrcode "github.com/skip2/go-qrcode"
)

// GenerateQRCode はテキストから QR コードを生成し、base64 エンコードされた PNG 文字列を返す。
// recoveryLevel は "L"/"M"/"Q"/"H" のいずれかを指定する。
// size は QR コード画像のピクセルサイズ（64〜2048）。
func GenerateQRCode(text string, size int, recoveryLevel string) (string, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return "", fmt.Errorf("テキストを入力してください")
	}
	if size < 64 || size > 2048 {
		return "", fmt.Errorf("サイズは 64〜2048 px の範囲で指定してください")
	}

	var level qrcode.RecoveryLevel
	switch strings.ToUpper(recoveryLevel) {
	case "L":
		level = qrcode.Low
	case "M":
		level = qrcode.Medium
	case "Q":
		level = qrcode.High
	case "H":
		level = qrcode.Highest
	default:
		return "", fmt.Errorf("不正な誤り訂正レベル: %s（L/M/Q/H のいずれかを指定してください）", recoveryLevel)
	}

	png, err := qrcode.Encode(text, level, size)
	if err != nil {
		return "", fmt.Errorf("QR コードの生成に失敗しました: %w", err)
	}

	return base64.StdEncoding.EncodeToString(png), nil
}
