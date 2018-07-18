panic
=====

A small util pkg to help with common panic things

If
==

```go
//old:
if err := doStuff(); err != nil {
    panic(err)
}
//new:
panic.If(doStuff())
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
    panic.If("uh oh")
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
    panic.If(1)
}(),func(){
    panic.If(2)
}(),func(){
    panic.If(3)
}())

//with a timeout, err == nil
ctx := context.Background()
ctx, cancel := context.WithTimeout(ctx, 1 * time.Second)
defer cancel()
err := panic.SafeGoGroup(func(){
    select {
    case <-time.After(2 * time.Second):
        panic(1)
    case <-ctx.Done():
        return
    }
}(),func(){
    select {
    case <-time.After(2 * time.Second):
        panic(2)
    case <-ctx.Done():
        return
    }
}(),func(){
    select {
    case <-time.After(2 * time.Second):
        panic(3)
    case <-ctx.Done():
        return
    }
}())
```
