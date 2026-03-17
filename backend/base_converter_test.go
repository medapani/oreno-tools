package backend

import (
	"testing"
)

func TestConvertBaseValue(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		base        int
		bitWidth    int
		signed      bool
		wantDecimal string
		wantHex     string
		wantBinary  string
		wantErr     bool
	}{
		{
			name:        "10進数 255 を 8bit unsigned",
			input:       "255",
			base:        10,
			bitWidth:    8,
			signed:      false,
			wantDecimal: "255",
			wantHex:     "FF",
			wantBinary:  "11111111",
			wantErr:     false,
		},
		{
			name:        "16進数 FF を 8bit unsigned",
			input:       "FF",
			base:        16,
			bitWidth:    8,
			signed:      false,
			wantDecimal: "255",
			wantHex:     "FF",
			wantBinary:  "11111111",
			wantErr:     false,
		},
		{
			name:        "2進数 1010 を 8bit unsigned",
			input:       "1010",
			base:        2,
			bitWidth:    8,
			signed:      false,
			wantDecimal: "10",
			wantHex:     "A",
			wantBinary:  "00001010",
			wantErr:     false,
		},
		{
			name:        "10進数 -1 を 8bit signed",
			input:       "-1",
			base:        10,
			bitWidth:    8,
			signed:      true,
			wantDecimal: "-1",
			wantHex:     "FF",
			wantBinary:  "11111111",
			wantErr:     false,
		},
		{
			name:        "オーバーフロー clamp: 256 を 8bit unsigned",
			input:       "256",
			base:        10,
			bitWidth:    8,
			signed:      false,
			wantDecimal: "255",
			wantErr:     false,
		},
		{
			name:     "不正なbase",
			input:    "10",
			base:     8,
			bitWidth: 8,
			signed:   false,
			wantErr:  true,
		},
		{
			name:     "空の入力",
			input:    "",
			base:     10,
			bitWidth: 8,
			signed:   false,
			wantErr:  true,
		},
		{
			name:     "不正なbitWidth",
			input:    "10",
			base:     10,
			bitWidth: 7,
			signed:   false,
			wantErr:  true,
		},
		{
			name:        "16bit unsigned",
			input:       "65535",
			base:        10,
			bitWidth:    16,
			signed:      false,
			wantDecimal: "65535",
			wantErr:     false,
		},
		{
			name:        "32bit unsigned",
			input:       "4294967295",
			base:        10,
			bitWidth:    32,
			signed:      false,
			wantDecimal: "4294967295",
			wantErr:     false,
		},
		{
			name:        "64bit unsigned",
			input:       "1",
			base:        10,
			bitWidth:    64,
			signed:      false,
			wantDecimal: "1",
			wantErr:     false,
		},
		{
			// 負数を unsigned に渡すと clampBigInt が 0 にクランプ
			name:        "負数 unsigned -> 0 にクランプ",
			input:       "-5",
			base:        10,
			bitWidth:    8,
			signed:      false,
			wantDecimal: "0",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertBaseValue(tt.input, tt.base, tt.bitWidth, tt.signed)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertBaseValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if tt.wantDecimal != "" && got.Decimal != tt.wantDecimal {
				t.Errorf("Decimal = %q, want %q", got.Decimal, tt.wantDecimal)
			}
			if tt.wantHex != "" && got.Hex != tt.wantHex {
				t.Errorf("Hex = %q, want %q", got.Hex, tt.wantHex)
			}
			if tt.wantBinary != "" && got.Binary != tt.wantBinary {
				t.Errorf("Binary = %q, want %q", got.Binary, tt.wantBinary)
			}
		})
	}
}

func TestCalculateBaseExpression(t *testing.T) {
	tests := []struct {
		name        string
		aInput      string
		aBase       int
		operator    string
		bInput      string
		bBase       int
		wantDecimal string
		wantErr     bool
	}{
		{
			name: "加算 10 + 5", aInput: "10", aBase: 10, operator: "+", bInput: "5", bBase: 10,
			wantDecimal: "15",
		},
		{
			name: "減算 10 - 3", aInput: "10", aBase: 10, operator: "-", bInput: "3", bBase: 10,
			wantDecimal: "7",
		},
		{
			name: "乗算 6 * 7", aInput: "6", aBase: 10, operator: "*", bInput: "7", bBase: 10,
			wantDecimal: "42",
		},
		{
			name: "除算 10 / 2", aInput: "10", aBase: 10, operator: "/", bInput: "2", bBase: 10,
			wantDecimal: "5",
		},
		{
			name: "剰余 10 % 3", aInput: "10", aBase: 10, operator: "%", bInput: "3", bBase: 10,
			wantDecimal: "1",
		},
		{
			name: "AND FF & 0F (hex)", aInput: "FF", aBase: 16, operator: "&", bInput: "F", bBase: 16,
			wantDecimal: "15",
		},
		{
			name: "OR 10 | 5", aInput: "10", aBase: 10, operator: "|", bInput: "5", bBase: 10,
			wantDecimal: "15",
		},
		{
			name: "XOR 15 ^ 9", aInput: "15", aBase: 10, operator: "^", bInput: "9", bBase: 10,
			wantDecimal: "6",
		},
		{
			name: "左シフト 1 << 4", aInput: "1", aBase: 10, operator: "<<", bInput: "4", bBase: 10,
			wantDecimal: "16",
		},
		{
			name: "右シフト 16 >> 2", aInput: "16", aBase: 10, operator: ">>", bInput: "2", bBase: 10,
			wantDecimal: "4",
		},
		{
			name: "ゼロ除算", aInput: "10", aBase: 10, operator: "/", bInput: "0", bBase: 10,
			wantErr: true,
		},
		{
			name: "ゼロ剰余", aInput: "10", aBase: 10, operator: "%", bInput: "0", bBase: 10,
			wantErr: true,
		},
		{
			name: "不正な演算子", aInput: "10", aBase: 10, operator: "**", bInput: "2", bBase: 10,
			wantErr: true,
		},
		{
			name:   "異なるbase間の計算: 16進A + 2進1010",
			aInput: "A", aBase: 16, operator: "+", bInput: "1010", bBase: 2,
			wantDecimal: "20",
		},
		{name: "Aが不正な入力", aInput: "", aBase: 10, operator: "+", bInput: "1", bBase: 10, wantErr: true},
		{name: "Bが不正な入力", aInput: "1", aBase: 10, operator: "+", bInput: "", bBase: 10, wantErr: true},
		{
			// formatRawHex / formatRawBinary の負値パスを通過
			name: "負の結果 3 - 10 = -7", aInput: "3", aBase: 10, operator: "-", bInput: "10", bBase: 10,
			wantDecimal: "-7",
		},
		{name: "左シフト 負のカウント", aInput: "1", aBase: 10, operator: "<<", bInput: "-1", bBase: 10, wantErr: true},
		{name: "右シフト 負のカウント", aInput: "16", aBase: 10, operator: ">>", bInput: "-1", bBase: 10, wantErr: true},
		{name: "左シフト 巨大カウント (>1<<30)", aInput: "1", aBase: 10, operator: "<<", bInput: "2000000000", bBase: 10, wantErr: true},
		{name: "右シフト 巨大カウント (>1<<30)", aInput: "1", aBase: 10, operator: ">>", bInput: "2000000000", bBase: 10, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateBaseExpression(tt.aInput, tt.aBase, tt.operator, tt.bInput, tt.bBase)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateBaseExpression() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.Decimal != tt.wantDecimal {
				t.Errorf("Decimal = %q, want %q", got.Decimal, tt.wantDecimal)
			}
		})
	}
}

func TestParseBigIntByBase(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		base    int
		wantErr bool
	}{
		{"正常な10進数", "42", 10, false},
		{"正常な16進数", "FF", 16, false},
		{"正常な2進数", "1010", 2, false},
		{"マイナス記号のみ", "-", 10, true},
		{"空文字", "", 10, true},
		{"サポート外のbase (8進) ", "10", 8, true},
		{"パース不可能な入力", "ZZ", 16, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseBigIntByBase(tt.input, tt.base)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseBigIntByBase(%q, %d) error = %v, wantErr %v", tt.input, tt.base, err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("parseBigIntByBase() returned nil without error")
			}
		})
	}
}

func TestGroupBinary(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"空文字", "", ""},
		{"4ビット", "1010", "1010"},
		{"8ビット", "11111111", "1111 1111"},
		{"16ビット", "1010101010101010", "1010 1010 1010 1010"},
		{"端数あり", "10101", "1 0101"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := groupBinary(tt.input)
			if got != tt.want {
				t.Errorf("groupBinary(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestGroupBinaryWithSign(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"正の値", "10101010", "1010 1010"},
		{"負の値", "-10101010", "-1010 1010"},
		{"空文字", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := groupBinaryWithSign(tt.input)
			if got != tt.want {
				t.Errorf("groupBinaryWithSign(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
