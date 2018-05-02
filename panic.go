package panic

import (
	"fmt"
	"sync"
	"runtime/debug"
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
func SafeGoGroup(fs ...func()) {
	if len(fs) < 2 {
		panic(fmt.Errorf("fs must be 2 or more funcs"))
	}
	wg := &sync.WaitGroup{}
	errsMtx := &sync.Mutex{}
	errs := make([]*Error, 0, len(fs))
	for _, f := range fs {
		func(f func()) {
			wg.Add(1)
			go func() {
				defer func() {
					defer wg.Done()
					if rVal := recover(); rVal != nil {
						errsMtx.Lock()
						defer errsMtx.Unlock()
						errs = append(errs, &Error{
							Stack: string(debug.Stack()),
							Value: rVal,
						})
					}
				}()
				f()
			}()
		}(f)
	}
	wg.Wait()
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