package exception

import (
	"fmt"
	"runtime/debug"
)

type Exception interface {
	GetErrString() string
	GetStackString() string
	GetDefaultLogString() string
}

type exception struct {
	ErrString        string
	StackString      string
	DefaultLogString string
}

func (e *exception) GetErrString() string {
	return e.ErrString
}

func (e *exception) GetStackString() string {
	return e.StackString
}

func (e *exception) GetDefaultLogString() string {
	return e.DefaultLogString
}

//统一异常处理
func CatchError(title string, err interface{}) Exception {
	ex := new(exception)
	ex.ErrString = fmt.Sprintln(err)
	ex.StackString = string(debug.Stack())
	ex.DefaultLogString = title + " CatchError [" + ex.GetErrString() + "] [" + ex.GetStackString() + "]"
	return ex
}
