# ecount

A simple module for tracking number of events occured.

Event map gets evicted at a specified interval and user can provide a before evict hook.


## example
```go
package main

import (
  "log"
  "time"
  "github.com/sankalpjonn/ecount"
)

const (
  EVICTION_TIME = time.Second * 1
)

func main() {
  beforeEvictHook := func(eventCntMap map[string]int) {
   log.Println("event map before eviction: ", eventCntMap)
  }

  ec := ecount.New(EVICTION_TIME, beforeEvictHook)

  for i:=0; i<10000; i++ {
    ec.Incr("test-event-1")
  }
  for i:=0; i<10000; i++ {
    ec.Incr("test-event-2")
  }

  time.Sleep(time.Second * 3)
  ec.Stop()
}

// output
// 2019/02/28 17:00:02 event map before eviction:  map[test-event-1:10000 test-event-2:10000]
// 2019/02/28 17:00:03 event map before eviction:  map[]
// 2019/02/28 17:00:04 event map before eviction:  map[]
// 2019/02/28 17:00:04 event map before eviction:  map[]
```
