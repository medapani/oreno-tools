package backend

import (
	"fmt"
	"net"
)

// CIDRCalculationResult はCIDR計算の結果を保持する構造体
type CIDRCalculationResult struct {
	NetworkAddress   string `json:"networkAddress"`
	BroadcastAddress string `json:"broadcastAddress"`
	SubnetMask       string `json:"subnetMask"`
	WildcardMask     string `json:"wildcardMask"`
	FirstHostAddress string `json:"firstHostAddress"`
	LastHostAddress  string `json:"lastHostAddress"`
	TotalHosts       int64  `json:"totalHosts"`
	UsableHosts      int64  `json:"usableHosts"`
	CIDR             string `json:"cidr"`
	BinarySubnetMask string `json:"binarySubnetMask"`
	IPClass          string `json:"ipClass"`
	IPType           string `json:"ipType"`
}

// CalculateCIDR はCIDR表記からネットワーク情報を計算する
func CalculateCIDR(cidr string) (CIDRCalculationResult, error) {
	// CIDR表記をパース
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return CIDRCalculationResult{}, fmt.Errorf("無効なCIDR表記です: %v", err)
	}

	// ネットワークアドレス
	networkAddr := ipNet.IP

	// サブネットマスク
	subnetMask := net.IP(ipNet.Mask)

	// ワイルドカードマスク
	wildcardMask := make(net.IP, len(subnetMask))
	for i := range subnetMask {
		wildcardMask[i] = ^subnetMask[i]
	}

	// ブロードキャストアドレス
	broadcastAddr := make(net.IP, len(networkAddr))
	for i := range networkAddr {
		broadcastAddr[i] = networkAddr[i] | wildcardMask[i]
	}

	// 最初のホストアドレス
	firstHost := make(net.IP, len(networkAddr))
	copy(firstHost, networkAddr)
	firstHost[len(firstHost)-1]++

	// 最後のホストアドレス
	lastHost := make(net.IP, len(broadcastAddr))
	copy(lastHost, broadcastAddr)
	lastHost[len(lastHost)-1]--

	// ホスト数を計算
	ones, bits := ipNet.Mask.Size()
	hostBits := bits - ones
	if hostBits < 0 {
		return CIDRCalculationResult{}, fmt.Errorf("invalid CIDR mask size")
	}
	if hostBits > 62 {
		return CIDRCalculationResult{}, fmt.Errorf("host count is too large to fit in int64")
	}
	totalHosts := int64(1) << hostBits
	usableHosts := totalHosts - 2
	if totalHosts <= 2 {
		usableHosts = 0
	}

	// バイナリ形式のサブネットマスク
	binaryMask := ""
	for i, b := range subnetMask {
		if i > 0 {
			binaryMask += "."
		}
		binaryMask += fmt.Sprintf("%08b", b)
	}

	// IPクラスを判定
	ipClass := determineIPClass(networkAddr)

	// IPタイプを判定（プライベート/パブリック）
	ipType := determineIPType(networkAddr)

	return CIDRCalculationResult{
		NetworkAddress:   networkAddr.String(),
		BroadcastAddress: broadcastAddr.String(),
		SubnetMask:       subnetMask.String(),
		WildcardMask:     wildcardMask.String(),
		FirstHostAddress: firstHost.String(),
		LastHostAddress:  lastHost.String(),
		TotalHosts:       totalHosts,
		UsableHosts:      usableHosts,
		CIDR:             cidr,
		BinarySubnetMask: binaryMask,
		IPClass:          ipClass,
		IPType:           ipType,
	}, nil
}

// determineIPClass はIPアドレスのクラスを判定する
func determineIPClass(ip net.IP) string {
	if ip == nil {
		return "Unknown"
	}

	ip = ip.To4()
	if ip == nil {
		return "IPv6"
	}

	firstOctet := ip[0]

	switch {
	case firstOctet >= 1 && firstOctet <= 126:
		return "A"
	case firstOctet >= 128 && firstOctet <= 191:
		return "B"
	case firstOctet >= 192 && firstOctet <= 223:
		return "C"
	case firstOctet >= 224 && firstOctet <= 239:
		return "D (Multicast)"
	case firstOctet >= 240:
		return "E (Reserved)"
	default:
		return "Unknown"
	}
}

// determineIPType はIPアドレスがプライベートかパブリックかを判定する
func determineIPType(ip net.IP) string {
	if ip == nil {
		return "Unknown"
	}

	ip = ip.To4()
	if ip == nil {
		return "IPv6"
	}

	// プライベートIPアドレスの範囲
	// 10.0.0.0/8
	// 172.16.0.0/12
	// 192.168.0.0/16

	if ip[0] == 10 {
		return "Private"
	}

	if ip[0] == 172 && ip[1] >= 16 && ip[1] <= 31 {
		return "Private"
	}

	if ip[0] == 192 && ip[1] == 168 {
		return "Private"
	}

	// ループバックアドレス
	if ip[0] == 127 {
		return "Loopback"
	}

	return "Public"
}
