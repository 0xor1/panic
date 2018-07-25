package panic

import (
	"context"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func Test_If(t *testing.T) {
	var e error
	If(e != nil, "an error")
	defer func() {
		r := recover()
		assert.Equal(t, "an error", r.(error).Error())
	}()

	If(true, "an error")
}

func Test_SafeGo(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	SafeGo(func() {
		panic(assert.AnError)
	}, func(i interface{}) {
		defer wg.Done()
		assert.Equal(t, assert.AnError, i)
	})
	wg.Wait()
}

func Test_SafeGoGroup(t *testing.T) {
	e := SafeGoGroup(func() {
		panic(0)
	}, func() {
		panic(1)
	}, func() {
		panic(2)
	})

	assert.Equal(t, 3, len(e.Errors))
	idxIsPresent := []bool{false, false, false}
	for _, e := range e.Errors {
		idxIsPresent[e.RecoverValue.(int)] = true
	}
	assert.True(t, idxIsPresent[0] && idxIsPresent[1] && idxIsPresent[2])
	e.Error()

	assert.Nil(t, SafeGoGroup(func() {
		time.Sleep(time.Second)
	}, func() {
		time.Sleep(time.Second)
	}, func() {
		time.Sleep(time.Second)
	}))

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	e = SafeGoGroup(func() {
		panic(0)
	}, func() {
		panic(1)
	}, func() {
		select {
		case <-time.After(2 * time.Second):
			panic(2)
		case <-ctx.Done():
			panic(3)
		}
	})

	assert.Equal(t, 3, len(e.Errors))
	idxIsPresent = []bool{false, false, false, false}
	for _, e := range e.Errors {
		idxIsPresent[e.RecoverValue.(int)] = true
	}
	assert.True(t, idxIsPresent[0] && idxIsPresent[1] && !idxIsPresent[2] && idxIsPresent[3])
	e.Error()
	assert.Nil(t, SafeGoGroup())
}
