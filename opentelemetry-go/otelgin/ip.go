package otelgin

import "net"

var LocalIp string

func init() {
	if ip := GetLocalIp(); ip != nil {
		LocalIp = *ip
	}
}

func GetLocalIp() *string {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	for i := 0; i < len(netInterfaces); i++ {
		if (netInterfaces[i].Flags & net.FlagUp) != 0 {
			addresses, _ := netInterfaces[i].Addrs()

			for _, address := range addresses {
				if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						var ip string
						ip = ipNet.IP.String()
						return &ip
					}
				}
			}
		}
	}
	return nil
}
