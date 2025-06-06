// utils/logger.go
package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogConfig struct {
	ToConsole bool
	ToFile    bool
	FilePath  string // 目录，如 logs/
	IsDaily   bool   // 是否按日期生成子文件
}

type Logger struct {
	logger *log.Logger
}

var (
	loggerRegistry = make(map[string]*Logger)
	loggerConfigs  = make(map[string]LogConfig)
	loggerMutex    sync.Mutex
)

func InitLogConfigs(configs map[string]LogConfig) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	loggerConfigs = configs
}

func Log(name string) *Logger {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()

	if logger, exists := loggerRegistry[name]; exists {
		return logger
	}

	cfg, ok := loggerConfigs[name]
	if !ok {
		cfg = LogConfig{ToConsole: true} // fallback 默认配置
	}

	var writers []io.Writer

	if cfg.ToConsole {
		writers = append(writers, os.Stdout)
	}

	if cfg.ToFile {
		logDir := filepath.Join(cfg.FilePath, name)
		_ = os.MkdirAll(logDir, os.ModePerm)

		var filename string
		if cfg.IsDaily {
			// logs/app/20250606.log 或 logs/app/2025/06/06.log
			date := time.Now().Format(time.DateOnly)
			filename = filepath.Join(logDir, fmt.Sprintf("%s.log", date))
		} else {
			filename = filepath.Join(logDir, fmt.Sprintf("%s.log", name))
		}

		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			writers = append(writers, file)
		} else {
			fmt.Println("日志文件打开失败:", err)
		}
	}

	multiWriter := io.MultiWriter(writers...)
	logger := log.New(multiWriter, "", log.LstdFlags|log.Lshortfile)

	wrap := &Logger{logger: logger}
	loggerRegistry[name] = wrap
	return wrap
}

func (l *Logger) Info(v ...interface{}) {
	l.logger.SetPrefix("[INFO] ")
	l.logger.Println(v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.logger.SetPrefix("[WARN] ")
	l.logger.Println(v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.logger.SetPrefix("[ERROR] ")
	l.logger.Println(v...)
}
