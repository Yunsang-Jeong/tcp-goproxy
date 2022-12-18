package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Yunsang-Jeong/tcp-goproxy/internal/logger"
	"github.com/jlaffaye/ftp"
)

type readOnlyFTPProxyHandler struct {
	newLoggingFunc logger.NewLoggingFunc
	user           string
	password       string
	resource       string
	conn           *ftp.ServerConn
}

type ftpLISTResponse struct {
	fileName string
	fileType string
	fileSize uint64
}

const timeOut = 5

const (
	errInvalidMethod       = "failed to process this method"
	errParseRequestURI     = "failed to parse request"
	errFTPConnectWithLogin = "failed to connect(login) ftp server"
	errFTPLIST             = "failed to run LIST command"
	errFTPRETR             = "failed to run RETR command"
	errFTPRETRRead         = "failed to read result of RETR command"
)

func NewReadOnlyFTPProxyHandler(newLoggingFunc logger.NewLoggingFunc) readOnlyFTPProxyHandler {
	return readOnlyFTPProxyHandler{
		newLoggingFunc: newLoggingFunc,
	}
}

func (h *readOnlyFTPProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logging := h.newLoggingFunc()

	logging(fmt.Sprintf("New  :: From %s", r.RemoteAddr))
	defer logging(fmt.Sprintf("Close:: From %s\n", r.RemoteAddr))

	switch r.Method {
	case "GET":
		if err := h.parseFTPRequestURI(r.RequestURI); err != nil {
			logging(fmt.Sprintf("%s: %s", errParseRequestURI, err.Error()), w)
			return
		}

		if err := h.ftpConnectWithLogin(r.Host); err != nil {
			logging(fmt.Sprintf("%s: %s", errFTPConnectWithLogin, err.Error()), w)
			return
		}
		defer h.conn.Quit()

		if strings.HasSuffix(h.resource, "/") {
			//
			// FTP LIST
			//
			logging(fmt.Sprintf("Try to send LIST(%s) Command to %s", h.resource, r.Host))
			resp, err := h.sendFTPLISTCommand()
			if err != nil {
				logging(fmt.Sprintf("%s: %s", errFTPLIST, err.Error()), w)
			}

			for _, r := range resp {
				logging(fmt.Sprintf("[%s] %s (%d bytes)", r.fileType, r.fileName, r.fileSize), w)
			}
			logging(fmt.Sprintf("Sucess to listing files on %s", h.resource))
		} else {
			//
			// FTP RETR (GET)
			//
			logging(fmt.Sprintf("Try to send RETR(%s) Command to %s", h.resource, r.Host))
			resp, err := h.sendFTPRETRCommand()
			if err != nil {
				logging(fmt.Sprintf("%s: %s", errFTPRETR, err.Error()), w)
				return
			}

			if size, err := w.Write(resp); err != nil {
				logging(fmt.Sprintf("%s: %s", errFTPRETRRead, err.Error()), w)
			} else {
				logging(fmt.Sprintf("Sucess to send data to client: (%d) bytes", size))
			}
		}

	default:
		logging(fmt.Sprintf("%s: This is read only ftp proxy server", errInvalidMethod), w)
	}
}

func (h *readOnlyFTPProxyHandler) sendFTPLISTCommand() ([]ftpLISTResponse, error) {
	result := []ftpLISTResponse{}
	resp, err := h.conn.List(h.resource)
	if err != nil {
		return nil, err
	}

	for _, r := range resp {
		result = append(result, ftpLISTResponse{
			fileName: r.Name,
			fileType: r.Type.String(),
			fileSize: r.Size,
		})
	}

	return result, nil
}

func (h *readOnlyFTPProxyHandler) sendFTPRETRCommand() ([]byte, error) {
	resp, err := h.conn.Retr(h.resource)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	result, err := io.ReadAll(resp)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (h *readOnlyFTPProxyHandler) ftpConnectWithLogin(host string) error {
	_, _, flag := strings.Cut(host, ":")
	if !flag {
		host = fmt.Sprintf("%s:21", host)
	}

	conn, err := ftp.Dial(host, ftp.DialWithTimeout(timeOut*time.Second))
	if err != nil {
		return err
	}
	h.conn = conn

	err = conn.Login(h.user, h.password)
	if err != nil {
		h.conn.Quit()
		return err
	}

	return nil
}

func (h *readOnlyFTPProxyHandler) parseFTPRequestURI(uri string) error {
	var flag bool
	var resource string

	proto, part1, flag := strings.Cut(uri, "://")
	if !flag || proto != "ftp" {
		return fmt.Errorf("invalid proto in request uri")
	}

	h.user = "anonymous"
	auth, part2, flag := strings.Cut(part1, "@")
	if flag {
		user, password, flag := strings.Cut(auth, ":")
		if flag {
			h.user = user
			h.password = password
		}
	}

	if flag {
		_, resource, flag = strings.Cut(part2, "/")
	} else {
		_, resource, flag = strings.Cut(part1, "/")
	}
	if !flag {
		return fmt.Errorf("failed to parse resource: %s", uri)
	}
	h.resource = fmt.Sprintf("/%s", resource)

	return nil
}
