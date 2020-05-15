package nologr

import (
	"github.com/go-logr/logr"
)

type NoInfoLogr struct{}

func (NoInfoLogr) Info(string, ...interface{}) {}

func (NoInfoLogr) Enabled() bool {
	return false
}

func (NoInfoLogr) fmtMsg(string, ...interface{}) string {
	return ""
}

type NoLogr struct {
	NoInfoLogr
}

func (NoLogr) Error(error, string, ...interface{}) {}

func (NoLogr) V(int) logr.InfoLogger {
	return new(NoInfoLogr)
}

func (NoLogr) WithValues(...interface{}) logr.Logger {
	return &NoLogr{NoInfoLogr: NoInfoLogr{}}
}

func (NoLogr) WithName(string) logr.Logger {
	return &NoLogr{NoInfoLogr: NoInfoLogr{}}
}

func New() logr.Logger {
	return &NoLogr{NoInfoLogr: NoInfoLogr{}}
}
