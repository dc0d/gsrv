# gsrv
Go Generic-Code Generator Server

# intro

Watch [here](https://youtu.be/T_v3qh0S1QQ).

# intro, text

The empty interface - `interface{}` - is used, in Go, to represent a value, that its type is unknown, at a certain point, in the code. Sometimes relying on empty interfaces, makes code hard to reason about and the type-safety provided by the compiler, would be of no use.

Since Go has not generics in its type system, the alternative approach is to use code generators, combined with copy/paste.

The problem with the copy/paste approach is, when you find and fix a bug in the original _generic template_ or _generic definition_ code, it's hard to find all the places that the code is copied to, and replace the fixed version.

This tool helps with that. In addition to keeping your _generic definitions_ in sync, and applying changes automatically, by using _blank imports_, it ensures that all the code in the chain - from the main _generic definition_ down the specialized _generic implementation_ - is valid Go code. Because they get compiled by Go compiler!

### get the tool

First `go get` the tool components:

```
$ go get -u github.com/dc0d/ggen
$ go get -u github.com/dc0d/gsrv
```

Then start the server:

```
$ gsrv
```

### generate the skeleton

Then generate a skeleton for defining the generic type:

```
$ ggen create --name mygenericmap
```

Now go on and modify the skeleton and finalize the _generic definition_.

### specialize the generic definition

Create a package and put a blank import of the _generic definition_ inside it. Say we have `intstrmap` package, and inside `intstrmap.go` we add:

```go
import (
    _ "mygenericmap"
)
```

Now this tool recognizes this blank import and copies the necessary files from the original _generic definition_, to this directory and keeps them in sync!

To specialize the original _generic definition_, we can not set the `interface{}`s, to concrete types. We can even specialize the implementation yet everything will remain synced!

To see it in action, watch the [intro](https://youtu.be/T_v3qh0S1QQ).