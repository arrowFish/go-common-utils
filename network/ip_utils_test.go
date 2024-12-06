package network

import (
	"net"
	"testing"
)

func TestGetLocalIPv4(t *testing.T) {
	// 保存原始函数
	originalInterfaceAddrs := interfaceAddrs
	originalGetInterfaces := getInterfaces
	defer func() {
		interfaceAddrs = originalInterfaceAddrs
		getInterfaces = originalGetInterfaces
	}()

	tests := []struct {
		name     string
		setup    func()
		expected string
	}{
		{
			name: "正常情况-返回公网IP",
			setup: func() {
				getInterfaces = func() ([]net.Interface, error) {
					return []net.Interface{{
						Index: 1,
						Name:  "eth0",
						Flags: net.FlagUp,
					}}, nil
				}
				interfaceAddrs = func(i *net.Interface) ([]net.Addr, error) {
					return []net.Addr{
						&net.IPNet{
							IP:   net.ParseIP("203.0.113.1"),
							Mask: net.CIDRMask(24, 32),
						},
					}, nil
				}
			},
			expected: "203.0.113.1",
		},
		{
			name: "正常情况-返回内网IP",
			setup: func() {
				getInterfaces = func() ([]net.Interface, error) {
					return []net.Interface{{
						Index: 1,
						Name:  "eth0",
						Flags: net.FlagUp,
					}}, nil
				}
				interfaceAddrs = func(i *net.Interface) ([]net.Addr, error) {
					return []net.Addr{
						&net.IPNet{
							IP:   net.ParseIP("192.168.1.100"),
							Mask: net.CIDRMask(24, 32),
						},
					}, nil
				}
			},
			expected: "192.168.1.100",
		},
		{ // 没有多网卡的环境，暂未测试通过
			name: "多网卡-优先返回公网IP",
			setup: func() {
				getInterfaces = func() ([]net.Interface, error) {
					return []net.Interface{
						{
							Index: 1,
							Name:  "eth0",
							Flags: net.FlagUp,
						},
						{
							Index: 2,
							Name:  "eth1",
							Flags: net.FlagUp,
						},
					}, nil
				}
				interfaceAddrs = func(i *net.Interface) ([]net.Addr, error) {
					if i.Name == "eth0" {
						return []net.Addr{
							&net.IPNet{
								IP:   net.ParseIP("192.168.1.100"),
								Mask: net.CIDRMask(24, 32),
							},
						}, nil
					}
					return []net.Addr{
						&net.IPNet{
							IP:   net.ParseIP("203.0.113.1"),
							Mask: net.CIDRMask(24, 32),
						},
					}, nil
				}
			},
			expected: "203.0.113.1",
		},
		{
			name: "错误情况-接口获取失败",
			setup: func() {
				getInterfaces = func() ([]net.Interface, error) {
					return nil, &net.OpError{Op: "mock", Net: "mock", Err: net.UnknownNetworkError("mock")}
				}
			},
			expected: "127.0.0.1",
		},
		{
			name: "特殊情况-跳过docker接口",
			setup: func() {
				getInterfaces = func() ([]net.Interface, error) {
					return []net.Interface{
						{
							Index: 1,
							Name:  "docker0",
							Flags: net.FlagUp,
						},
						{
							Index: 2,
							Name:  "eth0",
							Flags: net.FlagUp,
						},
					}, nil
				}
				interfaceAddrs = func(i *net.Interface) ([]net.Addr, error) {
					if i.Name == "docker0" {
						return []net.Addr{
							&net.IPNet{
								IP:   net.ParseIP("172.17.0.1"),
								Mask: net.CIDRMask(24, 32),
							},
						}, nil
					}
					return []net.Addr{
						&net.IPNet{
							IP:   net.ParseIP("192.168.1.100"),
							Mask: net.CIDRMask(24, 32),
						},
					}, nil
				}
			},
			expected: "192.168.1.100",
		},
		{
			name: "特殊情况-跳过回环地址",
			setup: func() {
				getInterfaces = func() ([]net.Interface, error) {
					return []net.Interface{{
						Index: 1,
						Name:  "lo",
						Flags: net.FlagUp | net.FlagLoopback,
					}}, nil
				}
				interfaceAddrs = func(i *net.Interface) ([]net.Addr, error) {
					return []net.Addr{
						&net.IPNet{
							IP:   net.ParseIP("127.0.0.1"),
							Mask: net.CIDRMask(8, 32),
						},
					}, nil
				}
			},
			expected: "127.0.0.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			result := GetLocalIPv4()
			if result != tt.expected {
				t.Errorf("GetLocalIPv4() = %v, want %v", result, tt.expected)
			}
		})
	}
}
