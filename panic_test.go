package panic

import(
	"testing"
	"github.com/stretchr/testify/assert"
	"sync"
	"time"
)

func Test_If(t *testing.T) {
	If(nil)
	defer func(){
		r := recover()
		assert.Equal(t, assert.AnError, r)
	}()
	If(assert.AnError)
}

func Test_SafeGo(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	SafeGo(func(){
		panic(assert.AnError)
	}, func(i interface{}){
		defer wg.Done()
		assert.Equal(t, assert.AnError, i)
	})
	wg.Wait()
}

func Test_SafeGoGroup(t *testing.T) {
	e := SafeGoGroup(0, func(){
		panic(0)
	}, func(){
		panic(1)
	}, func(){
		panic(2)
	})

	assert.Equal(t, 3, len(e.(*err).Value.([]*err)))
	idxIsPresent := []bool{false, false, false}
	for _, e := range e.(*err).Value.([]*err) {
		idxIsPresent[e.Value.(int)] = true
	}
	assert.True(t, idxIsPresent[0] && idxIsPresent[1] && idxIsPresent[2])
	e.(*err).Error()

	assert.Nil(t, SafeGoGroup(2 * time.Second, func(){
		time.Sleep(time.Second)
	}, func(){
		time.Sleep(time.Second)
	}, func(){
		time.Sleep(time.Second)
	}))

	e = SafeGoGroup(time.Second, func(){
		panic(0)
	}, func(){
		panic(1)
	}, func(){
		time.Sleep(2 * time.Second)
	})

	assert.Equal(t, 3, e.(*timeoutErr).GoRoutineCount)
	assert.Equal(t, time.Second, e.(*timeoutErr).Timeout)
	assert.Equal(t, 2, len(e.(*timeoutErr).ReceivedErrors))
	e.(*timeoutErr).Error()
}