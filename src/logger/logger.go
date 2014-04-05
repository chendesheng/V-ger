package logger

import (
	"io"
	"log"
	"log/syslog"
	"os"
	// "util"
)

// func init() {
// 	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
// 	w := logWriter{}
// 	w.writers = append(w.writers, os.Stdout)

// 	if logPath := util.ReadConfig("log"); logPath != "" {
// 		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
// 		if err == nil {
// 			os.Stderr = f
// 			w.writers = append(w.writers, f)
// 		} else {
// 			log.Print(err)
// 		}
// 	}

// 	log.SetOutput(w)
// 	log.Print("log initialized.")
// }

type logWriter struct {
	writers []io.Writer
}

func (l logWriter) Write(p []byte) (int, error) {
	for _, w := range l.writers {
		w.Write(p)
	}

	return len(p), nil
}

func InitLog(prefix string) {
	log.SetFlags(log.Lshortfile)
	w := logWriter{}
	w.writers = append(w.writers, os.Stdout)

	l, _ := syslog.New(syslog.LOG_NOTICE, prefix)
	w.writers = append(w.writers, l)

	// if filename != "" {
	// 	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	// 	if err == nil {
	// 		// os.Stderr = f
	// 		w.writers = append(w.writers, f)
	// 	} else {
	// 		log.Print(err)
	// 	}
	// }

	log.SetOutput(w)
}
