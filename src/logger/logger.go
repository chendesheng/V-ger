package logger

import (
	"io"
	"log"
	"log/syslog"
	"os"
	"syscall"
	// "util"
)

func InitLog(prefix string, crashLogFile string) {
	log.SetFlags(log.Lshortfile)

	l, _ := syslog.New(syslog.LOG_NOTICE, prefix)
	w := io.MultiWriter(os.Stdout, l)

	log.SetOutput(w)

	if len(crashLogFile) > 0 {
		crashLog(crashLogFile)
	}
}

func crashLog(file string) {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Print(err.Error())
	} else {
		defer f.Close()

		// syscall.Dup2(int(f.Fd()), 1)
		syscall.Dup2(int(f.Fd()), 2)
	}
}
