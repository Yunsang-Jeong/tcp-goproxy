package logger

import (
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/fatih/color"
	"github.com/google/uuid"
)

type LoggingFunc func(message string, writers ...io.Writer)
type NewLoggingFunc func() LoggingFunc
type colorPrinter func(format string, a ...interface{}) string

type loggerManager struct {
	colorPrinterList []colorPrinter
	nextColorPrinter int
	m                sync.Mutex
}

type logger struct {
	identifier   string
	colorPrinter colorPrinter
}

func NewLoggerManager() NewLoggingFunc {
	colorPrinterList := []colorPrinter{
		color.HiWhiteString,
		color.HiYellowString,
		color.HiCyanString,
		color.HiGreenString,
		color.HiBlueString,
		color.HiMagentaString,
	}

	lm := loggerManager{
		colorPrinterList: colorPrinterList,
	}

	return lm.newLoggingFunc
}

func (lm *loggerManager) newLoggingFunc() LoggingFunc {
	uuid := uuid.New()

	lm.m.Lock()
	colorPrinter := lm.colorPrinterList[lm.nextColorPrinter]
	lm.nextColorPrinter = (lm.nextColorPrinter + 1) % len(lm.colorPrinterList)
	lm.m.Unlock()

	l := logger{
		identifier:   uuid.String(),
		colorPrinter: colorPrinter,
	}
	return l.logging
}

func (l *logger) logging(message string, writers ...io.Writer) {
	msg := fmt.Sprintf("[%s] %s", l.identifier, message)
	msgWithNewline := fmt.Sprintln(message)

	log.Print(l.colorPrinter(msg))

	for _, w := range writers {
		w.Write([]byte(msgWithNewline))
	}
}
