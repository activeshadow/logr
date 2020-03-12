package stdlogr

import (
	"errors"
	"fmt"
	"io"
	"time"
)

// Mimic github.com/pkg/errors w/out requiring it in the actual module.
type wrapped struct {
	cause error
	msg   string
}

func (this *wrapped) Error() string { return this.msg + ": " + this.cause.Error() }

func (this *wrapped) Cause() error { return this.cause }

func (this *wrapped) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", this.Cause())
			io.WriteString(s, this.msg)
			return
		}

		fallthrough
	case 's', 'q':
		io.WriteString(s, this.Error())
	}
}

func ExampleInfoLogger() {
	mock, _ := time.Parse("2006-01-02", "2015-12-15")

	logger := New("foo")
	logger.(*StdLogr).clock = &clock{mock: mock}
	logger.Info("test log", "hello", "world")
	// Output: level=info ts="2015/12/15 00:00:00" epoch=1450137600 name=foo msg="test log" hello=world
}

func ExampleErrLogger() {
	mock, _ := time.Parse("2006-01-02", "2015-12-15")
	err := errors.New("BOOM SUCKA!")

	logger := New("bar")
	logger.(*StdLogr).clock = &clock{mock: mock}
	logger.Error(err, "test error log", "hello", "world")
	// Output: level=error ts="2015/12/15 00:00:00" epoch=1450137600 name=bar msg="test error log" hello=world error="BOOM SUCKA!"
}

func ExampleWrappedErrLogger() {
	mock, _ := time.Parse("2006-01-02", "2015-12-15")
	err := errors.New("BOOM SUCKA!")

	err = &wrapped{cause: err, msg: "foo bar"}

	logger := New("bar")
	logger.(*StdLogr).clock = &clock{mock: mock}
	logger.Error(err, "test error log", "hello", "world")
	// Output: level=error ts="2015/12/15 00:00:00" epoch=1450137600 name=bar msg="test error log" hello=world error="foo bar: BOOM SUCKA!"
}

func ExampleNonVerboseLogger() {
	mock, _ := time.Parse("2006-01-02", "2015-12-15")

	logger := New("sucka")
	logger.(*StdLogr).clock = &clock{mock: mock}
	logger.V(1).Info("test verbose log", "hello", "crazy world")
	// Output:
}

func ExampleVerboseLogger() {
	mock, _ := time.Parse("2006-01-02", "2015-12-15")

	SetVerbosity(1)

	logger := New("sucka")
	vLogger := logger.V(1)
	vLogger.(*StdLogr).clock = &clock{mock: mock}
	vLogger.Info("test verbose log", "hello", "crazy world")
	// Output: level=info ts="2015/12/15 00:00:00" epoch=1450137600 name=sucka msg="test verbose log" v=1 hello="crazy world"
}

func ExampleNamedLogger() {
	mock, _ := time.Parse("2006-01-02", "2015-12-15")

	logger := New("foo")
	namedLogger := logger.WithName("bar")
	namedLogger.(*StdLogr).clock = &clock{mock: mock}
	namedLogger.Info("test log", "hello", "world")
	// Output: level=info ts="2015/12/15 00:00:00" epoch=1450137600 name=foo.bar msg="test log" hello=world
}

func ExampleValuesLogger() {
	mock, _ := time.Parse("2006-01-02", "2015-12-15")

	logger := New("foo")
	valuesLogger := logger.WithValues("goodbye", "crazy world")
	valuesLogger.(*StdLogr).clock = &clock{mock: mock}
	valuesLogger.Info("test log", "hello", "world")
	// Output: level=info ts="2015/12/15 00:00:00" epoch=1450137600 name=foo msg="test log" goodbye="crazy world" hello=world
}

func ExampleLimitedLogger() {
	mock, _ := time.Parse("2006-01-02", "2015-12-15")

	SetVerbosity(1)
	LimitToLoggers("bar")

	logger := New("foo")
	vLogger := logger.V(1)
	vLogger.(*StdInfoLogr).clock = &clock{mock: mock}
	vLogger.Info("test verbose log", "hello", "crazy world")
	// Output:
}
