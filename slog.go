package slog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"sync"
	"time"
)

const (
	infoLevel  = "INF"
	errorLevel = "ERR"
)

type Logger interface {
	Info(msg string, keyvalues ...interface{})
	Error(err error, msg string, keyvalues ...interface{})
}

type lg struct {
	w io.Writer
	p sync.Pool
}

func newlg(w io.Writer) *lg {
	return &lg{
		w: w,
		p: sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
	}
}

type JSON lg

func NewJSON(w io.Writer) Logger {
	return (*JSON)(newlg(w))
}

func (l *JSON) print(msg, level string, kvs ...interface{}) {
	b := l.p.Get().(*bytes.Buffer)
	b.Reset()
	defer l.p.Put(b)

	fmt.Fprintf(b, `{"time":"%s", "level":"%s", "msg":"%s"`, time.Now().Format(time.RFC3339), level, msg)
	for i := 0; i+1 < len(kvs); i += 2 {
		kb, err := json.Marshal(kvs[i])
		if err != nil {
			continue
		}
		vb, err := json.Marshal(kvs[i+1])
		if err != nil {
			continue
		}
		fmt.Fprintf(b, `, %s:%s`, kb, vb)
	}

	b.WriteString("}\n")

	b.WriteTo(l.w)
}

func (l *JSON) Info(msg string, keyvalues ...interface{}) {
	l.print(msg, infoLevel, keyvalues...)
}

func (l *JSON) Error(err error, msg string, keyvalues ...interface{}) {
	l.print(msg, errorLevel, append([]interface{}{"err", err.Error()}, keyvalues...)...)
}

type Text lg

func NewText(w io.Writer) Logger {
	return (*Text)(newlg(w))
}

func (l *Text) print(msg, level string, kvs ...interface{}) {
	b := l.p.Get().(*bytes.Buffer)
	b.Reset()
	defer l.p.Put(b)

	fmt.Fprintf(b, `%s %s msg=%q`, time.Now().Format(time.RFC3339), level, msg)
	for i := 0; i+1 < len(kvs); i += 2 {
		fmt.Fprintf(b, ` %v=%s`, kvs[i], strconv.Quote(fmt.Sprintf("%v", kvs[i+1])))
	}
	b.WriteRune('\n')

	b.WriteTo(l.w)
}

func (l *Text) Info(msg string, keyvalues ...interface{}) {
	l.print(msg, infoLevel, keyvalues...)
}

func (l *Text) Error(err error, msg string, keyvalues ...interface{}) {
	l.print(msg, errorLevel, append([]interface{}{"err", err.Error()}, keyvalues...)...)
}

// StdLogger returns a stdlib log.Logger
// for use with net/http.Server, etc.
// logging at INF level
func StdLogger(l Logger) *log.Logger {
	w := stdlog{l}
	return log.New(w, "", 0)
}

type stdlog struct {
	l Logger
}

func (s stdlog) Write(b []byte) (int, error) {
	s.l.Info(string(b))
	return len(b), nil
}
