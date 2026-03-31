package internal

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	log *logrus.Logger
}

func NewLogger() *Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.DebugLevel)
	log.SetOutput(os.Stdout)

	if !envBoolDefault("LOG_TO_FILE", true) {
		return &Logger{log: log}
	}

	logDir := envOrDefault("LOG_DIR", "log")
	if err := os.MkdirAll(logDir, 0o755); err != nil {
		log.WithError(err).Warn("failed to create log directory, using stdout only")
		return &Logger{log: log}
	}

	filename := filepath.Join(logDir, time.Now().UTC().Format("2006-01-02")+".log")
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
	if err != nil {
		log.WithError(err).Warn("failed to open log file, using stdout only")
		return &Logger{log: log}
	}

	log.SetOutput(io.MultiWriter(os.Stdout, file))

	return &Logger{log: log}
}

func (l *Logger) GetLogger() *logrus.Logger {
	return l.log
}

func envOrDefault(key string, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}

	return fallback
}

func envBoolDefault(key string, fallback bool) bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	switch value {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	case "":
		return fallback
	default:
		return fallback
	}
}
