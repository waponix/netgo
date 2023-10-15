package main

import (
	"github.com/waponix/netgo/logger"
)

func main() {
	appLog := logger.New()

	appLog.Filename = "dev.log"
	appLog.LogLevels = []string{logger.INFO, logger.ERROR}

	appLog.Info("App was initialized")
}
