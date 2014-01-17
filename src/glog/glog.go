package glog

import (
	"log"
	"os"
	"time"
	"util"
)

const (
	Ldebug = iota
	Linfo
	Lwarn
	Lerror
	Lpanic
	Lfatal
)

var Levels = []string{
	"[DEBUG]",
	"[INFO]",
	"[WARN]",
	"[ERROR]",
	"[PANIC]",
	"[FATAL]",
}

type Logger struct {
	ch      chan string
	fname   string
	logger  *log.Logger
	fhanler *os.File
	fsize   int
	loglv   int32
}

//max log file size(bytes)
const (
	MAX_LOGFILE_SIZE = 5 * 1024 * 1024
)

//replace log.Logger and file handler
func (l *Logger) resetLogger(f *os.File, ll *log.Logger) {
	l.logger = ll
	l.fhanler = f
	l.fsize = 0
}

//create a log file and log.Logger
func createFL(fname string) (*os.File, *log.Logger) {
	path := "logs/" + fname
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("failed to create logfile: %s: %s", fname, err.Error())
		return nil, nil
	}
	l := log.New(f, "", log.LstdFlags)

	return f, l
}

//create glog.Logger, a proxy of log.Logger
//when log file size exceeds MAX_LOGFILE_SIZE,close the log file handler,
//create a new log.Logger and a new log file
func NewLogger(fname string, lv int32) (l *Logger) {
	_now := time.Now()
	tstr := util.Time2Str(&_now)
	newName := fname + tstr
	f, lg := createFL(newName)
	if f == nil || lg == nil {
		log.Fatalf("failed to create logfile: %s", newName)
		return nil
	}

	l = &Logger{
		fname:   fname,
		logger:  lg,
		fhanler: f,
		fsize:   0,
		loglv:   lv,
	}

	l.ch = make(chan string, 2048)

	return l
}

func (l *Logger) PrintLog(str string) {
	lgLen := len(str)
	l.logger.Print(str)
	l.fsize += lgLen
	if l.fsize >= MAX_LOGFILE_SIZE {
		l.createNewLogger()
	}
}

//push logs into ch
func (l *Logger) PushLog(str string) {
	l.ch <- str
}

//Get Logger.ch
func (l *Logger) CH() chan string {
	return l.ch
}

//get logger lv
func (l *Logger) GetLogLv() int32 {
	return l.loglv
}

//close old file handler,create a new log.Logger
func (l *Logger) createNewLogger() {
	_now := time.Now()
	tstr := util.Time2Str(&_now)

	l.fhanler.Close()

	fname := l.fname + tstr
	f, lg := createFL(fname)
	if f == nil || lg == nil {
		return
	}
	l.resetLogger(f, lg)
}
