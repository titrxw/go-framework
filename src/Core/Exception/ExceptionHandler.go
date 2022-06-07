package exception

import (
	"bytes"
	"fmt"
	server "github.com/titrxw/go-framework/src/Core/Server"
	"io/ioutil"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ExceptionHandler struct {
	ExceptionHandlerInterface
	server.Response
}

func (this *ExceptionHandler) Reporter(logger *zap.Logger, err error) {
	logger.Debug(err.Error())
}

func (this *ExceptionHandler) Handle(ctx *gin.Context, err error) {
	if ctx.Request == nil {
		return
	}

	if gin.Mode() == gin.DebugMode {
		this.JsonResponseWithError(ctx, string(this.Stack(4)), http.StatusInternalServerError)
	} else {
		this.JsonResponseWithError(ctx, "系统内部错误", http.StatusInternalServerError)
	}
}

//以下代码来自gin recover

func (this *ExceptionHandler) Stack(skip int) []byte {
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
		fmt.Fprintf(buf, "\t%s: %s\n", this.function(pc), this.source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func (this *ExceptionHandler) source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func (this *ExceptionHandler) function(pc uintptr) []byte {
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
