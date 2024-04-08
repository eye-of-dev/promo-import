package logger

import (
	"fmt"
	"log"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

type FileSLogger struct {
	logDir         string
	defaultLogFile string
	env            string
	Logger         *slog.Logger
}

func NewFileSLogger(logDir, defaultLogFile, env string) *FileSLogger {
	return &FileSLogger{logDir, defaultLogFile, env, nil}
}

func (l *FileSLogger) SetUpLogger() *FileSLogger {
	switch l.env {
	case envLocal:
		l.Logger = slog.New(
			slog.NewTextHandler(l.openFileTarget(), &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		l.Logger = slog.New(
			slog.NewTextHandler(l.openFileTarget(), &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		l.Logger = slog.New(
			slog.NewTextHandler(l.openFileTarget(), &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return l
}

func (l *FileSLogger) GetLevel() slog.Level {
	var level slog.Level
	switch l.env {
	case envLocal:
		level = slog.LevelDebug
	case envDev:
		level = slog.LevelDebug
	case envProd:
		level = slog.LevelInfo
	}

	return level
}

func (l *FileSLogger) openFileTarget() *os.File {
	f, err := os.OpenFile(
		fmt.Sprintf("%s/%s.txt", l.logDir, l.defaultLogFile),
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0666,
	)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	return f
}

func (l *FileSLogger) closeFileTarget(f *os.File) {
	err := f.Close()
	if err != nil {
		log.Fatalf("error closing file: %v", err)
	}
}
