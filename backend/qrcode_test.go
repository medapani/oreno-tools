package backend

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestGenerateQRCode_Valid(t *testing.T) {
	result, err := GenerateQRCode("https://example.com", 256, "M")
	if err != nil {
		t.Fatalf("予期しないエラー: %v", err)
	}
	if result == "" {
		t.Fatal("結果が空です")
	}
	png, err := base64.StdEncoding.DecodeString(result)
	if err != nil {
		t.Fatalf("base64 デコード失敗: %v", err)
	}
	// PNG シグネチャ確認
	if len(png) < 8 || string(png[:8]) != "\x89PNG\r\n\x1a\n" {
		t.Fatal("PNG シグネチャが一致しません")
	}
}

func TestGenerateQRCode_AllRecoveryLevels(t *testing.T) {
	for _, level := range []string{"L", "M", "Q", "H"} {
		t.Run(level, func(t *testing.T) {
			_, err := GenerateQRCode("test", 128, level)
			if err != nil {
				t.Errorf("レベル %s でエラー: %v", level, err)
			}
		})
	}
}

func TestGenerateQRCode_EmptyText(t *testing.T) {
	_, err := GenerateQRCode("", 256, "M")
	if err == nil {
		t.Fatal("空テキストでエラーが返されるべきです")
	}
}

func TestGenerateQRCode_WhitespaceOnly(t *testing.T) {
	_, err := GenerateQRCode("   ", 256, "M")
	if err == nil {
		t.Fatal("空白のみのテキストでエラーが返されるべきです")
	}
}

func TestGenerateQRCode_InvalidRecoveryLevel(t *testing.T) {
	_, err := GenerateQRCode("test", 256, "X")
	if err == nil {
		t.Fatal("不正なレベルでエラーが返されるべきです")
	}
}

func TestGenerateQRCode_InvalidSize(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{"zero", 0},
		{"negative", -1},
		{"too small", 63},
		{"too large", 2049},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := GenerateQRCode("test", tc.size, "M")
			if err == nil {
				t.Errorf("サイズ %d でエラーが返されるべきです", tc.size)
			}
		})
	}
}

func TestGenerateQRCode_UTF8Text(t *testing.T) {
	_, err := GenerateQRCode("こんにちは世界 🌍", 256, "M")
	if err != nil {
		t.Fatalf("UTF-8 テキストで予期しないエラー: %v", err)
	}
}

func TestGenerateQRCode_SizeBoundary(t *testing.T) {
	for _, size := range []int{64, 2048} {
		t.Run(strings.Repeat("size", 1), func(t *testing.T) {
			_, err := GenerateQRCode("test", size, "M")
			if err != nil {
				t.Errorf("サイズ %d でエラー: %v", size, err)
			}
		})
	}
}
