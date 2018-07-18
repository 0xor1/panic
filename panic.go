package panic

import (
	"fmt"
	"sync"
	"runtime/debug"
)

// If condition is true, then panic with error format string and args
func If(condition bool, format string, args ...interface{}) {
	if condition {
		IfNotNil(fmt.Errorf(format, args...))
	}
}

// If err is not nil, panic with value err
func IfNotNil(e interface{}) {
	if e != nil {
		panic(e)
	}
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
func SafeGoGroup(fs ...func()) error {
	If(len(fs) < 2, "fs must be 2 or more funcs")
	doneChan := make(chan bool)
	defer close(doneChan)
	errsMtx := &sync.Mutex{}
	errs := make([]*err, 0, len(fs))
	for _, f := range fs {
		func(f func()) {
			go func() {
				defer func() {
					if rVal := recover(); rVal != nil {
						errsMtx.Lock()
						defer errsMtx.Unlock()
						errs = append(errs, &err{
							Stack: string(debug.Stack()),
							Value: rVal,
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
	if len(errs) > 0{
		return &err{
			Value: errs,
		}
	}
	return nil
}

type err struct {
	Stack string
	Value interface{}
}

func (e *err) Error() string {
	return fmt.Sprintf("%v\n%s", e.Value, e.Stack)
}
