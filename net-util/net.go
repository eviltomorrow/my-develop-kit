package netutil

import (
	"net"
	"strings"
	"time"
)

// GetExternalIP get external ip
func GetExternalIP() (string, error) {
	data, err := GetHTTP("https://myexternalip.com/raw", 10*time.Second, nil)
	if err != nil {
		return "", err
	}
	return data, nil
}

// IsPublicIP is public ip
func IsPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

// GetLocalIPv4s get local ipv4
func GetLocalIPv4s() ([]string, error) {
	address, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	var ips = make([]string, 0, len(address))
	for _, addr := range address {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			ips = append(ips, ipNet.IP.String())
		}
	}
	return ips, nil
}

// GetLocalIP get local ip
func GetLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx], nil
}
