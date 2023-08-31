package app

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

func SetupLogrus(level, logsPath string) error {
	loggerLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(loggerLevel)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetOutput(io.Discard)

	if err = createLogsDir(logsPath); err != nil {
		return fmt.Errorf("failed to detect logs directory by path \"%s\" due to error: %w", logsPath, err)
	}

	logsFile, err := os.OpenFile(logsPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		return fmt.Errorf("failed to create or open logs file by path \"%s\" due to error: %w", logsPath, err)
	}

	logrus.AddHook(&writerHook{
		Writer:    []io.Writer{logsFile, os.Stdout},
		LogLevels: logrus.AllLevels,
	})

	return nil
}

// createLogsDir создаёт директорию, в которой будет храниться
// текстовый файл с логами, в случае её отсутствия.
func createLogsDir(path string) error {
	path = filepath.Dir(path)
	// Создадим папку для хранения логов.
	exists, err := func(path string) (bool, error) {
		_, err := os.Stat(path)
		if err == nil {
			return true, nil
		}
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}(path)
	if err != nil {
		return err
	}
	if !exists {
		if err := os.Mkdir(path, 0666); err != nil {
			panic(err)
		}
	}
	return nil
}

// Хук на запись для модификации логирования logrus.
type writerHook struct {
	Writer    []io.Writer
	LogLevels []logrus.Level
}

func (hook *writerHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range hook.Writer {
		_, err = w.Write([]byte(line))
		if err != nil {
			break
		}
	}
	return err
}

func (hook *writerHook) Levels() []logrus.Level {
	return hook.LogLevels
}
