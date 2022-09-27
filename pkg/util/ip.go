package util

import (
	"net"
	"strings"
)

var outBoundIP string

func GetOutBoundIP() (string, error) {
	if outBoundIP != "" {
		return outBoundIP, nil
	}
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return "", err
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	outBoundIP = strings.Split(localAddr.String(), ":")[0]
	return outBoundIP, nil
}
