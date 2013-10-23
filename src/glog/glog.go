package glog

import (
	"fmt"
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

var levels = []string{
	"[DEBUG]",
	"[INFO]",
	"[WARN]",
	"[ERROR]",
	"[PANIC]",
	"[FATAL]",
}

type Logger struct {
	ch     chan string
	fname  string
	logger *log.Logger
	fd     *os.File
	fsize  int
	loglv  int
}

//if a log file size beyond 10M,create a new one
const (
	MAX_LOGFILE_SIZE = 10 * 1024 * 1024
)

//replace log.Logger and file handler
func (l *Logger) resetLogger(f *os.File, ll *log.Logger) {
	l.logger = ll
	l.fd = f
	l.fsize = 0
}

//create a log file,return the file handler and assign this file handler to a new log.Logger
//return file handler and log.Logger
func createFL(fname string) (*os.File, *log.Logger) {
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("failed to create logfile: %s", fname)
		return nil, nil
	}
	l := log.New(f, "", log.LstdFlags|log.Lshortfile)

	return f, l
}

//create glog.Logger, a proxy of log.Logger
//a glog.Logger contains a log file
//when log file size exceeds MAX_LOGFILE_SIZE,close the log file handler and
//create a new log.Logger with a new log file
func NewLogger(fname string, lv int) (l *Logger) {
	_now := time.Now()
	tstr := util.Time2Str(&_now)
	ln := fname + tstr
	f, ll := createFL(ln)
	if f == nil || ll == nil {
		log.Fatalf("failed to create logfile: %s", ln)
		return nil
	}

	l = &Logger{
		fname:  fname,
		logger: ll,
		fd:     f,
		fsize:  0,
		loglv:  lv}

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
func (l *Logger) PushLog(lv int, str string) {
	if lv >= Ldebug && lv <= Lfatal && lv >= l.loglv {
		newstr := fmt.Sprintf("%s %s", levels[lv], str)
		l.ch <- newstr
	}
}

//Get Logger.ch
func (l *Logger) CH() chan string {
	return l.ch
}

//create a new log.Logger,close formal file handler first
func (l *Logger) createNewLogger() {
	_now := time.Now()
	tstr := util.Time2Str(&_now)

	l.fd.Close()

	fname := l.fname + tstr
	f, ll := createFL(fname)
	if f == nil || ll == nil {
		log.Fatalf("failed to create logfile: %s", fname)
		return
	}
	l.resetLogger(f, ll)
}
