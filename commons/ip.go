package commons

import (
	"net"
	"strconv"
)

// IpPort 支持 ipv6:port, [ipv6]:port, ip:port
func IpPort(ip string, port int) string {
	if IsIPV6(ip) {
		return "[" + ip + "]:" + strconv.Itoa(port)
	}
	return ip + ":" + strconv.Itoa(port)
}

// IsIPV6 支持ipv6, 不支持 [ipv6]
func IsIPV6(s string) bool {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return false
		case ':':
			return net.ParseIP(s) != nil
		}
	}
	return false
}
