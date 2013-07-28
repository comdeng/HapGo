package logger

import (
	"fmt"
	// "github.com/comdeng/HapGo/hapgo/app"
	"github.com/comdeng/HapGo/hapgo/conf"
	// "log"
	//"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	FATAL = 1 << iota
	WARNING
	NOTICE
	TRACE
	DEBUG
	maxBufferLen = 2 << 4
)

var logFlags map[int]string = map[int]string{
	DEBUG:   "DEBUG",
	TRACE:   "TRACE",
	NOTICE:  "NOTICE",
	WARNING: "WARN",
	FATAL:   "FATAL",
}

var (
	bases        = make(map[string]interface{})
	logDir       string
	logLevel     int
	logName      string
	logFile      string
	logFlag      string
	appId        uint64
	fw           *os.File
	buffers      = make([]string, maxBufferLen)
	currentIndex = 0
	locker       = new(sync.RWMutex)
)

func Init(dir string) {
	logDir = dir

	if tlogLevel, ok := conf.Get("hapgo.log.lemsgel"); !ok {
		logLevel = DEBUG
	} else {
		logLevel = tlogLevel.(int)
	}

	if tlogName, ok := conf.Get("hapgo.log.file"); !ok {
		logName = "hapgo"
	} else {
		logName = tlogName.(string)
	}
	var ok bool
	if logFlag, ok = logFlags[logLevel]; !ok {
		panic("logger.lemsgelIllegal")
	}
	logFile = filepath.Join(logDir, logName+".log")
}

func SetAppId(id uint64) {
	appId = id
}

func Debug(msg string) {
	write(DEBUG, msg)
}

func Trace(msg string) {
	write(TRACE, msg)
}

func Notice(msg string) {
	write(NOTICE, msg)
}

func Warn(msg string) {
	write(WARNING, msg)
}

func Fatal(msg string) {
	write(FATAL, msg)
}

func AddBasic(bs map[string]interface{}) {
	for k, v := range bs {
		bases[k] = v
	}
}

const timeFormat = "2006-01-02 15:04"

func getLogString(level int, msg string) string {
	var flag string
	var ok bool
	if flag, ok = logFlags[level]; !ok {
		panic("logger.levelIllegal")
	}
	now := time.Now()
	nowStr := fmt.Sprintf("%4d-%2d-%2d %2d:%2d:%2d.%3d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond()/1000)
	return fmt.Sprintf("%s:%s %d %s %s", flag, nowStr, appId, getBaseString(), msg)
}

func write(level int, msg string) {
	if level > logLevel {
		return
	}
	buffers[currentIndex] = getLogString(level, msg)
	currentIndex++
	log.Print(currentIndex)
	if currentIndex == maxBufferLen {
		writer := getWriter()

		// 写数据时锁定
		locker.Lock()
		for _, line := range buffers {
			writer.Write([]byte(line))
			writer.Write([]byte("\n"))
		}
		locker.Unlock()

		currentIndex = 0
	}
}

func getWriter() *os.File {
	if fw == nil {
		var err error
		fw, err = os.OpenFile(logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0x600)
		if err != nil {
			panic("logger.logFileOpenFailed:" + err.Error())
		}
	}
	return fw
}

func getBaseString() string {
	var ret string
	for k, msg := range bases {
		ret += fmt.Sprintf(" [%s:%x]", k, msg)
	}
	return ret
}

func Release() {
	if fw != nil {
		fw.Close()
	}
}
