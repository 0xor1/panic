package panic

import (
	"fmt"
	"sync"
	"runtime/debug"
	"time"
)

// If i is not nil, then panic with value i
func If(i interface{}) {
	if i != nil {
		panic(i)
	}
}

// Runs f in a go routine, if a panic happens r will be passed the value from calling recover
func SafeGo(f func(), r func(i interface{})) {
	if f == nil {
		panic(fmt.Errorf("f must be none nil go routine func"))
	}
	if r == nil {
		panic(fmt.Errorf("r must be none nil recover func"))
	}
	go func(){
		defer func(){
			if rVal := recover(); rVal != nil {
				r(rVal)
			}
		}()
		f()
	}()
}

// Runs each of fs in its own go routine, waits for them all to complete and panics with a collection of all the recover
// values if there are any
func SafeGoGroup(timeout time.Duration, fs ...func()) {
	if len(fs) < 2 {
		panic(fmt.Errorf("fs must be 2 or more funcs"))
	}
	doneChan := make(chan bool)
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
	if timeout > 0 {
		timer := time.NewTimer(timeout)
		for doneCount < len(fs) {
			select {
			case <- doneChan:
				doneCount++
			case <- timer.C:
				errsMtx.Lock()
				defer errsMtx.Unlock()
				panic(&TimeoutError{
					Timeout: timeout,
					GoRoutineCount: len(fs),
					ReceivedErrors: append(make([]*Error,0, len(errs)), errs...),
				})
			}
		}
	} else {
		for doneCount < len(fs) {
			select {
			case <- doneChan:
				doneCount++
			}
		}
	}
	if len(errs) > 0 {
		panic(&Error{
			Value: errs,
		})
	}
}

type Error struct{
	Stack string
	Value interface{}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v\n%s", e.Value, e.Stack)
}

type TimeoutError struct{
	Timeout time.Duration
	GoRoutineCount int
	ReceivedErrors []*Error
}

func (e *TimeoutError) Error() string {
	return fmt.Sprintf("go routine group timed out, timeout %s, go routine count: %d, received errors: %v", e.Timeout, e.GoRoutineCount, e.ReceivedErrors)
}