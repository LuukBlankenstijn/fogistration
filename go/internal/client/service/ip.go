package service

import "net"

func GetIp() (net.IP, error) {
	conn, err := net.Dial("udp", "1.1.1.1:80") // use random ip, does not matter
	if err != nil {
		return nil, err
	}
	defer func() { _ = conn.Close() }()

	addr := conn.LocalAddr().(*net.UDPAddr)

	return addr.IP, nil
}
