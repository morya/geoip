package iputil

import "net"

// IsIntranet 判断是否局域网IP
func IsIntranet(ip net.IP) bool {
	if ip.IsLoopback() {
		return true
	}

	ipV4 := ip.To4()
	switch true {
	case ipV4[0] == 10:
		return true
	case ipV4[0] == 172 && ipV4[1] >= 16 && ipV4[1] <= 31:
		return true
	case ipV4[0] == 192 && ipV4[1] == 168:
		return true
	default:
		break
	}
	return false
}
