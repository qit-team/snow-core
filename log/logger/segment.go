package logger

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"io"
	"time"
)

func LogSegment(l *logrus.Logger, logDir, name string) (*logrus.Logger, error) {
	var (
		err         error
		infoWriter  io.Writer
		warnWriter  io.Writer
		errorWriter io.Writer
	)

	if l.GetLevel() >= logrus.DebugLevel {
		infoPath := fmt.Sprintf("%s/%s.%s", logDir, name, "INFO.%Y%m%d.log")
		linkInfoPath := fmt.Sprintf("%s/%s.%s", logDir, name, "INFO.log")
		infoWriter, err = rotatelogs.New(
			infoPath,
			rotatelogs.WithLinkName(linkInfoPath),
			rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
			rotatelogs.WithRotationTime(time.Duration(604800)*time.Second),
		)
		if err != nil {
			return nil, err
		}
	}
	if l.GetLevel() >= logrus.WarnLevel {
		warnPath := fmt.Sprintf("%s/%s.%s", logDir, name, "WARN.%Y%m%d.log")
		linkWarnPath := fmt.Sprintf("%s/%s.%s", logDir, name, "WARN.log")
		warnWriter, err = rotatelogs.New(
			warnPath,
			rotatelogs.WithLinkName(linkWarnPath),
			rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
			rotatelogs.WithRotationTime(time.Duration(604800)*time.Second),
		)
		if err != nil {
			return nil, err
		}
	}
	if l.GetLevel() >= logrus.ErrorLevel {
		errorPath := fmt.Sprintf("%s/%s.%s", logDir, name, "ERROR.%Y%m%d.log")
		linkErrorPath := fmt.Sprintf("%s/%s.%s", logDir, name, "ERROR.log")
		errorWriter, err = rotatelogs.New(
			errorPath,
			rotatelogs.WithLinkName(linkErrorPath),
			rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
			rotatelogs.WithRotationTime(time.Duration(604800)*time.Second),
		)
		if err != nil {
			return nil, err
		}
	}

	writerMap := lfshook.WriterMap{}
	if infoWriter != nil {
		writerMap[logrus.DebugLevel] = infoWriter
		writerMap[logrus.InfoLevel] = infoWriter
	}
	if warnWriter != nil {
		writerMap[logrus.WarnLevel] = warnWriter
	}
	if errorWriter != nil {
		writerMap[logrus.ErrorLevel] = errorWriter
	}

	l.Hooks.Add(lfshook.NewHook(writerMap, &logrus.JSONFormatter{}))
	return l, nil
}
