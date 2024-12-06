package network

import (
	"bytes"
	"net"
	"strings"
)

// 定义接口获取器
var interfaceAddrs = func(i *net.Interface) ([]net.Addr, error) {
	return i.Addrs()
}

// 定义网络接口获取器
var getInterfaces = func() ([]net.Interface, error) {
	return net.Interfaces()
}

func GetLocalIPv4() string {
	interfaces, err := getInterfaces()
	if err != nil {
		return "127.0.0.1"
	}

	var candidateIPs []string

	for _, iface := range interfaces {
		// 跳过 loopback 和 down 状态的接口
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 ||
			strings.Contains(iface.Name, "docker") || strings.Contains(iface.Name, "veth") {
			continue
		}

		addrs, err := interfaceAddrs(&iface)
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			ip4 := ipNet.IP.To4()
			if ip4 == nil || ip4.IsLoopback() {
				continue
			}

			// 排除特殊地址段
			if ip4[0] == 169 && ip4[1] == 254 { // 排除 link-local 地址
				continue
			}

			// 优先选择非内网地址
			if !isPrivateIP(ip4) {
				return ip4.String() // 立即返回公网 IP
			}

			candidateIPs = append(candidateIPs, ip4.String())
		}
	}

	// 如果有内网地址，返回第一个内网地址
	if len(candidateIPs) > 0 {
		return candidateIPs[0]
	}

	return "127.0.0.1"
}

// 判断是否为内网 IP
func isPrivateIP(ip net.IP) bool {
	privateIPBlocks := []struct {
		start net.IP
		end   net.IP
	}{
		{net.ParseIP("10.0.0.0"), net.ParseIP("10.255.255.255")},
		{net.ParseIP("172.16.0.0"), net.ParseIP("172.31.255.255")},
		{net.ParseIP("192.168.0.0"), net.ParseIP("192.168.255.255")},
	}

	for _, block := range privateIPBlocks {
		if bytes.Compare(ip, block.start) >= 0 && bytes.Compare(ip, block.end) <= 0 {
			return true
		}
	}
	return false
}
