// Package log package log
package log

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"trade_agent/global"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	globalLogger *logrus.Logger
	once         sync.Once
)

func initLogger() {
	if globalLogger != nil {
		return
	}
	// Get current path
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	basePath := filepath.Dir(ex)
	// create new instance
	globalLogger = logrus.New()
	deployment := os.Getenv(global.EnvDeployment)
	if deployment == global.DockerDeployment {
		globalLogger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: global.LongTimeLayout,
			PrettyPrint:     false,
		})
	} else {
		globalLogger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat:  "2006/01/02 15:04:05",
			FullTimestamp:    true,
			QuoteEmptyFields: true,
			PadLevelText:     false,
			ForceColors:      true,
			ForceQuote:       true,
		})
	}
	// Log.SetReportCaller(true)
	folderName := time.Now().Format(time.RFC3339)[:16]
	folderName = strings.ReplaceAll(folderName, ":", "")
	globalLogger.SetLevel(logrus.TraceLevel)
	globalLogger.SetOutput(os.Stdout)
	pathMap := lfshook.PathMap{
		logrus.PanicLevel: basePath + "/logs/" + folderName + "/panic.json",
		logrus.FatalLevel: basePath + "/logs/" + folderName + "/fetal.json",
		logrus.ErrorLevel: basePath + "/logs/" + folderName + "/error.json",
		logrus.WarnLevel:  basePath + "/logs/" + folderName + "/warn.json",
		logrus.InfoLevel:  basePath + "/logs/" + folderName + "/info.json",
		logrus.DebugLevel: basePath + "/logs/" + folderName + "/debug.json",
		logrus.TraceLevel: basePath + "/logs/" + folderName + "/error.json",
	}
	globalLogger.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.JSONFormatter{},
	))
}

// Get Get
func Get() *logrus.Logger {
	if globalLogger != nil {
		return globalLogger
	}
	once.Do(initLogger)
	return globalLogger
}
