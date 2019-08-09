panic
=====

A small util pkg to help with common panic things

If
==

```go
//old:
if shouldBeOne := returnOne(); shouldBeOne != 1 {
    panic(fmt.Errorf("returnOne did not return one, it returned %d", shouldBeOne))
}
//new:
shouldBeOne := returnOne()
panic.If(shouldBeOne != 1, "returnOne did not return one, it returned %d", shouldBeOne)
```

IfNotNil
========

```go
if err := doStuff(); err != nil {
    panic(err)
}
//new:
panic.IfNotNil(doStuff())
```

SafeGo
======

Prevents your application from crashing if you are running a go routine and a panic occurs in the go routine but there
is no deferred `recover` call

```go
//crash entire application
go func(){
    panic("uh oh")
}()

//application continues as normal
// non blocking
panic.SafeGo(func(){
    panic.If(true, "uh oh")
}(), func(i interface{}) {
    // i == "uh oh"
})
```

SafeGoGroup
===========

Runs a collection of routines in a wait group and returns an error containing all the panicked values and stack traces for each

```go
//blocking call but safe, err contains all the panicked values and stack traces for each
err := panic.SafeGoGroup(func(){
    panic.If(true, 1)
}(),func(){
    panic.If(true, 2)
}(),func(){
    panic.If(true, 3)
}())

//with a timeout, err == nil
ctx := context.Background()
ctx, cancel := context.WithTimeout(ctx, 1 * time.Second)
defer cancel()
err := panic.SafeGoGroup(func(){
    select {
    case <-time.After(2 * time.Second):
        panic.If(true, 1)
    case <-ctx.Done():
        return
    }
}(),func(){
    select {
    case <-time.After(2 * time.Second):
        panic.If(true, 2)
    case <-ctx.Done():
        return
    }
}(),func(){
    select {
    case <-time.After(2 * time.Second):
        panic.If(true, 3)
    case <-ctx.Done():
        return
    }
}())
```
