package panic

import (
	"fmt"
	"sync"
	"time"
	"runtime/debug"
)

// If condition is true, then panic with error format string and args
func IfTruef(condition bool, format string, args ...interface{}) {
	IfTrue(condition, fmt.Errorf(format, args...))
}

// If condition is true, then panic with value i
func IfTrue(condition bool, i interface{}) {
	if condition {
		panic(i)
	}
}

// If i is not nil then panic with value i
func If(i interface{}) {
	IfTrue(i != nil, i)
}

// Runs f in a go routine, if a panic happens r will be passed the value from calling recover
func SafeGo(f func(), r func(i interface{})) {
	IfTruef(f == nil, "f must be none nil go routine func")
	IfTruef(r == nil, "r must be none nil recover func")
	go func() {
		defer func() {
			if rVal := recover(); rVal != nil {
				r(rVal)
			}
		}()
		f()
	}()
}

// Runs each of fs in its own go routine and panics with a collection of all the recover
// values if there are any, if a timeout of <=0 is passed in then it will not timeout the group,
// if a timeout of >0 is passed in it will panic after this duration if any go routines are still running.
func SafeGoGroup(timeout time.Duration, fs ...func()) error {
	IfTruef(len(fs) < 2, "fs must be 2 or more funcs")
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
					defer func() {
						recover() //incase doneChan has been closed after a timeout
					}()
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
			case <-doneChan:
				doneCount++
			case <-timer.C:
				errsMtx.Lock()
				defer errsMtx.Unlock()
				return &timeoutErr{
					Timeout:        timeout,
					GoRoutineCount: len(fs),
					ReceivedErrors: append(make([]*err, 0, len(errs)), errs...),
				}
			}
		}
	} else {
		for doneCount < len(fs) {
			select {
			case <-doneChan:
				doneCount++
			}
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

type timeoutErr struct {
	Timeout        time.Duration
	GoRoutineCount int
	ReceivedErrors []*err
}

func (e *timeoutErr) Error() string {
	return fmt.Sprintf("go routine group timed out, timeout %s, go routine count: %d, received errors: %v", e.Timeout, e.GoRoutineCount, e.ReceivedErrors)
}
