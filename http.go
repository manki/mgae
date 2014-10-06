// Google App Engine support library.
// Heavily inspired by
// http://blog.golang.org/2011/07/error-handling-and-go.html.

// +build appengine

package mgae

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"appengine"
)

type Error struct {
	cause        error
	message      string
	responseCode int
	stackTrace   string
}

// Creates a new instance of Error.
// cause - root cause of this error; can be nil.
// message - error message.
// responseCode - HTTP status code to send to the client.
func NewError(cause error, message string, responseCode int) *Error {
	return &Error{
		cause:        cause,
		message:      message,
		responseCode: responseCode,
		stackTrace:   StackTrace(),
	}
}

// Same as NewError, but uses http.StatusInternalServerError as the response
// code.
func NewInternalError(cause error, message string) *Error {
	return NewError(cause, message, http.StatusInternalServerError)
}

// Returns current call stack trace.  Can be used for error logging.
// TODO: This method doesn't belong in this package; move this out.
func StackTrace() string {
	calls := make([]string, 0)
	for i := 2; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		calls = append(calls, fmt.Sprintf("    %x %s:%d", pc, file, line))
	}
	return strings.Join(calls, "\n")
}

func serveError(w http.ResponseWriter, r *http.Request, err *Error) {
	c := appengine.NewContext(r)
	c.Errorf("%s: %v\n%s", err.message, err.cause, err.stackTrace)
	http.Error(w, err.message, err.responseCode)
}

// Default HTTP handler function.
type Handler func(w http.ResponseWriter, r *http.Request) *Error

// Serves HTTP request r and writes the response to w.
// Implementation of http.HandlerFunc interface.
func (fn Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		serveError(w, r, err)
	}
}

// Validating HTTP handler that calls Validate function before serving a
// request.  If Validate returns a non-nil value, the error is propagated to
// the client.  Otherwise Handler is invoked to serve the request.
type ValidatingHandler struct {
	Validator func(r *http.Request) *Error
	Handler   Handler
}

// Serves HTTP request r and writes the response to w.
// Implementation of http.HandlerFunc interface.
func (vh ValidatingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := vh.Validator(r); err != nil {
		serveError(w, r, err)
		return
	}
	vh.Handler.ServeHTTP(w, r)
}

type PreprocessingHandler struct {
	Preprocess func(w http.ResponseWriter, r *http.Request) *Error
	Handler    Handler
}

// Serves HTTP request r and writes the response to w.
// Implementation of http.HandlerFunc interface.
func (ph PreprocessingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := ph.Preprocess(w, r); err != nil {
		serveError(w, r, err)
		return
	}
	ph.Handler.ServeHTTP(w, r)
}
