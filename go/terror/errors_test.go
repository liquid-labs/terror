package terror_test

import (
  "errors"
  "fmt"
  "testing"

  "github.com/stretchr/testify/assert"

  "github.com/Liquid-Labs/strkit/go/strkit"

  // the package we're testing
  "github.com/Liquid-Labs/terror/go/terror"
)

const testMessage = `It's all foobar!`
var causeError error = errors.New(`I'm the real cause.`)

type userErrTestStruct struct {
  createFunc func(string) terror.Terror
  code int
}

type backendErrTestStruct struct {
  createFunc func(string, error) terror.Terror
  code int
}

func checkError(t *testing.T, testMessage string, code int, lno int, terror terror.Terror) {
  // test the user message
  assert.Equal(t, testMessage, terror.Error(), `Terror does not reflect the user message.`)
  // test report on underlying error
  assert.Contains(t, terror.Cause(), `/terror_test.TestTerrors`, `Annotated cause does not identify source function.`)
  assert.Contains(t, terror.Cause(), `errors_test.go`, `Annotated cause does not identify source file.`)
  assert.Contains(t, terror.Cause(), fmt.Sprintf(`:%d]`, lno), `Annotated cause does not identify source line number.`)
  // test the codes
  assert.Equal(t, code, terror.Code())
}

func TestTerrors(t *testing.T) {
  userTerrors := []userErrTestStruct {
    userErrTestStruct{terror.BadRequestError, 400},
    userErrTestStruct{terror.UnauthenticatedError, 401},
    userErrTestStruct{terror.ForbiddenError, 403},
    userErrTestStruct{terror.NotFoundError, 404},
    userErrTestStruct{terror.UnprocessableEntityError, 422},
  }

  backendTerrors := []backendErrTestStruct {
    backendErrTestStruct{terror.ServerError, 500},
  }

  for _, tStruct := range userTerrors {
    funcName := strkit.FuncNameOnly(tStruct.createFunc)
    t.Run(funcName, func (t *testing.T) {
      terror := tStruct.createFunc(testMessage)
      checkError(t, testMessage, tStruct.code, 56/* line no where error created */, terror)
    })
  }

  for _, tStruct := range backendTerrors {
    funcName := strkit.FuncNameOnly(tStruct.createFunc)
    t.Run(funcName, func (t *testing.T) {
      terror := tStruct.createFunc(testMessage, causeError)
      checkError(t, testMessage, tStruct.code, 64/* line no where error created */, terror)
    })
  }
}
