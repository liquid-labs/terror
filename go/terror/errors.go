package terror

import (
  "fmt"
  "log"
  "net/http"
  "os"
  "runtime"
)

const debugKey string = `DEBUG_TERROR`
var debug string = os.Getenv(debugKey)

func EchoErrorLog() {
  debug = `true`
}

// Terror interface supporting immutable results.
type Terror interface {
  Error() string
  Code() int
  HasCause() bool
  Cause() string
  CauseError() error
}

// errorData is the internal struct implementing the Terror interface.
type errorData struct {
  message string
  code int
  cause string
  causeError error
}
// Error implements the Error interface. This returns the user message.
func (e errorData) Error() string {
  return e.message
}
// Code returns the numeric HTTP code which defines the Terror 'type'.
func (e errorData) Code() int {
  return e.code
}
// HasCause indicates whether the Terror has an underlying cause or not. E.g., a Terror may wrap a regular error or be created directly be terror aware implementations.
func (e errorData) HasCause() bool {
  return e.causeError != nil
}
// Cause provides an annoted version of the wrapped error message, if any. This provides the funciton name, file, and line number of the point where the terror was created and is meant exclusively for internal logging. This value should never be expposed to end users.
func (e errorData) Cause() string {
  return e.cause
}
// CauseError provides the wrapped error, if any.
func (e errorData) CauseError() error {
  return e.causeError
}

// annotateError add function, file, and line number information to errors.
func annotateError(cause interface{}) string {
  var causeS string
  switch cause.(type) {
  case string:
    causeS = cause.(string)
  case error:
    causeS = (cause.(error)).Error()
  }
  // '0' is 'annotateError'
  // '1' is the Terror creator function
  // '2' is the user error creation point
  pc, file, line, _ := runtime.Caller(2)
  causeLog := fmt.Sprintf("(%s[%s:%d]) %s", runtime.FuncForPC(pc).Name(), file, line, causeS)
  if debug != `` {
    log.Println(causeLog)
  }
  return causeLog
}

// BadRequestError (400) indicates a malformed request. See [RFC2616](https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html#sec10.4.1)
func BadRequestError(message string) Terror {
  return errorData{message, http.StatusBadRequest, annotateError(`Bad request.`), nil}
}

// UnautthenticatedError (401) indicates a properly formed request that lacks required authentication. Note that the HTTP standards call this "Unauthorized", rather than "unauthenticated". We prefer the latter to maintain a clear distinction between authentication and authorization. Use `ForbiddenError` for unauthorized requests. See [RFC2616](https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html#sec10.4.2)
func UnauthenticatedError(message string) Terror {
  return errorData{message, http.StatusUnauthorized, annotateError(`Lacks required authentication.`), nil}
}

// ForbiddenError (403) indicates a properly formed request by an authenticated with is none-the-less unauthorized. While `UnauthorizedError` would be the better name, this would conflict with the HTTP standards use of the term, so in this case we keep "forbidden" to avoid confusion. See [RFC2616](https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html#sec10.4.4)
func ForbiddenError(message string) Terror {
  return errorData{message, http.StatusForbidden, annotateError(`Unauthorized (forbidden) request.`), nil}
}

// NotFoundError (404) indicates a properly formed request for something which is not there. See [RFC2616](https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html#sec10.4.5)
func NotFoundError(message string) Terror {
  return errorData{message, http.StatusNotFound, annotateError(`Resource not found.`), nil}
}

// MethodNotAllowedError (405) indicates a properly formed request that is understood by the server, but whose method is not allowed or not supported for the target resource. This is distinct from where a method is possible, but not allowed under current circumstances, in which case `UnauthorizedError`, `ForbiddenError`, or others may be appropriate. See [RFC7231](https://tools.ietf.org/html/rfc7231#section-6.5.5)
func MethodNotAllowedError(message string) Terror {
  return errorData{message, http.StatusMethodNotAllowed, annotateError(`Method not allowed.`), nil}
}

// UnprocessableEntityError (422) indicates a properly formed request with semantically invalid content. See [RFC4918](https://tools.ietf.org/html/rfc4918#page-78)
func UnprocessableEntityError(message string) Terror {
  return errorData{message, http.StatusUnprocessableEntity, annotateError(`Unprocessable entity (bad user data).`), nil}
}

// ServerError (500) indicates an unexpected server side error. See [RFC2616](https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html#sec10.5.1)
func ServerError(message string, cause error) Terror  {
  return errorData{message, http.StatusInternalServerError, annotateError(cause), cause}
}
