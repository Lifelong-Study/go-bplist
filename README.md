# go-bplist

A parser using Go to implement bplist file

## First

I would like to thank these two articles for helping me complete this difficult task.

[https://doubleblak.com/blogPost.php?k=plist](https://doubleblak.com/blogPost.php?k=plist)

[https://medium.com/@karaiskc/understanding-apples-binary-property-list-format-281e6da00dbd](https://medium.com/@karaiskc/understanding-apples-binary-property-list-format-281e6da00dbd)



## Getting Started
#### Installing

Use go get to retrieve the SDK to add it to your project's Go module dependencies.

```
go get github.com/Lifelong-Study/go-bplist
```



## Quick Examples

```
package main

import (
    "fmt"

    bplist "github.com/Lifelong-Study/go-bplist"
)

func main() {
    // Read bplist format file
    data, err := bplist.Read("info.plist")

    if err != nil {
        panic(err.Error())
    }

    nodes, err := bplist.Parse(data)

    //
    if err != nil {
        panic(err.Error())
    }

    // save to XML format file
    bplist.Save(nodes, "out.plist")
}
```

