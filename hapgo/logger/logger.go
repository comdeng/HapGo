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
	maxBufferLen = 2 << 0
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

func Debug(msg string, params ...interface{}) {
	write(DEBUG, fmt.Sprintf(msg, params...))
}

func Trace(msg string, params ...interface{}) {
	write(TRACE, fmt.Sprintf(msg, params...))
}

func Notice(msg string, params ...interface{}) {
	write(NOTICE, fmt.Sprintf(msg, params...))
}

func Warn(msg string, params ...interface{}) {
	write(WARNING, fmt.Sprintf(msg, params...))
}

func Fatal(msg string, params ...interface{}) {
	write(FATAL, fmt.Sprintf(msg, params...))
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
	return fmt.Sprintf("%x:%s %d %s %s", flag, now.Format("2006/01/02 15:04:05.999"), appId, getBaseString(), msg)
}

func write(level int, msg string) {
	if level > logLevel {
		return
	}
	buffers[currentIndex] = getLogString(level, msg)
	currentIndex++
	log.Print(msg)
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
