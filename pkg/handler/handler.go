package handler

type Handler interface {
	Start() error
	startListener(string) error
}
