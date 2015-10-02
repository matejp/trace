package trace

import (
	"fmt"
	"path/filepath"
	"runtime"
)

var debug bool

func EnableDebug() {
	debug = true
}

// Wrap takes the original error and wraps it into the Trace struct
// memorizing the context of the error.
func Wrap(err error, args ...interface{}) error {
	t := newTrace(runtime.Caller(1))
	t.error = err
	if len(args) != 0 {
		t.Message = fmt.Sprintf(fmt.Sprintf("%v", args[0]), args[1:]...)
	}
	return t
}

// Errorf is similar to fmt.Errorf except that it captures
// more information about the origin of error, such as
// callee, line number and function that simplifies debugging
func Errorf(format string, args ...interface{}) error {
	t := newTrace(runtime.Caller(1))
	t.error = fmt.Errorf(format, args...)
	return t
}

// Fatalf. If debug is false Fatalf calls Errorf. If debug is
// true Fatalf calls panic
func Fatalf(format string, args ...interface{}) error {
	if debug {
		panic(fmt.Sprintf(format, args))
	} else {
		return Errorf(format, args)
	}
}

func newTrace(pc uintptr, filePath string, line int, ok bool) *TraceErr {
	if !ok {
		return &TraceErr{
			File: "unknown_file",
			Path: "unknown_path",
			Func: "unknown_func",
			Line: 0,
		}
	}
	return &TraceErr{
		File: filepath.Base(filePath),
		Path: filePath,
		Func: runtime.FuncForPC(pc).Name(),
		Line: line,
	}
}

// TraceErr contains error message and some additional
// information about the error origin
type TraceErr struct {
	error
	Message string
	File    string
	Path    string
	Func    string
	Line    int
}

func (e *TraceErr) Error() string {
	return fmt.Sprintf("[%v:%v] %v %v", e.File, e.Line, e.Message, e.error)
}

func (e *TraceErr) OrigError() error {
	return e.error
}

// Error is an interface that helps to adapt usage of trace in the code
// When applications define new error types, they can implement the interface
// So error handlers can use OrigError() to retrieve error from the wrapper
type Error interface {
	OrigError() error
}
