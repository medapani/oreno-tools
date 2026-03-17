package backend

import (
	"math"
	"testing"
)

const floatTolerance = 1e-9

func floatEqual(a, b float64) bool {
	return math.Abs(a-b) < floatTolerance
}

func TestConvertUnit(t *testing.T) {
	tests := []struct {
		name    string
		size    float64
		unit    string
		wantB   float64
		wantMB  float64
		wantMiB float64
		wantErr bool
	}{
		{
			name:    "1 GB",
			size:    1,
			unit:    "GB",
			wantB:   1e9,
			wantMB:  1000,
			wantMiB: 1e9 / (1024 * 1024),
			wantErr: false,
		},
		{
			name:    "1 GiB",
			size:    1,
			unit:    "GiB",
			wantB:   1024 * 1024 * 1024,
			wantMB:  float64(1024*1024*1024) / 1e6,
			wantMiB: 1024,
			wantErr: false,
		},
		{
			name:    "1 KB",
			size:    1,
			unit:    "KB",
			wantB:   1000,
			wantErr: false,
		},
		{
			name:    "8 bits = 1 B",
			size:    8,
			unit:    "bits",
			wantB:   1,
			wantErr: false,
		},
		{
			name:    "不正な単位",
			size:    1,
			unit:    "XB",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertUnit(tt.size, tt.unit)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertUnit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if !floatEqual(got.B, tt.wantB) {
				t.Errorf("B = %v, want %v", got.B, tt.wantB)
			}
			if tt.wantMB != 0 && !floatEqual(got.MB, tt.wantMB) {
				t.Errorf("MB = %v, want %v", got.MB, tt.wantMB)
			}
			if tt.wantMiB != 0 && !floatEqual(got.MiB, tt.wantMiB) {
				t.Errorf("MiB = %v, want %v", got.MiB, tt.wantMiB)
			}
		})
	}
}

func TestConvertUnitBitsRelation(t *testing.T) {
	// 1 B = 8 bits の関係を検証
	got, err := ConvertUnit(1, "B")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !floatEqual(got.Bits, 8) {
		t.Errorf("1B Bits = %v, want 8", got.Bits)
	}
	if !floatEqual(got.B, 1) {
		t.Errorf("1B B = %v, want 1", got.B)
	}
}

func TestConvertDataTransferRate(t *testing.T) {
	tests := []struct {
		name        string
		speed       float64
		unit        string
		wantBPerSec float64
		wantErr     bool
	}{
		{
			name:        "1 GB/s",
			speed:       1,
			unit:        "GB/s",
			wantBPerSec: 1e9,
			wantErr:     false,
		},
		{
			name:        "1 GiB/s",
			speed:       1,
			unit:        "GiB/s",
			wantBPerSec: 1024 * 1024 * 1024,
			wantErr:     false,
		},
		{
			name:        "1000 Mbps",
			speed:       1000,
			unit:        "Mbps",
			wantBPerSec: 1000 * 1e6 / 8,
			wantErr:     false,
		},
		{
			name:    "不正な単位",
			speed:   1,
			unit:    "XB/s",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertDataTransferRate(tt.speed, tt.unit)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertDataTransferRate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if !floatEqual(got.BytesPerSec, tt.wantBPerSec) {
				t.Errorf("BytesPerSec = %v, want %v", got.BytesPerSec, tt.wantBPerSec)
			}
		})
	}
}

func TestCalculateTransferTime(t *testing.T) {
	tests := []struct {
		name        string
		dataSize    float64
		dataUnit    string
		speed       float64
		speedUnit   string
		wantSeconds float64
		wantErr     bool
	}{
		{
			name:        "1 GB を 1 GB/s で転送",
			dataSize:    1,
			dataUnit:    "GB",
			speed:       1,
			speedUnit:   "GB/s",
			wantSeconds: 1,
			wantErr:     false,
		},
		{
			name:        "100 MB を 10 MB/s で転送",
			dataSize:    100,
			dataUnit:    "MB",
			speed:       10,
			speedUnit:   "MB/s",
			wantSeconds: 10,
			wantErr:     false,
		},
		{
			name:        "1 TB を 1 GiB/s で転送",
			dataSize:    1,
			dataUnit:    "TB",
			speed:       1,
			speedUnit:   "GiB/s",
			wantSeconds: 1e12 / (1024 * 1024 * 1024),
			wantErr:     false,
		},
		{
			name:      "不正なデータ単位",
			dataSize:  1,
			dataUnit:  "XX",
			speed:     1,
			speedUnit: "MB/s",
			wantErr:   true,
		},
		{
			name:      "不正な速度単位",
			dataSize:  1,
			dataUnit:  "GB",
			speed:     1,
			speedUnit: "XX/s",
			wantErr:   true,
		},
		{
			name:      "速度ゼロ",
			dataSize:  1,
			dataUnit:  "GB",
			speed:     0,
			speedUnit: "GB/s",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateTransferTime(tt.dataSize, tt.dataUnit, tt.speed, tt.speedUnit)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateTransferTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if !floatEqual(got.Seconds, tt.wantSeconds) {
				t.Errorf("Seconds = %v, want %v", got.Seconds, tt.wantSeconds)
			}
			if !floatEqual(got.Minutes, tt.wantSeconds/60) {
				t.Errorf("Minutes = %v, want %v", got.Minutes, tt.wantSeconds/60)
			}
			if !floatEqual(got.Hours, tt.wantSeconds/3600) {
				t.Errorf("Hours = %v, want %v", got.Hours, tt.wantSeconds/3600)
			}
		})
	}
}
