package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
	"unsafe"

	"github.com/google/uuid"

	bflag "github.com/vaiktorg/grimoire/bitflag"
)

type Log struct {
	Timestamp string      `json:"timestamp"`
	Level     Level       `json:"level"`
	Msg       string      `json:"msg"`
	SourceId  string      `json:"sourceId,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

type msgstore []string

type Logger struct {
	sync.Once
	//msgstore  strings.Builder
	msgstore  msgstore
	LogLevels bflag.Bit

	inChan chan Log

	//network//
	connection Connection
}

const (
	DefaultLogSize = 5 * MB
	MaxStoreSize   = DefaultLogSize
)

func NewLogger() *Logger {
	l := &Logger{
		LogLevels: bflag.Bit(LevelTrace | LevelInfo | LevelDebug | LevelError),
		inChan:    make(chan Log),

		//network//
		connection: Connection{
			user:     "logger",
			password: uuid.New().String(),
			appkey:   uuid.New().String(),
		},
	}

	go l.handlePersistence()

	return l
}
func (l *Logger) TRACE(info string, obj ...interface{}) {
	l.newMsg(info, obj, LevelTrace)
}
func (l *Logger) INFO(info string, obj ...interface{}) {
	l.newMsg(info, obj, LevelInfo)
}
func (l *Logger) DEBUG(procstep string, obj ...interface{}) {
	l.newMsg(procstep, obj, LevelDebug)
}
func (l *Logger) WARN(warn string, obj ...interface{}) {
	l.newMsg(warn, obj, LevelWarn)
}
func (l *Logger) ERROR(errmsg error, obj ...interface{}) {
	l.newMsg(errmsg.Error(), obj, LevelError)
}
func (l *Logger) FATAL(breakage string, obj ...interface{}) {
	l.newMsg(breakage, obj, LevelFatal)
}
func (l *Logger) INBOUND(msg Log) {
	if l.HasLevel(msg.Level) {
		l.inChan <- msg
	}
}

func (l *Logger) Messages() string {
	return l.msgstore.Messages()
}
func (l *Logger) Close() {
	close(l.inChan)
	defer l.toFile(l.msgstore.Messages())

	fmt.Println("Logger Closed!")
}

func (l *Logger) newMsg(msg string, data interface{}, level Level) {
	if l.HasLevel(level) {
		l.inChan <- Log{
			Timestamp: time.Now().Format("20060102150405"),
			Level:     level,
			Msg:       msg,
			Data:      data,
		}
	}
}
func (l *Logger) writeMsg(msg Log) {
	// Marshal your json

	data, err := json.Marshal(msg)
	l.must(err)

	//Persist to file
	l.msgstore = append(l.msgstore, string(data))
}
func (l *Logger) handlePersistence() {
	go func() {
		for msg := range l.inChan {
			memsize := int(unsafe.Sizeof(l.msgstore) + unsafe.Sizeof(""))
			if memsize >= MaxStoreSize.Val() {
				l.dumpLog()
			}

			l.writeMsg(msg)
		}
	}()
}

func (l *Logger) dumpLog() {
	l.toFile(l.msgstore.Messages())
	l.msgstore = []string{}
}

func (l *Logger) toFile(msg string) {
	filename := time.Now().Format("20060102150405") + ".log"

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	_, err = file.WriteString(msg)
	l.must(err)

	err = file.Close()
	l.must(err)
}

func (l *Logger) must(e error) {
	defer func() {
		if e, ok := recover().(error); ok {
			_, _ = fmt.Fprint(os.Stdout, e)
		}
	}()

	if e != nil {
		panic(e)
	}
}

//===================================

func (m *msgstore) Messages() string {
	wrtr := bytes.Buffer{}
	_ = json.NewEncoder(&wrtr).Encode(*m)

	return wrtr.String()
}
