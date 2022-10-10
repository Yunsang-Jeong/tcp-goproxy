package proxy

import (
	"io"
	"log"
	"net"

	"github.com/fatih/color"
)

func StartTCPProxy(cConn *net.TCPConn, targetAddr string) {
	addr, err := net.ResolveTCPAddr("tcp", targetAddr)
	if err != nil {
		log.Println(color.HiRedString("failed to resolve address: %w", err))
		return
	}

	tConn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Println(color.HiRedString("failed to dail target server: %w", err))
		return
	}
	defer tConn.Close()

	log.Println(color.HiCyanString("Proxy the tcp packet to target(%s --> %s)", tConn.LocalAddr().String(), tConn.RemoteAddr().String()))

	go func() {
		if _, err := io.Copy(tConn, cConn); err != nil {
			log.Println(color.HiRedString("failed to send payload (server --> target): %w", err))
		} else {
			log.Println(color.HiGreenString("success to send payload (server --> target)"))
		}
	}()

	if _, err := io.Copy(cConn, tConn); err != nil {
		log.Println(color.HiRedString("failed to send payload (target --> server): %w", err))
	} else {
		log.Println(color.HiGreenString("success to send payload (target --> server)"))
	}
}
