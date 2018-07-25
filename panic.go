package panic

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime/debug"
	"sync"
)

// If condition is true, then panic with error format string and args
func If(condition bool, format string, args ...interface{}) {
	if condition {
		IfNotNil(fmt.Errorf(format, args...))
	}
}

// If err is not nil, panic with value err
func IfNotNil(e error) {
	panickedDoingReflectionCheck := true
	defer func() {
		if panickedDoingReflectionCheck {
			recover()
		}
	}()
	if e == nil || reflect.ValueOf(e).IsNil() {
		return
	}
	panickedDoingReflectionCheck = false
	panic(e)
}

// Runs f in a go routine, if a panic happens r will be passed the value from calling recover, f should use
// context.Context to handle freeing of resources if necessary on a timeout/deadline/cancellation signal.
func SafeGo(f func(), r func(i interface{})) {
	If(f == nil, "f must be none nil go routine func")
	If(r == nil, "r must be none nil recover func")
	go func() {
		defer func() {
			if rVal := recover(); rVal != nil {
				r(rVal)
			}
		}()
		f()
	}()
}

// Runs each of fs in its own go routine and returns with a collection of all the recover
// values if there are any, each routine should use context.Context to handle freeing of resources
// if necessary on a timeout/deadline/cancellation signal.
func SafeGoGroup(fs ...func()) *Errors {
	doneChan := make(chan bool)
	defer close(doneChan)
	errsMtx := &sync.Mutex{}
	errs := make([]*Error, 0, len(fs))
	for _, f := range fs {
		func(f func()) {
			go func() {
				defer func() {
					if rVal := recover(); rVal != nil {
						errsMtx.Lock()
						defer errsMtx.Unlock()
						errs = append(errs, &Error{
							StackTrace:   string(debug.Stack()),
							RecoverValue: rVal,
						})
					}
					doneChan <- true
				}()
				f()
			}()
		}(f)
	}
	doneCount := 0
	for doneCount < len(fs) {
		select {
		case <-doneChan:
			doneCount++
		}
	}
	if len(errs) > 0 {
		return &Errors{
			Errors: errs,
		}
	}
	return nil
}

type Errors struct {
	Errors []*Error
}

func (e Errors) Error() string {
	buf := bytes.NewBufferString("")
	for _, err := range e.Errors {
		buf.WriteString(err.Error())
	}
	return buf.String()
}

type Error struct {
	StackTrace   string
	RecoverValue interface{}
}

func (e Error) Error() string {
	return fmt.Sprintf("%v\n%s\n", e.RecoverValue, e.StackTrace)
}
