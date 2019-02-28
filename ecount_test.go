package ecount

import (
  "log"
  "time"
  "testing"
)

func TestEventCounter(*testing.T) {
  ec := New(time.Second * 1, func(eventCtrMap map[string]int) {
    log.Println("got map : ", eventCtrMap)
  })

  for i:=0; i< 10000; i++ {
    ec.Incr("one")
  }

  ec.Incr("two")
  ec.Incr("two")

  ec.Incr("three")
  ec.Incr("three")
  ec.Incr("three")

  time.Sleep(time.Second * 3)
  ec.Stop()
}
