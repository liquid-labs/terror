package terror

import (
  "fmt"
  "net/http"
  "runtime"
)

type Terror interface {
  Error() string
  Code() int
  HasCause() bool
  Cause() string
  CauseError() error
}

type errorData struct {
  message string
  code int
  cause string
  causeError error
}
func (e errorData) Error() string {
  return e.message
}
func (e errorData) Code() int {
  return e.code
}
func (e errorData) HasCause() bool {
  return e.causeError != nil
}
func (e errorData) Cause() string {
  return e.cause
}
func (e errorData) CauseError() error {
  return e.causeError
}

func annotateError(cause error) string {
  if cause == nil {
    return ``
  }
  // '1' is the 'annotateError' call itself
  // '2' is the error creation point
  pc, file, line, _ := runtime.Caller(2)
  return fmt.Sprintf("(%s[%s:%d]) %s", runtime.FuncForPC(pc).Name(), file, line, cause)
}

func BadRequestError(message string, cause error) errorData {
  return errorData{message, http.StatusBadRequest, annotateError(cause), cause}
}
func ForbiddenError(message string, cause error) errorData {
  return errorData{message, http.StatusForbidden, annotateError(cause), cause}
}
func AuthorizationError(message string, cause error) errorData {
  return errorData{message, http.StatusUnauthorized, annotateError(cause), cause}
}
func NotFoundError(message string, cause error) errorData {
  return errorData{message, http.StatusNotFound, annotateError(cause), cause}
}
func UnprocessableEntityError(message string, cause error) errorData {
  return errorData{message, http.StatusUnprocessableEntity, annotateError(cause), cause}
}
func ServerError(message string, cause error) errorData  {
  return errorData{message, http.StatusInternalServerError, annotateError(cause), cause}
}
