package utils

import (
	"log"
	"net"
)

func GetOpenPort() int {
	// 1) Ask the OS for a free TCP port on all interfaces.
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer ln.Close()

	// 2) Extract the chosen port.
	addr := ln.Addr().(*net.TCPAddr)
	port := addr.Port
	return port
}
