package app

import (
	"os"

	"github.com/google/logger"
)

func Init_Logger(logPath string) {
	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}
	// defer lf.Close()

	// defer logger.Init("LoggerExample", *verbose, true, lf).Close()
	logger.Init("ETH-Transfer", true, true, lf)

	logger.Info("Successfully Run")
}
