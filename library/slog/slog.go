package slog

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
)

type log struct {
	*glog.Logger
}

func Init(name ...string) *log {
	return &log{g.Log(name...).Async(true)}
}

func (l *log) Redis() *glog.Logger {
	l.Logger.SetPrefix("[redis]")
	return l.Logger
}

func (l *log) Mongodb() *glog.Logger {
	l.Logger.SetPrefix("[mongodb]")
	return l.Logger
}

func (l *log) Mysql() *glog.Logger {
	l.Logger.SetPrefix("[mysql]")
	return l.Logger
}

func (l *log) Cache() *glog.Logger {
	l.Logger.SetPrefix("[cache]")
	return l.Logger
}

func (l *log) S3() *glog.Logger {
	l.Logger.SetPrefix("[s3]")
	return l.Logger
}

func (l *log) ThemeZip() *glog.Logger {
	l.Logger.SetPrefix("[theme_zip]")
	return l.Logger
}
