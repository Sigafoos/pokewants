package logger

import (
	"io"
	"log"
)

func New(out io.Writer) *Logger {
	l := log.New(out, "", log.LstdFlags)
	return &Logger{l: l}
}

type Logger struct {
	l *log.Logger
}

func (l *Logger) Event(eventName string) {
	//l.l.Println("event: " + eventName)
}

func (l *Logger) EventKv(eventName string, kvs map[string]string) {
	l.l.Printf("event: %s: %+v\n", eventName, kvs)
}

func (l *Logger) EventErr(eventName string, err error) error {
	l.l.Printf("event error: %s: %s\n", eventName, err)
	return err
}

func (l *Logger) EventErrKv(eventName string, err error, kvs map[string]string) error {
	l.l.Printf("event error: %s: %s: %+v\n", eventName, err, kvs)
	return err
}

func (l *Logger) Timing(eventName string, nanoseconds int64) {
	//l.l.Printf("event timing: %s: %v\n", eventName, nanoseconds)
}

func (l *Logger) TimingKv(eventName string, nanoseconds int64, kvs map[string]string) {
	//l.l.Printf("event timing: %s: %v: %+v\n", eventName, nanoseconds, kvs)
}
