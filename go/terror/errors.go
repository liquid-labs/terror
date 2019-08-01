package terror

import (
  "fmt"
  "net/http"
  "runtime"
)

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
func annotateError(cause error) string {
  if cause == nil {
    return ``
  }
  // '1' is the 'annotateError' call itself
  // '2' is the error creation point
  pc, file, line, _ := runtime.Caller(2)
  return fmt.Sprintf("(%s[%s:%d]) %s", runtime.FuncForPC(pc).Name(), file, line, cause)
}

// BadRequestError indicates a malformed request. See [RFC2616](https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html#sec10.4.1)
func BadRequestError(message string, cause error) errorData {
  return errorData{message, http.StatusBadRequest, annotateError(cause), cause}
}
// AuthorizationError indicates a properly formed request lacking proper authorization. See [RFC2616](https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html#sec10.4.2)
func AuthorizationError(message string, cause error) errorData {
  return errorData{message, http.StatusUnauthorized, annotateError(cause), cause}
}
// ForbiddenError indicates a properly formed request which is forbidden regardless of authorization. See [RFC2616](https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html#sec10.4.4)
func ForbiddenError(message string, cause error) errorData {
  return errorData{message, http.StatusForbidden, annotateError(cause), cause}
}
// NotFoundError indicates a properly formed request for something which is not there. See [RFC2616](https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html#sec10.4.5)
func NotFoundError(message string, cause error) errorData {
  return errorData{message, http.StatusNotFound, annotateError(cause), cause}
}
// UnprocessableEntityError indicates a properly formed request with semantically invalid content. See [RFC4918](https://tools.ietf.org/html/rfc4918#page-78)
func UnprocessableEntityError(message string, cause error) errorData {
  return errorData{message, http.StatusUnprocessableEntity, annotateError(cause), cause}
}
// ServerError indicates an unexpected server side error. See [RFC2616](https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html#sec10.5.1)
func ServerError(message string, cause error) errorData  {
  return errorData{message, http.StatusInternalServerError, annotateError(cause), cause}
}
