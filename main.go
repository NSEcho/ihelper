package main

import (
	"github.com/charmbracelet/log"
	"github.com/lateralusd/ihelper/cmd"
)

func main() {
	logger := log.New()
	logger.SetReportTimestamp(false)
	if err := cmd.Execute(logger); err != nil {
		logger.Error(err)
	}
}
