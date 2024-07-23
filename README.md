# go-bplist

A parser using Go to implement bplist file

## First

I would like to thank these two articles for helping me complete this difficult task.

[https://doubleblak.com/blogPost.php?k=plist](https://doubleblak.com/blogPost.php?k=plist)

[https://medium.com/@karaiskc/understanding-apples-binary-property-list-format-281e6da00dbd](https://medium.com/@karaiskc/understanding-apples-binary-property-list-format-281e6da00dbd)

<br />

## Getting Started
#### Installing

Use go get to retrieve the SDK to add it to your project's Go module dependencies.

```
go get github.com/Lifelong-Study/go-bplist
```

<br />


## Quick Examples

```
package main

import (
    "fmt"

    bplist "github.com/Lifelong-Study/go-bplist"
)

func main() {
    // Read bplist format file
    nodes, err := bplist.Parse("info.plist")

    //
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    // save to XML format file
    bplist.Save(nodes, "out.plist")
}
```

