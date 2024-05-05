package util

import "net"

type NetworkUtil struct {
	networks []*net.IPNet
}

func NewNetworkUtil() *NetworkUtil {
	nu := new(NetworkUtil)

	networks := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"127.0.0.0/8",
		"::1/128",
	}

	for _, v := range networks {
		_, ipnet, err := net.ParseCIDR(v)
		if err != nil {
			continue
		}

		nu.networks = append(nu.networks, ipnet)
	}

	return nu
}

// IsPrivateAddress はプライベートアドレスかどうかを調べる
func (nu NetworkUtil) IsPrivateAddress(ipString string) bool {
	ip := net.ParseIP(ipString)

	if ip == nil {
		return false
	}

	for _, v := range nu.networks {
		if v.Contains(ip) {
			return true
		}
	}

	return false
}
