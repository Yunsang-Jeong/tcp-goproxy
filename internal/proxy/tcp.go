package proxy

import (
	"fmt"
	"net"

	"github.com/Yunsang-Jeong/tcp-goproxy/internal/handler"
	"github.com/Yunsang-Jeong/tcp-goproxy/internal/logger"

	"github.com/fatih/color"
)

type tcpProxyServer struct {
	serverAddr string
	targetAddr string
}

func NewTCPProxyServer(serverAddr string, targetAddr string) tcpProxyServer {
	return tcpProxyServer{
		serverAddr: serverAddr,
		targetAddr: targetAddr,
	}
}

func (ps *tcpProxyServer) Start() error {
	newLoggingFunc := logger.NewLoggerManager()
	psLogging := newLoggingFunc()

	laddr, err := net.ResolveTCPAddr("tcp", ps.serverAddr)
	if err != nil {
		psLogging(fmt.Sprintf("failed to resolve service address: %s", err))
		return err
	}

	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		psLogging(fmt.Sprintf("failed to open service port: %s", err))
		return err
	}
	defer listener.Close()

	psLogging(fmt.Sprintf("Start to tcp proxy server: %s", ps.serverAddr))

	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			color.HiRedString("failed to resolve service address: %w", err)
			return err
		}
		conn.Close()

		hLogging := newLoggingFunc()
		hLogging(fmt.Sprintf("New  :: From %s", conn.LocalAddr().String()))

		handler := handler.NewTCPProxyHandler(conn, ps.targetAddr, hLogging)
		go handler.Handle()
	}
}
