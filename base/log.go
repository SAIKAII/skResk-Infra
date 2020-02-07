package base

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	// 定义日志格式
	formatter := &log.TextFormatter{}
	formatter.FullTimestamp = true
	formatter.TimestampFormat = "2006-01-02.15:04:05.000000"
	log.SetFormatter(formatter)
	// 日志级别
	level := os.Getenv("log.debug")
	if level == "true" {
		log.SetLevel(log.DebugLevel)
	}
	// 控制台高亮显示
	formatter.ForceColors = true
	formatter.DisableColors = false
}
