package ecount

import (
  "log"
  "time"
  "testing"
  "sync"
)

func TestEventCounter(*testing.T) {
  ec := New(time.Second * 1)
  var wg sync.WaitGroup
  wg.Add(1)
  go func(){
      for m := range ec.Evicted() {
        log.Println("got evicted: ", m)
      }
      wg.Done()
  }()

  log.Println("started incrementing")
  for i:=0; i< 10000001; i++ {
    ec.Incr("one")
  }

  for i:=0; i< 10000002; i++ {
    ec.Incr("two")
  }

  for i:=0; i< 10000003; i++ {
    ec.Incr("three")
  }

  log.Println("called stop")
  ec.Stop()
  wg.Wait()
}
