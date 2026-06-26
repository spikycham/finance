package logger

import (
	"log"
	"os"
)

var (
	errorLog = log.New(os.Stderr, "[ERROR]", log.LstdFlags|log.Lshortfile)
	warnLog  = log.New(os.Stderr, "[WARN]", log.LstdFlags)
	infoLog  = log.New(os.Stdout, "[INFO]", log.LstdFlags)
)

func Error(msg string, err error) {
	errorLog.Printf("%s: %v", msg, err)
}

func Warn(msg string) {
	warnLog.Println(msg)
}

func Info(msg string) {
	infoLog.Println(msg)
}

func Infof(format string, v ...any) {
	infoLog.Printf(format, v...)
}
