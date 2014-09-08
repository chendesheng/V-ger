package logger

import (
	"log"
	"os"
	"syscall"
)

func InitLog(prefix string, file string) {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Print(err.Error())
		return
	}
	syscall.Dup2(int(f.Fd()), 2)

	log.SetFlags(log.Lshortfile | log.Ltime)
	log.SetPrefix(prefix)

	log.SetOutput(f)

	println("\n\n==============================================================================")
}

func crashLog(file string) {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Print(err.Error())
	} else {
		// defer f.Close()

		// syscall.Dup2(int(f.Fd()), 1)
		syscall.Dup2(int(f.Fd()), 2)
	}
}
