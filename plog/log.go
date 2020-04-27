package plog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// Heartf monitor logging
func Heartf(format string, args ...interface{}) {
	loggerHeart.Infof(format, args...)
}

// Infof info logging
func Infof(format string, args ...interface{}) {
	loggerNormal.Infof(format, args...)
}

// Debugf debug logging
func Debugf(format string, args ...interface{}) {
	loggerNormal.Debugf(format, args...)
}

// Errorf info logging
func Errorf(format string, args ...interface{}) {
	loggerNormal.Errorf(format, args...)
}

// Fatalf info logging
func Fatalf(format string, args ...interface{}) {
	loggerNormal.Fatalf(format, args...)
}

// Warningf info logging
func Warningf(format string, args ...interface{}) {
	loggerNormal.Warningf(format, args...)
}

// logTypeToColor converts the Level to a color string.
func logTypeToColor(level logrus.Level) string {
	switch level {
	case logrus.DebugLevel:
		return "[0;37"
	case logrus.InfoLevel:
		return "[0;36"
	case logrus.WarnLevel:
		return "[0;33"
	case logrus.ErrorLevel:
		return "[0;31"
	case logrus.FatalLevel:
		return "[0;31"
	case logrus.PanicLevel:
		return "[0;31"
	}

	return "[0;37"
}

// textFormatter is for compatibility with ngaut/log
type textFormatter struct {
	DisableTimestamp bool
	EnableColors     bool
	EnableEntryOrder bool
}

// Format implements logrus.Formatter
func (f *textFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	if f.EnableColors {
		colorStr := logTypeToColor(entry.Level)
		fmt.Fprintf(b, "\033%sm ", colorStr)
	}

	if !f.DisableTimestamp {
		fmt.Fprintf(b, "%s ", entry.Time.Format(defaultLogTimeFormat))
	}
	if file, ok := entry.Data["file"]; ok {
		fmt.Fprintf(b, "%s:%v:", file, entry.Data["line"])
	}
	fmt.Fprintf(b, " [%s] %s", entry.Level.String(), entry.Message)

	if f.EnableEntryOrder {
		keys := make([]string, 0, len(entry.Data))
		for k := range entry.Data {
			if k != "file" && k != "line" {
				keys = append(keys, k)
			}
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(b, " %v=%v", k, entry.Data[k])
		}
	} else {
		for k, v := range entry.Data {
			if k != "file" && k != "line" {
				fmt.Fprintf(b, " %v=%v", k, v)
			}
		}
	}

	b.WriteByte('\n')

	if f.EnableColors {
		b.WriteString("\033[0m")
	}
	return b.Bytes(), nil
}

func stringToLogLevel(level string) logrus.Level {
	switch strings.ToLower(level) {
	case "fatal":
		return logrus.FatalLevel
	case "error":
		return logrus.ErrorLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	}
	return defaultLogLevel
}

func stringToLogFormatter(format string) logrus.Formatter {
	switch strings.ToLower(format) {
	case "text":
		return &textFormatter{
			// TimestampFormat:  defaultLogTimeFormat,
			DisableTimestamp: false,
		}
	case "json":
		return &logrus.JSONFormatter{
			TimestampFormat:  defaultLogTimeFormat,
			DisableTimestamp: false,
		}
	case "console":
		return &logrus.TextFormatter{
			FullTimestamp:    true,
			TimestampFormat:  defaultLogTimeFormat,
			DisableTimestamp: false,
		}
	default:
		return &textFormatter{}
	}
}

// initFileLog initializes file based logging options.
func addHook(path string, format logrus.Formatter, level ...logrus.Level) *lfshook.LfsHook {
	var writer *rotatelogs.RotateLogs
	var err error

	switch runtime.GOOS {
	case "windows":
		writer, err = rotatelogs.New(
			path+".%Y%m%d%H",
			// rotatelogs.WithLinkName(path),          // 生成软链，指向最新日志文件
			rotatelogs.WithMaxAge(24*90*time.Hour),    // 文件最大保存时间
			rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
		)
		if err != nil {
			panic(err)
		}

	case "linux":
		writer, err = rotatelogs.New(
			path+".%Y%m%d%H",
			// rotatelogs.WithLinkName(path),             // 生成软链，指向最新日志文件
			rotatelogs.WithMaxAge(24*90*time.Hour),    // 文件最大保存时间
			rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
		)
		if err != nil {
			panic(err)
		}

	default:
		writer, err = rotatelogs.New(
			path+".%Y%m%d%H",
			rotatelogs.WithMaxAge(24*90*time.Hour),    // 文件最大保存时间
			rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
		)
		if err != nil {
			panic(err)
		}
	}

	option := lfshook.WriterMap{}
	for _, l := range level {
		option[l] = writer
	}
	return lfshook.NewHook(option, format)

}

var (
	defaultLogTimeFormat = "2006-01-02 15:04:05"
	defaultLogLevel      = logrus.InfoLevel

	loggerHeart  = logrus.New()
	loggerNormal = logrus.New()
)

// LogConfig serializes log related config in toml/json.
type LogConfig struct {
	// Log level.
	Level string
	// Log format. one of json, text, or console.
	Format string
	// LogFilePath log file path
	LogFileDir string
}

func createDIR(path string) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		err = os.MkdirAll(path, 0766)
		if err != nil {
			panic(err)
		}
		return
	}
	if !fileInfo.IsDir() {
		err = os.MkdirAll(path, 0766)
		if err != nil {
			panic(err)
		}
	}
}

// InitLogger initialize logger
func InitLogger(config *LogConfig) {
	createDIR(config.LogFileDir)

	format := stringToLogFormatter(config.Format)

	// Normal logger
	loggerNormal.SetLevel(stringToLogLevel(config.Level))
	loggerNormal.SetFormatter(format)
	loggerNormal.AddHook(addHook(
		path.Join(config.LogFileDir, "error.log"),
		format,
		logrus.ErrorLevel, logrus.FatalLevel,
	))
	loggerNormal.AddHook(addHook(
		path.Join(config.LogFileDir, "data.log"),
		format,
		stringToLogLevel(config.Level), logrus.InfoLevel, logrus.ErrorLevel,
	))

	// Heart logger
	loggerHeart.SetLevel(logrus.InfoLevel)
	loggerHeart.SetFormatter(format)
	loggerHeart.AddHook(addHook(
		path.Join(config.LogFileDir, "heart.log"),
		format,
		logrus.InfoLevel,
	))

}

type jsonFormatter struct {
}

func (j *jsonFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Note this doesn't include Time, Level and Message which are available on
	// the Entry. Consult `godoc` on information about those fields or read the
	// source of the official loggers.
	serialized, err := json.Marshal(entry.Data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
