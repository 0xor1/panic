panic
=====

A small util pkg to help with common panic things

IfTrue
======

Syntax sugar

```go
//old:
if something == somethingElse {
    panic(myError)
}
//new:
panic.IfTrue(something == somethingElse, myError)
```

IfTruef
=======

Syntax sugar for formatted string error

```go
//old:
if something == somethingElse {
    panic(fmt.Errorf("something is wrong %v", something))
}
//new:
panic.IfTruef(something == somethingElse, "something is wrong %v", something)
```

If
==

Syntax sugar for a not nil bool check, typically useful for checking returned errors

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
err := panic.SafeGoGroup(0, func(){
    panic.If(1)
}(),func(){
    panic.If(2)
}(),func(){
    panic.If(3)
}())

//with a timeout, err == nil
err := panic.SafeGoGroup(2 * time.Second, func(){
    time.Sleep(time.Second)
}(),func(){
    time.Sleep(time.Second)
}(),func(){
    time.Sleep(time.Second)
}())
```
