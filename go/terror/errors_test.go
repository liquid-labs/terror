package terror_test

import (
  "errors"
  "fmt"
  "testing"

  "github.com/stretchr/testify/assert"

  // the package we're testing
  "github.com/Liquid-Labs/terror/go/terror"
)

const testMessage = `It's all foobar!`
var causeError error = errors.New(`I'm the real cause.`)

type testStruct struct {
  createFunc func(string, error) terror.Terror
  code int
}

func TestTerrors(t *testing.T) {
  funcs := []testStruct {
    testStruct{func (msg string, cause error) terror.Terror { return terror.BadRequestError(msg, cause) }, 400},
    testStruct{func (msg string, cause error) terror.Terror { return terror.AuthorizationError(msg, cause) }, 401},
    testStruct{func (msg string, cause error) terror.Terror { return terror.ForbiddenError(msg, cause) }, 403},
    testStruct{func (msg string, cause error) terror.Terror { return terror.NotFoundError(msg, cause) }, 404},
    testStruct{func (msg string, cause error) terror.Terror { return terror.UnprocessableEntityError(msg, cause) }, 422},
    testStruct{func (msg string, cause error) terror.Terror { return terror.ServerError(msg, cause) }, 500},
  }

  for i, tStruct := range funcs {
    terror := tStruct.createFunc(testMessage, causeError)

    // test the user message
    assert.Equal(t, testMessage, terror.Error(), `Terror does not reflect the user message.`)
    // test report on underlying error
    assert.Contains(t, terror.Cause(), fmt.Sprintf(`TestTerrors.func%d`, i + 1), `Annotated cause does not identify source function.`)
    assert.Contains(t, terror.Cause(), `errors_test.go`, `Annotated cause does not identify source file.`)
    assert.Contains(t, terror.Cause(), fmt.Sprintf(`:%d]`, 24 + i), `Annotated cause does not identify source line number.`)
    // test the codes
    assert.Equal(t, tStruct.code, terror.Code())
  }
}
