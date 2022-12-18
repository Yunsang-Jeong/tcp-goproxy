package handler

import (
	"fmt"
	"io"
	"net"

	"github.com/Yunsang-Jeong/tcp-goproxy/internal/logger"
)

type tcpProxyHandler struct {
	cConn      *net.TCPConn
	targetAddr string
	logging    logger.LoggingFunc
}

func NewTCPProxyHandler(cConn *net.TCPConn, targetAddr string, loggingFunc logger.LoggingFunc) tcpProxyHandler {
	return tcpProxyHandler{
		cConn:      cConn,
		targetAddr: targetAddr,
		logging:    loggingFunc,
	}
}

func (h *tcpProxyHandler) Handle() {
	addr, err := net.ResolveTCPAddr("tcp", h.targetAddr)
	if err != nil {
		h.logging(fmt.Sprintf("failed to resolve address: %s", err))
		return
	}

	tConn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		h.logging(fmt.Sprintf("failed to dail target server: %s", err))
		return
	}
	defer tConn.Close()

	h.logging(fmt.Sprintf("Proxy the tcp packet to target(%s --> %s)", tConn.LocalAddr().String(), tConn.RemoteAddr().String()))

	go func() {
		if _, err := io.Copy(tConn, h.cConn); err != nil {
			h.logging(fmt.Sprintf("failed to send payload (server --> target): %s", err))
		} else {
			h.logging("success to send payload (server --> target)")
		}
	}()

	if _, err := io.Copy(h.cConn, tConn); err != nil {
		h.logging(fmt.Sprintf("failed to send payload (target --> server): %s", err))
	} else {
		h.logging("success to send payload (target --> server)")
	}
}
