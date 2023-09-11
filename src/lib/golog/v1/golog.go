package golog

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"sync"
	"time"
)

// ■■■■■■■■■■ Log level ■■■■■■■■■■

type TypLogLevel byte

const (
	LogLvlEmerg  TypLogLevel = 0
	LogLvlAlert  TypLogLevel = 1
	LogLvlCrit   TypLogLevel = 2
	LogLvlErr    TypLogLevel = 3
	LogLvlWarn   TypLogLevel = 4
	LogLvlNotice TypLogLevel = 5
	LogLvlInfo   TypLogLevel = 6
	LogLvlDebug  TypLogLevel = 7
)

// ■■■■■■■■■■ Log structure ■■■■■■■■■■

type TypLog struct {
	// Level of logging
	logLvl TypLogLevel

	// Output
	outputs      []io.Writer
	outputWriter io.Writer

	// Competitive write access
	mutex sync.Mutex
}

// ■■■■■■■■■■ Constructor ■■■■■■■■■■

func NewLogger(_LogLvl TypLogLevel) (*TypLog, error) {
	log := TypLog{
		logLvl:       _LogLvl,
		outputs:      make([]io.Writer, 0),
		outputWriter: io.MultiWriter(),
	}

	// Mutex to synchronize competitive access
	log.mutex = sync.Mutex{}

	return &log, nil
}

// ■■■■■■■■■■ Add file output ■■■■■■■■■■

func (log *TypLog) FctAddFile(_fileLocation string) error {
	if log.outputs == nil {
		return errors.New("log outputs are closed")
	}

	// Directory existence check
	_fileLocation = filepath.FromSlash(_fileLocation)
	_fileDir := path.Dir(_fileLocation)
	if _, err := os.Stat(_fileDir); os.IsNotExist(err) {
		return err
	}

	/*if _, err := os.Stat(_fileLocation); os.IsNotExist(err) {
		return err
	}*/

	ficLog, err := os.OpenFile(_fileLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	return log.FctAddOutput(ficLog)
}

// ■■■■■■■■■■ Add io output ■■■■■■■■■■

func (log *TypLog) FctAddOutput(_buffer io.Writer) error {
	if log.outputs == nil {
		return errors.New("log outputs are closed")
	}

	log.outputs = append(log.outputs, _buffer)
	log.outputWriter = io.MultiWriter(log.outputs...)
	return nil
}

// ■■■■■■■■■■ Start logging and close outputs ■■■■■■■■■■

func (log *TypLog) Start() {
	log.outputs = nil
	log.FctLog(LogLvlDebug, "■■■■■■■■■■ Debugging starts ■■■■■■■■■■")
}

// ■■■■■■■■■■ Log writing ■■■■■■■■■■

func (log *TypLog) FctLog(_TriggerLvl TypLogLevel, _Str string, _Obj ...interface{}) {
	if _TriggerLvl > log.logLvl {
		return
	}

	log.mutex.Lock()
	defer log.mutex.Unlock()

	if _Obj != nil {
		r := regexp.MustCompile(`[%]`)
		matches := r.FindAll([]byte(_Str), -1)
		if len(matches) != len(_Obj) {
			return
		}
		_Str = fmt.Sprintf(_Str, _Obj...)
	}

	// 2006 : years
	// 01 : month
	// 02 : days
	// 15 : hours (24h format)
	// 04 : minutes
	// 05 : seconds
	dateTime := time.Now().Format("2006/01/02 15:04:05")

	_, execFileFullName, execLine, _ := runtime.Caller(1)
	execFileShortName := filepath.Base(execFileFullName)

	var lvl string
	switch _TriggerLvl {
	case LogLvlDebug:
		lvl = "DEBUG "
	case LogLvlInfo:
		lvl = "INFO  "
	case LogLvlNotice:
		lvl = "NOTICE"
	case LogLvlWarn:
		lvl = "WARN  "
	case LogLvlErr:
		lvl = "ERROR "
	case LogLvlCrit:
		lvl = "CRIT  "
	case LogLvlAlert:
		lvl = "ALERT "
	case LogLvlEmerg:
		lvl = "EMERG "
	}

	strBytes := []byte(fmt.Sprintf("%s %s %s:%d : %v", dateTime, lvl, execFileShortName, execLine, _Str))
	if len(strBytes) == 0 || strBytes[len(strBytes)-1] != '\n' {
		strBytes = append(strBytes, '\n')
	}

	log.outputWriter.Write(strBytes)
}

// ■■■■■■■■■■ Log level change ■■■■■■■■■■

func (log *TypLog) FctDefineLogLevel(_lvlLog TypLogLevel) {
	log.logLvl = _lvlLog
}

// ■■■■■■■■■■ Log location according to the environment ■■■■■■■■■■

func fctLogLocation() string {
	switch runtime.GOOS {
	case "linux":
		return path.Join("/var", "log")
	case "windows":
		return os.Getenv("Temp")
	default:
		// Where binary is executed
		return ""
	}
}
