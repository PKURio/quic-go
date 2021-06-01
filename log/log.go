/*
Package log 提供基本的日志，分级，输出，轮转等功能
- 考虑使用公司的Log
*/
package log

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"sync"
)

const (
	maxLogSize = 512 // 单个日志文件最大存储容量，单位MB
	maxBackup  = 3   // 日志文件最大备份数
	maxAge     = 7   // 备份日志文件最大存储天数
	logName    = "pcdn_node.log"
)

var (
	instance *Logger
	once     sync.Once
)

type Logger struct {
	*logrus.Logger
}

// newLogger 初始化一个日志服务
func newLogger() *Logger {
	logger := &Logger{&logrus.Logger{
		Out:          os.Stderr,
		Formatter:    &logrus.TextFormatter{FullTimestamp: true},
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}}
	return logger
}

// setLevel 根据输入的级别字符串 (debug, info, ...)
// 设置日志级别，无法识别则导致 panic
func (l *Logger) SetLevel(level string) {
	lv, err := logrus.ParseLevel(level)
	if err != nil {
		panic(err)
	}
	l.Logger.SetLevel(lv)
}

// SetFileOutPut 根据指定的存储路径设置日志文件的存储路径，同时使用 lumberjack 实现日志文件的循环记录
// 如果循环记录的hook失败，直接panic退出
func (l *Logger) SetFileOutPut(logPath string, level string) {
	err := os.MkdirAll(logPath, 0755)
	if err != nil {
		l.Fatalln("cannot creat log_dir")
	}
	fileName := logPath + logName
	if level == "" {
		level = "debug"
	}
	lv, err := logrus.ParseLevel(level)
	if err != nil {
		l.Fatalln(err)
	}
	hook, err := NewRotateFileHook(
		RotateFileConfig{
			Filename:   fileName,
			MaxSize:    maxLogSize,
			MaxBackups: maxBackup,
			MaxAge:     maxAge,
			Level:      lv,
			Formatter:  &logFileFormat{},
		})
	if err != nil {
		panic(err)
	}
	l.Logger.AddHook(hook)
}

// GetLogger 保证获取全局唯一的 Logger 单例
// 返回默认的日志服务
func GetLogger() *Logger {
	once.Do(func() { instance = newLogger() })
	return instance
}

type logFileFormat struct {
}

func (f *logFileFormat) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer

	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	b.WriteByte('[')
	b.WriteString(strings.ToUpper(entry.Level.String()))
	b.WriteString("]:")
	b.WriteString(entry.Time.Format("2006-01-02 15:04:05"))

	if entry.Message != "" {
		b.WriteString(" - ")
		b.WriteString(entry.Message)
	}

	if len(entry.Data) > 0 {
		b.WriteString(" || ")
	}
	for key, value := range entry.Data {
		b.WriteString(key)
		b.WriteByte('=')
		b.WriteByte('{')
		fmt.Fprint(b, value)
		b.WriteString("}, ")
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}
