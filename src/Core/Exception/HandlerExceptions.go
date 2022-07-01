package exception

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"runtime"

	"go.uber.org/zap"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

type HandlerExceptions struct {
	ExceptionHandler ExceptionHandlerInterface
	Logger           *zap.Logger
}

func (handlerExceptions *HandlerExceptions) SetExceptionHandler(exceptionHandler ExceptionHandlerInterface) {
	handlerExceptions.ExceptionHandler = exceptionHandler
}

func (handlerExceptions *HandlerExceptions) GetExceptionHandler() ExceptionHandlerInterface {
	if handlerExceptions.ExceptionHandler == nil {
		handlerExceptions.ExceptionHandler = new(ExceptionHandler)
	}
	return handlerExceptions.ExceptionHandler
}

func (handlerExceptions *HandlerExceptions) RegisterExceptionHandle() func() {
	return func() {
		if err := recover(); err != nil {
			handlerExceptions.GetExceptionHandler().Reporter(handlerExceptions.Logger, fmt.Errorf("%v", err), string(handlerExceptions.Stack(4)))
			handlerExceptions.GetExceptionHandler().Handle(fmt.Errorf("%v", err), string(handlerExceptions.Stack(4)))
		}
	}
}

//以下代码来自gin recover
func (handlerExceptions *HandlerExceptions) Stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", handlerExceptions.function(pc), handlerExceptions.source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func (handlerExceptions *HandlerExceptions) source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func (handlerExceptions *HandlerExceptions) function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
