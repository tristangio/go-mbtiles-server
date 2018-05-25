package main

import (
	"log"
	"net"
	"strings"
)

// GetOutboundIP Get preferred outbound ip of this machine
func GetOutboundIP() string {
	conn, err := net.Dial("udp", "9.9.9.9:80")
	if err != nil {
		log.Printf("GetOutboundIP ERR: %s\n", err)
		return ""
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")

	return localAddr[0:idx]
}
