package common

import (
	"net"
	"strings"
)

// TODO this requires a lot of cleaning up, and need to fix the
// windows version of this.
//
// naive way to get your local network IP (eg: 192.169.1.10)
func GetLocalIP() ([]net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var ips []net.IP
	for _, el := range ifaces {
		addrs, err := el.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			if !(strings.Contains(el.Name, "enp") ||
				strings.Contains(el.Name, "wlp")) {
				continue
			}

			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.To4() == nil {
					continue
				}
				ips = append(ips, v.IP)
			case *net.IPAddr:
				if v.IP.To4() == nil {
					continue
				}
				ips = append(ips, v.IP)
			}

		}
	}

	return ips, nil
}
