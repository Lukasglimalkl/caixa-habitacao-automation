package logger

import (
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
)

// Init - inicializa o logger
func Init() {
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Info - loga mensagem de informação
func Info(message string) {
	if infoLogger == nil {
		Init() // Auto-inicializa se não foi inicializado
	}
	infoLogger.Println(message)
}

// Error - loga mensagem de erro
func Error(message string) {
	if errorLogger == nil {
		Init() // Auto-inicializa se não foi inicializado
	}
	errorLogger.Println(message)
}