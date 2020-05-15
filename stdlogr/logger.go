package stdlogr

import (
	"fmt"
	"time"

	"actshad.dev/logr/util"
	"github.com/go-logr/logr"
)

var (
	verbosity int
	loggers   []string
)

func SetVerbosity(v int) {
	verbosity = v
}

func LimitToLoggers(names ...string) {
	loggers = append(loggers, names...)
}

type clock struct {
	mock time.Time // set in tests
}

func (this *clock) now() time.Time {
	if this == nil { // normal operations
		return time.Now()
	}

	return this.mock // testing only
}

type StdInfoLogr struct {
	enabled bool
	name    string
	kvs     map[string]interface{}

	// Normally nil, set by test code only.
	clock *clock
}

func (this StdInfoLogr) Info(msg string, kvs ...interface{}) {
	if !this.enabled {
		return
	}

	fmt.Println("level=info " + this.fmtMsg(msg, kvs...))
}

func (this StdInfoLogr) Enabled() bool {
	return this.enabled
}

func (this StdInfoLogr) fmtMsg(msg string, kvs ...interface{}) string {
	now := this.clock.now().UTC()

	msg = util.QuoteSpaces(msg)
	msg = fmt.Sprintf(`ts="%s" epoch=%d name=%s msg=%s`, now.Format("2006/01/02 15:04:05"), now.Unix(), this.name, msg)

	for k, v := range this.kvs {
		str := util.QuoteSpaces(v)
		msg = fmt.Sprintf("%s %s=%s", msg, k, str)
	}

	for i := 0; i < len(kvs); i += 2 {
		str := util.QuoteSpaces(kvs[i+1])
		msg = fmt.Sprintf("%s %s=%s", msg, kvs[i], str)
	}

	return msg
}

type StdLogr struct {
	StdInfoLogr
}

func (this StdLogr) Error(err error, msg string, kvs ...interface{}) {
	kvs = append(kvs, "error", err.Error())
	fmt.Println("level=error " + this.fmtMsg(msg, kvs...))
}

func (this StdLogr) V(level int) logr.InfoLogger {
	if level <= verbosity {
		if len(loggers) == 0 || util.StringSliceContains(loggers, this.name) {
			return this.WithValues("v", level)
		}
	}

	return &StdInfoLogr{enabled: false}
}

func (this StdLogr) WithValues(kvs ...interface{}) logr.Logger {
	newKVs := make(map[string]interface{})

	for k, v := range this.kvs {
		newKVs[k] = v
	}

	for i := 0; i < len(kvs); i += 2 {
		newKVs[kvs[i].(string)] = kvs[i+1]
	}

	return &StdLogr{
		StdInfoLogr: StdInfoLogr{
			enabled: this.enabled,
			name:    this.name,
			kvs:     newKVs,
		},
	}
}

func (this StdLogr) WithName(name string) logr.Logger {
	name = fmt.Sprintf("%s.%s", this.name, name)

	return &StdLogr{
		StdInfoLogr: StdInfoLogr{
			enabled: this.enabled,
			name:    name,
			kvs:     this.kvs,
		},
	}
}

func New(name string) logr.Logger {
	return &StdLogr{
		StdInfoLogr: StdInfoLogr{
			enabled: true,
			name:    name,
			kvs:     make(map[string]interface{}),
		},
	}
}
