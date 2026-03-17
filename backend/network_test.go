package backend

import (
	"net"
	"testing"
)

func TestCalculateCIDR(t *testing.T) {
	tests := []struct {
		name            string
		cidr            string
		wantNetwork     string
		wantBroadcast   string
		wantSubnetMask  string
		wantFirstHost   string
		wantLastHost    string
		wantTotalHosts  int64
		wantUsableHosts int64
		wantIPClass     string
		wantIPType      string
		wantErr         bool
	}{
		{
			name:            "192.168.1.0/24 (クラスC プライベート)",
			cidr:            "192.168.1.0/24",
			wantNetwork:     "192.168.1.0",
			wantBroadcast:   "192.168.1.255",
			wantSubnetMask:  "255.255.255.0",
			wantFirstHost:   "192.168.1.1",
			wantLastHost:    "192.168.1.254",
			wantTotalHosts:  256,
			wantUsableHosts: 254,
			wantIPClass:     "C",
			wantIPType:      "Private",
			wantErr:         false,
		},
		{
			name:            "10.0.0.0/8 (クラスA プライベート)",
			cidr:            "10.0.0.0/8",
			wantNetwork:     "10.0.0.0",
			wantBroadcast:   "10.255.255.255",
			wantSubnetMask:  "255.0.0.0",
			wantTotalHosts:  16777216,
			wantUsableHosts: 16777214,
			wantIPClass:     "A",
			wantIPType:      "Private",
			wantErr:         false,
		},
		{
			name:            "172.16.0.0/12 (クラスB プライベート)",
			cidr:            "172.16.0.0/12",
			wantNetwork:     "172.16.0.0",
			wantBroadcast:   "172.31.255.255",
			wantSubnetMask:  "255.240.0.0",
			wantTotalHosts:  1048576,
			wantUsableHosts: 1048574,
			wantIPClass:     "B",
			wantIPType:      "Private",
			wantErr:         false,
		},
		{
			name:            "8.8.8.0/24 (パブリック)",
			cidr:            "8.8.8.0/24",
			wantNetwork:     "8.8.8.0",
			wantBroadcast:   "8.8.8.255",
			wantTotalHosts:  256,
			wantUsableHosts: 254,
			wantIPClass:     "A",
			wantIPType:      "Public",
			wantErr:         false,
		},
		{
			name:            "/30 ネットワーク (usable hosts = 2)",
			cidr:            "192.168.1.0/30",
			wantNetwork:     "192.168.1.0",
			wantBroadcast:   "192.168.1.3",
			wantTotalHosts:  4,
			wantUsableHosts: 2,
			wantErr:         false,
		},
		{
			name:            "/32 ホストルート",
			cidr:            "192.168.1.1/32",
			wantNetwork:     "192.168.1.1",
			wantTotalHosts:  1,
			wantUsableHosts: 0,
			wantErr:         false,
		},
		{
			name:    "不正なCIDR",
			cidr:    "invalid",
			wantErr: true,
		},
		{
			name:    "不正なIPアドレス",
			cidr:    "999.999.999.999/24",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateCIDR(tt.cidr)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateCIDR(%q) error = %v, wantErr %v", tt.cidr, err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if tt.wantNetwork != "" && got.NetworkAddress != tt.wantNetwork {
				t.Errorf("NetworkAddress = %q, want %q", got.NetworkAddress, tt.wantNetwork)
			}
			if tt.wantBroadcast != "" && got.BroadcastAddress != tt.wantBroadcast {
				t.Errorf("BroadcastAddress = %q, want %q", got.BroadcastAddress, tt.wantBroadcast)
			}
			if tt.wantSubnetMask != "" && got.SubnetMask != tt.wantSubnetMask {
				t.Errorf("SubnetMask = %q, want %q", got.SubnetMask, tt.wantSubnetMask)
			}
			if tt.wantFirstHost != "" && got.FirstHostAddress != tt.wantFirstHost {
				t.Errorf("FirstHostAddress = %q, want %q", got.FirstHostAddress, tt.wantFirstHost)
			}
			if tt.wantLastHost != "" && got.LastHostAddress != tt.wantLastHost {
				t.Errorf("LastHostAddress = %q, want %q", got.LastHostAddress, tt.wantLastHost)
			}
			if got.TotalHosts != tt.wantTotalHosts {
				t.Errorf("TotalHosts = %d, want %d", got.TotalHosts, tt.wantTotalHosts)
			}
			if got.UsableHosts != tt.wantUsableHosts {
				t.Errorf("UsableHosts = %d, want %d", got.UsableHosts, tt.wantUsableHosts)
			}
			if tt.wantIPClass != "" && got.IPClass != tt.wantIPClass {
				t.Errorf("IPClass = %q, want %q", got.IPClass, tt.wantIPClass)
			}
			if tt.wantIPType != "" && got.IPType != tt.wantIPType {
				t.Errorf("IPType = %q, want %q", got.IPType, tt.wantIPType)
			}
		})
	}
}

func TestDetermineIPClass(t *testing.T) {
	tests := []struct {
		name string
		ip   net.IP
		want string
	}{
		{"クラスA (1.x)", net.ParseIP("1.0.0.0").To4(), "A"},
		{"クラスA (100.x)", net.ParseIP("100.0.0.0").To4(), "A"},
		{"クラスB (128.x)", net.ParseIP("128.0.0.0").To4(), "B"},
		{"クラスB (191.x)", net.ParseIP("191.255.0.0").To4(), "B"},
		{"クラスC (192.x)", net.ParseIP("192.0.0.0").To4(), "C"},
		{"クラスC (223.x)", net.ParseIP("223.255.255.0").To4(), "C"},
		{"クラスD マルチキャスト", net.ParseIP("224.0.0.0").To4(), "D (Multicast)"},
		{"クラスE 予約済み", net.ParseIP("240.0.0.0").To4(), "E (Reserved)"},
		{"nil IP", nil, "Unknown"},
		// To4() を呼ばずに IPv6 アドレスを渡す → 関数内の ip.To4() が nil を返す
		{"IPv6アドレス", net.ParseIP("2001:db8::1"), "IPv6"},
		// 127.x は 1-126 の範囲外 → default → "Unknown"
		{"127.x ループバック → Unknown", net.ParseIP("127.0.0.1").To4(), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := determineIPClass(tt.ip)
			if got != tt.want {
				t.Errorf("determineIPClass(%v) = %q, want %q", tt.ip, got, tt.want)
			}
		})
	}
}

func TestDetermineIPType(t *testing.T) {
	tests := []struct {
		name string
		ip   net.IP
		want string
	}{
		{"10.x.x.x プライベート", net.ParseIP("10.0.0.1").To4(), "Private"},
		{"172.16.x.x プライベート", net.ParseIP("172.16.0.1").To4(), "Private"},
		{"172.31.x.x プライベート", net.ParseIP("172.31.255.255").To4(), "Private"},
		{"192.168.x.x プライベート", net.ParseIP("192.168.0.1").To4(), "Private"},
		{"8.8.8.8 パブリック", net.ParseIP("8.8.8.8").To4(), "Public"},
		{"127.0.0.1 ループバック", net.ParseIP("127.0.0.1").To4(), "Loopback"},
		{"nil IP", nil, "Unknown"},
		// To4() を呼ばずに IPv6 アドレスを渡す → 関数内の ip.To4() が nil を返す
		{"IPv6アドレス", net.ParseIP("2001:db8::1"), "IPv6"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := determineIPType(tt.ip)
			if got != tt.want {
				t.Errorf("determineIPType(%v) = %q, want %q", tt.ip, got, tt.want)
			}
		})
	}
}
