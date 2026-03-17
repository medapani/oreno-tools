package backend

import (
	"testing"
)

func TestBase64Encode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		urlSafe bool
		want    string
	}{
		{
			name:    "標準エンコード - 基本文字列",
			input:   "Hello, World!",
			urlSafe: false,
			want:    "SGVsbG8sIFdvcmxkIQ==",
		},
		{
			name:    "標準エンコード - 空文字列",
			input:   "",
			urlSafe: false,
			want:    "",
		},
		{
			name:    "URLセーフエンコード - 基本文字列",
			input:   "Hello, World!",
			urlSafe: true,
			want:    "SGVsbG8sIFdvcmxkIQ==",
		},
		{
			name:    "URLセーフエンコード - 特殊文字を含む",
			input:   "\xfb\xff\xfc",
			urlSafe: true,
			want:    "-__8",
		},
		{
			name:    "標準エンコード - 特殊文字を含む",
			input:   "\xfb\xff\xfc",
			urlSafe: false,
			want:    "+//8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Base64Encode(tt.input, tt.urlSafe)
			if got != tt.want {
				t.Errorf("Base64Encode(%q, %v) = %q, want %q", tt.input, tt.urlSafe, got, tt.want)
			}
		})
	}
}

func TestBase64Decode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		urlSafe bool
		want    string
		wantErr bool
	}{
		{
			name:    "標準デコード - 基本文字列",
			input:   "SGVsbG8sIFdvcmxkIQ==",
			urlSafe: false,
			want:    "Hello, World!",
			wantErr: false,
		},
		{
			name:    "標準デコード - 空文字列",
			input:   "",
			urlSafe: false,
			want:    "",
			wantErr: false,
		},
		{
			name:    "URLセーフデコード - URLセーフ文字",
			input:   "-__8",
			urlSafe: true,
			want:    "\xfb\xff\xfc",
			wantErr: false,
		},
		{
			name:    "標準デコード - 前後の空白を無視",
			input:   "  SGVsbG8sIFdvcmxkIQ==  ",
			urlSafe: false,
			want:    "Hello, World!",
			wantErr: false,
		},
		{
			name:    "標準デコード - 不正な入力",
			input:   "invalid!!!base64",
			urlSafe: false,
			want:    "",
			wantErr: true,
		},
		{
			// URLEncoding はパディング必須なので失敗し、RawURLEncoding フォールバックで成功
			name:    "URLセーフデコード - パディングなし (RawURLEncoding フォールバック)",
			input:   "YQ",
			urlSafe: true,
			want:    "a",
			wantErr: false,
		},
		{
			name:    "URLセーフデコード - 不正な入力 (両デコーダが失敗)",
			input:   "!!!invalid!!!",
			urlSafe: true,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Base64Decode(tt.input, tt.urlSafe)
			if (err != nil) != tt.wantErr {
				t.Errorf("Base64Decode(%q, %v) error = %v, wantErr %v", tt.input, tt.urlSafe, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("Base64Decode(%q, %v) = %q, want %q", tt.input, tt.urlSafe, got, tt.want)
			}
		})
	}
}

func TestBase64RoundTrip(t *testing.T) {
	inputs := []string{
		"Hello, World!",
		"",
		"special chars: !@#$%^&*()",
	}

	for _, input := range inputs {
		encoded := Base64Encode(input, false)
		decoded, err := Base64Decode(encoded, false)
		if err != nil {
			t.Errorf("standard round-trip decode error for %q: %v", input, err)
			continue
		}
		if decoded != input {
			t.Errorf("standard round-trip failed for %q: got %q", input, decoded)
		}

		encodedURL := Base64Encode(input, true)
		decodedURL, err := Base64Decode(encodedURL, true)
		if err != nil {
			t.Errorf("URL-safe round-trip decode error for %q: %v", input, err)
			continue
		}
		if decodedURL != input {
			t.Errorf("URL-safe round-trip failed for %q: got %q", input, decodedURL)
		}
	}
}
