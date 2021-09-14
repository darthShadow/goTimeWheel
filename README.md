## goTimeWheel

[![Build Status](https://travis-ci.org/darthShadow/goTimeWheel.svg)](https://travis-ci.org/darthShadow/goTimeWheel)
[![codecov](https://codecov.io/gh/darthShadow/goTimeWheel/branch/master/graph/badge.svg)](https://codecov.io/gh/darthShadow/goTimeWheel)
[![GoDoc](https://godoc.org/github.com/darthShadow/goTimeWheele?status.svg)](https://godoc.org/github.com/darthShadow/goTimeWheel)

TimeWheel Implemented By Go.

![goTimeWheel](https://github.com/darthShadow/goTimeWheel/blob/master/goTimeWheel.png)

Feature
--------

* Effective at Space Usage
* Each Timer can Customize its Task
* Add, Remove, Refresh & Get Task


Installation
-------------

```
go get github.com/darthShadow/goTimeWheel
```

Example
-------

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/darthShadow/goTimeWheel"
)

func main() {
    // timer ticker
    tw := goTimeWheel.New(1*time.Second, 3600)
    tw.Start()

    // "ID1" means the timer's name
    // Specify a function and params, it will run after 3s later
    name := "ID1"
    params := map[string]int{"age": 1}
    fn := func(data interface{}) {
        fmt.Printf("hello, %v\n", data)
    }
    tw.AddTimer(3*time.Second, name, fn, params)

    // Your Logic Code
    select{}

}
```

License
-------

Copyright (c) 2019 by [zheng-ji](http://zheng-ji.info) released under MIT License.
Copyright (c) 2021 by [darthShadow]() released under MIT License.

