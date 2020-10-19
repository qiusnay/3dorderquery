package main

import (
	"os"
	"github.com/google/logger"
	"github.com/qiusnay/3dorderquery/service"
)

func main() {
	const logPath = "./log/3dorderquery.log"
	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
	  logger.Fatalf("Failed to open log file: %v", err)
	}
	logger.Init("Logger", false, true, lf)
	service := new(service.Jdsdk)
	service.GetOrders()
}