package backend

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
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

func decodeQRCodePNG(base64PNG string) ([]byte, error) {
	trimmed := strings.TrimSpace(base64PNG)
	if trimmed == "" {
		return nil, fmt.Errorf("保存するQR画像がありません")
	}

	pngBytes, err := base64.StdEncoding.DecodeString(trimmed)
	if err != nil {
		return nil, fmt.Errorf("QR画像データのデコードに失敗しました: %w", err)
	}

	return pngBytes, nil
}

func ensurePNGFilename(filename string) string {
	safeFilename := sanitizeFilename(filename)
	if !strings.HasSuffix(strings.ToLower(safeFilename), ".png") {
		safeFilename += ".png"
	}
	return safeFilename
}

// SaveQRCodePNG は base64 文字列の PNG を Downloads/oreno-tools-qr 配下へ保存する。
func SaveQRCodePNG(base64PNG string, filename string) (string, error) {
	pngBytes, err := decodeQRCodePNG(base64PNG)
	if err != nil {
		return "", err
	}

	safeFilename := ensurePNGFilename(filename)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("home directory の取得に失敗しました: %w", err)
	}

	targetDir := filepath.Join(homeDir, "Downloads", "oreno-tools-qr")
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return "", fmt.Errorf("保存先ディレクトリの作成に失敗しました: %w", err)
	}

	targetPath := filepath.Join(targetDir, safeFilename)
	if err := os.WriteFile(targetPath, pngBytes, 0o600); err != nil {
		return "", fmt.Errorf("PNGファイルの保存に失敗しました: %w", err)
	}

	return targetPath, nil
}

// SaveQRCodePNGToPath は base64 文字列の PNG を指定パスへ保存する。
func SaveQRCodePNGToPath(base64PNG string, targetPath string) (string, error) {
	pngBytes, err := decodeQRCodePNG(base64PNG)
	if err != nil {
		return "", err
	}

	trimmedPath := strings.TrimSpace(targetPath)
	if trimmedPath == "" {
		return "", fmt.Errorf("保存先が指定されていません")
	}

	if !strings.HasSuffix(strings.ToLower(trimmedPath), ".png") {
		trimmedPath += ".png"
	}

	targetDir := filepath.Dir(trimmedPath)
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return "", fmt.Errorf("保存先ディレクトリの作成に失敗しました: %w", err)
	}

	if err := os.WriteFile(trimmedPath, pngBytes, 0o600); err != nil {
		return "", fmt.Errorf("PNGファイルの保存に失敗しました: %w", err)
	}

	return trimmedPath, nil
}
