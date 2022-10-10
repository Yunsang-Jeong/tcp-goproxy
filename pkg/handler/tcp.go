package handler

import (
	"log"
	"net"

	"github.com/Yunsang-Jeong/tcp-goproxy/pkg/proxy"

	"github.com/fatih/color"
)

func StartTCPHandler(serverAddr string) error {
	laddr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		color.HiRedString("failed to resolve service address: %w", err)
		return err
	}

	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		color.HiRedString("failed to open service port: %w", err)
		return err
	}
	defer listener.Close()

	log.Printf("Listening on %s ..\n", serverAddr)

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			color.HiRedString("failed to resolve service address: %w", err)
			return err
		}

		log.Println(color.HiCyanString("Accpeted:: %s --> %s", conn.LocalAddr().String(), conn.RemoteAddr().String()))
		go proxy.StartTCPProxy(conn, "127.0.0.1:9090")
	}
}
