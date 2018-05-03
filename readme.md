panic
=====

A small util pkg to help with common panic things

If
==

Syntax sugar

```go
//old:
if err := doStuff(); err != nil {
 panic(err)
}
//new:
panic.If(doStuff())
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
    panic("uh oh")
}(), func(i interface{}) {
    // i == "uh oh"
})
```

SafeGoGroup
===========

Runs a collection of routines in a wait group and panics with a collection of panicked values

```go
defer func(){
    r := recover()
    // r contains each of the panic values and stack traces for each
}()

//blocking call but safe
panic.SafeGoGroup(0, func(){
    panic(1)
}(),func(){
    panic(2)
}(),func(){
    panic(3)
}())

//with a timeout
panic.SafeGoGroup(2 * time.Second, func(){
    time.Sleep(time.Second)
}(),func(){
    time.Sleep(time.Second)
}(),func(){
    time.Sleep(time.Second)
}())
```