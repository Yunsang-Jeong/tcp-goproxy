package proxy

import (
	"net/http"

	"github.com/Yunsang-Jeong/tcp-goproxy/internal/handler"
	"github.com/Yunsang-Jeong/tcp-goproxy/internal/logger"
)

type readOnlyFTPProxyServer struct {
	serverAddr string
}

func NewFTPProxyServer(serverAddr string) readOnlyFTPProxyServer {
	return readOnlyFTPProxyServer{
		serverAddr: serverAddr,
	}
}

func (ps *readOnlyFTPProxyServer) Start() error {
	newLoggingFunc := logger.NewLoggerManager()

	logging := newLoggingFunc()
	logging("Start to read-only ftp proxy server")

	handler := handler.NewReadOnlyFTPProxyHandler(newLoggingFunc)
	if err := http.ListenAndServe(ps.serverAddr, &handler); err != nil {
		return err
	}

	return nil
}
