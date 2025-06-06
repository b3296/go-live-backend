// utils/logger.go
package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

type LogConfig struct {
	ToConsole bool
	ToFile    bool
	FilePath  string // 目录，如 logs/
	IsDaily   bool   // 是否按日期生成子文件
	AsJSON    bool   // 是否以 JSON 格式输出
	Level     LogLevel
}

type Logger struct {
	name   string
	logger *log.Logger
	config LogConfig
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
		cfg = LogConfig{ToConsole: true, Level: INFO} // fallback 默认配置
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
	goLogger := log.New(multiWriter, "", log.LstdFlags|log.Lshortfile)

	wrap := &Logger{logger: goLogger, config: cfg, name: name}
	loggerRegistry[name] = wrap
	return wrap
}

func (l *Logger) logf(level LogLevel, format string, v ...interface{}) {
	if level < l.config.Level {
		return
	}

	if l.config.AsJSON {
		entry := map[string]interface{}{
			"level":   level.String(),
			"name":    l.name,
			"time":    time.Now().Format(time.RFC3339),
			"message": fmt.Sprintf(format, v...),
		}
		jsonData, _ := json.Marshal(entry)
		l.logger.Println(string(jsonData))
	} else {
		prefix := fmt.Sprintf("[%s] ", level.String())
		l.logger.SetPrefix(prefix)
		l.logger.Output(3, fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logf(DEBUG, format, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.logf(INFO, format, v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.logf(WARN, format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logf(ERROR, format, v...)
}
