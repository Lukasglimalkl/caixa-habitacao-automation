package logger

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

// Inicializa os loggers
func Init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Info loga mensagens de informação
func Info(message string) {
	InfoLogger.Println(message)
}

// Error loga mensagens de erro
func Error(message string) {
	ErrorLogger.Println(message)
}

