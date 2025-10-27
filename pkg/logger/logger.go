package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type date struct {
	Year  int
	Month time.Month
	Day   int
}

func (d date) String() string {
	return fmt.Sprintf("%d-%02d-%02d", d.Year, d.Month, d.Day)
}

func nowDate() (d date) {
	d.Year, d.Month, d.Day = time.Now().Date()
	return d
}

const logDir = "logs"

type Logger struct {
	l             *log.Logger
	fileCreatedAt date
	logFile       *os.File
}

func NewLogger() (*Logger, error) {
	l := &Logger{
		l: log.New(os.Stderr, "", log.Ltime),
	}

	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		l.l.Fatal("failed to create log directory", err)
	}
	l.setFile()

	go func() {
		for {
			l.cleanLogs()
			time.Sleep(time.Hour * 24)
		}
	}()

	return l, nil
}

func (l *Logger) setFile() {
	now := nowDate()
	if l.fileCreatedAt == now && l.logFile != nil {
		return
	}
	defer l.logFile.Close()

	filename := now.String() + ".log"

	f, err := os.OpenFile(filepath.Join(logDir, filename), os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		l.l.Fatalln("failed to open log file:", err)
	}
	out := io.MultiWriter(f, os.Stdout)
	l.l.SetOutput(out)

	l.fileCreatedAt = now
	l.logFile = f
}

func (l *Logger) Println(v ...any) {
	l.setFile()
	l.l.Println(v...)
}

func (l *Logger) Printf(format string, v ...any) {
	l.setFile()
	l.l.Printf(format, v...)
}

func (l *Logger) Fatal(v ...any) {
	l.setFile()
	l.l.Fatal(v...)
}

func (l *Logger) Close() error {
	return l.logFile.Close()
}

func (l *Logger) cleanLogs() {
	logFiles, err := ioutil.ReadDir(logDir)
	if err != nil {
		l.l.Println("failed to read log directory:", err)
	}
	deleteBefore := time.Now().Add(-time.Hour * 24 * 30)
	for _, f := range logFiles {
		if !f.IsDir() {
			t, err := time.Parse("2006-01-02", strings.TrimRight(f.Name(), ".log"))
			if err != nil {
				continue
			}

			if t.Before(deleteBefore) {
				os.Remove(logDir + "/" + f.Name())
			}
		}
	}
}
