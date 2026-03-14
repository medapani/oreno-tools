package backend

import (
	"fmt"
)

// ConversionResult は容量変換の結果を保持する構造体
type ConversionResult struct {
	B      float64 `json:"B"`
	KB     float64 `json:"KB"`
	MB     float64 `json:"MB"`
	GB     float64 `json:"GB"`
	TB     float64 `json:"TB"`
	KiB    float64 `json:"KiB"`
	MiB    float64 `json:"MiB"`
	GiB    float64 `json:"GiB"`
	TiB    float64 `json:"TiB"`
	Bits   float64 `json:"bits"`
	Kbits  float64 `json:"Kbits"`
	Mbits  float64 `json:"Mbits"`
	Gbits  float64 `json:"Gbits"`
	Tbits  float64 `json:"Tbits"`
	Kibits float64 `json:"Kibits"`
	Mibits float64 `json:"Mibits"`
	Gibits float64 `json:"Gibits"`
	Tibits float64 `json:"Tibits"`
}

// DataTransferRateResult はデータ転送速度変換の結果を保持する構造体
type DataTransferRateResult struct {
	BytesPerSec float64 `json:"B/s"`
	KBPerSec    float64 `json:"KB/s"`
	MBPerSec    float64 `json:"MB/s"`
	GBPerSec    float64 `json:"GB/s"`
	TBPerSec    float64 `json:"TB/s"`
	KiBPerSec   float64 `json:"KiB/s"`
	MiBPerSec   float64 `json:"MiB/s"`
	GiBPerSec   float64 `json:"GiB/s"`
	TiBPerSec   float64 `json:"TiB/s"`
	BitsPerSec  float64 `json:"bit/s"`
	KbitPerSec  float64 `json:"Kbit/s"`
	MbitPerSec  float64 `json:"Mbit/s"`
	GbitPerSec  float64 `json:"Gbit/s"`
	TbitPerSec  float64 `json:"Tbit/s"`
	KibitPerSec float64 `json:"Kibit/s"`
	MibitPerSec float64 `json:"Mibit/s"`
	GibitPerSec float64 `json:"Gibit/s"`
	TibitPerSec float64 `json:"Tibit/s"`
}

// TransferTimeResult はデータ転送時間の結果を保持する構造体
type TransferTimeResult struct {
	Seconds float64 `json:"seconds"`
	Minutes float64 `json:"minutes"`
	Hours   float64 `json:"hours"`
	Days    float64 `json:"days"`
}

// 単位変換用の基本係数
const (
	KB  = 1e3
	MB  = 1e6
	GB  = 1e9
	TB  = 1e12
	KiB = 1024
	MiB = 1024 * KiB
	GiB = 1024 * MiB
	TiB = 1024 * GiB
)

var convertUnitToBytes = map[string]float64{
	"B":      1,
	"KB":     KB,
	"MB":     MB,
	"GB":     GB,
	"TB":     TB,
	"KiB":    KiB,
	"MiB":    MiB,
	"GiB":    GiB,
	"TiB":    TiB,
	"bits":   1.0 / 8,
	"Kbits":  KB / 8,
	"Mbits":  MB / 8,
	"Gbits":  GB / 8,
	"Tbits":  TB / 8,
	"Kibits": KiB / 8,
	"Mibits": MiB / 8,
	"Gibits": GiB / 8,
	"Tibits": TiB / 8,
}

var dataSizeUnitToBytes = map[string]float64{
	"B":   1,
	"KB":  KB,
	"MB":  MB,
	"GB":  GB,
	"TB":  TB,
	"KiB": KiB,
	"MiB": MiB,
	"GiB": GiB,
	"TiB": TiB,
}

var transferRateToBytesPerSec = map[string]float64{
	"B/s":   1,
	"KB/s":  KB,
	"MB/s":  MB,
	"GB/s":  GB,
	"TB/s":  TB,
	"KiB/s": KiB,
	"MiB/s": MiB,
	"GiB/s": GiB,
	"TiB/s": TiB,
	"bps":   1.0 / 8,
	"Kbps":  KB / 8,
	"Mbps":  MB / 8,
	"Gbps":  GB / 8,
	"Tbps":  TB / 8,
	"Kibps": KiB / 8,
	"Mibps": MiB / 8,
	"Gibps": GiB / 8,
	"Tibps": TiB / 8,
}

func ConvertUnit(size float64, unit string) (ConversionResult, error) {
	factor, ok := convertUnitToBytes[unit]
	if !ok {
		return ConversionResult{}, fmt.Errorf("invalid unit: %s", unit)
	}

	baseSize := size * factor

	return ConversionResult{
		B:      baseSize,
		KB:     baseSize / KB,
		MB:     baseSize / MB,
		GB:     baseSize / GB,
		TB:     baseSize / TB,
		KiB:    baseSize / KiB,
		MiB:    baseSize / MiB,
		GiB:    baseSize / GiB,
		TiB:    baseSize / TiB,
		Bits:   baseSize * 8,
		Kbits:  baseSize * 8 / KB,
		Mbits:  baseSize * 8 / MB,
		Gbits:  baseSize * 8 / GB,
		Tbits:  baseSize * 8 / TB,
		Kibits: baseSize * 8 / KiB,
		Mibits: baseSize * 8 / MiB,
		Gibits: baseSize * 8 / GiB,
		Tibits: baseSize * 8 / TiB,
	}, nil
}

// ConvertDataTransferRate はデータ転送速度を変換する
func ConvertDataTransferRate(speed float64, unit string) (DataTransferRateResult, error) {
	speedFactor, ok := transferRateToBytesPerSec[unit]
	if !ok {
		return DataTransferRateResult{}, fmt.Errorf("invalid transfer rate unit: %s", unit)
	}
	baseBytesPerSec := speed * speedFactor

	return DataTransferRateResult{
		BytesPerSec: baseBytesPerSec,
		KBPerSec:    baseBytesPerSec / KB,
		MBPerSec:    baseBytesPerSec / MB,
		GBPerSec:    baseBytesPerSec / GB,
		TBPerSec:    baseBytesPerSec / TB,
		KiBPerSec:   baseBytesPerSec / KiB,
		MiBPerSec:   baseBytesPerSec / MiB,
		GiBPerSec:   baseBytesPerSec / GiB,
		TiBPerSec:   baseBytesPerSec / TiB,
		BitsPerSec:  baseBytesPerSec * 8,
		KbitPerSec:  baseBytesPerSec * 8 / KB,
		MbitPerSec:  baseBytesPerSec * 8 / MB,
		GbitPerSec:  baseBytesPerSec * 8 / GB,
		TbitPerSec:  baseBytesPerSec * 8 / TB,
		KibitPerSec: baseBytesPerSec * 8 / KiB,
		MibitPerSec: baseBytesPerSec * 8 / MiB,
		GibitPerSec: baseBytesPerSec * 8 / GiB,
		TibitPerSec: baseBytesPerSec * 8 / TiB,
	}, nil
}

// CalculateTransferTime はデータ転送時間を計算する
// dataSize: データ容量（バイト単位）
// dataUnit: データ容量の単位 (B, KB, MB, GB, TB, KiB, MiB, GiB, TiB)
// speed: 転送速度（数値）
// speedUnit: 転送速度の単位 (B/s, KB/s, MB/s, GB/s, TB/s, KiB/s, MiB/s, GiB/s, TiB/s, bit/s, Kbit/s, Mbit/s, Gbit/s, Tbit/s, Kibit/s, Mibit/s, Gibit/s, Tibit/s)
func CalculateTransferTime(dataSize float64, dataUnit string, speed float64, speedUnit string) (TransferTimeResult, error) {
	dataFactor, ok := dataSizeUnitToBytes[dataUnit]
	if !ok {
		return TransferTimeResult{}, fmt.Errorf("invalid data unit: %s", dataUnit)
	}
	baseSizeBytes := dataSize * dataFactor

	speedFactor, ok := transferRateToBytesPerSec[speedUnit]
	if !ok {
		return TransferTimeResult{}, fmt.Errorf("invalid speed unit: %s", speedUnit)
	}
	baseBytesPerSec := speed * speedFactor

	// 転送速度が0の場合はエラー
	if baseBytesPerSec == 0 {
		return TransferTimeResult{}, fmt.Errorf("transfer speed cannot be zero")
	}

	// 転送時間を秒で計算
	seconds := baseSizeBytes / baseBytesPerSec

	return TransferTimeResult{
		Seconds: seconds,
		Minutes: seconds / 60,
		Hours:   seconds / 3600,
		Days:    seconds / 86400,
	}, nil
}
