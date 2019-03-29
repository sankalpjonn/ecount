package ecount

import (
    "time"
    "sync"
)

type Ecount interface{
    Incr(string)
    Stop()
    Evicted() chan map[string]int
}

type ecount struct {
  eventCntMap     map[string]int
  eventC          chan string
  evictC          chan map[string]int
  stopC           chan bool
  stopEvictC      chan bool
  mu              sync.RWMutex
  wg              sync.WaitGroup
  evictionTime    time.Duration
}

func New(evictionTime time.Duration) Ecount {
  ec := &ecount{
    eventCntMap: map[string]int{},
    eventC: make(chan string, 1000),
    evictC: make(chan map[string]int, 1000),
    stopC: make(chan bool),
    stopEvictC: make(chan bool),
    evictionTime: evictionTime,
  }
  go ec.startEventListener()
  ec.wg.Add(1)

  go ec.startEvictor()
  ec.wg.Add(1)
  return ec
}

func(self *ecount) onEvent(event string) {
    self.mu.Lock()
    defer self.mu.Unlock()

    if _, ok := self.eventCntMap[event]; !ok {
      self.eventCntMap[event] = 0
    }
    self.eventCntMap[event]++
}


func(self *ecount) startEventListener() {
  for event := range self.eventC {
    self.onEvent(event)
  }

  self.wg.Done()
}

func(self *ecount) startEvictor() {
  ticker := time.NewTicker(self.evictionTime)
	for {
		select {
		case <-ticker.C:
			self.evict()
		case <-self.stopEvictC:
			ticker.Stop()
      self.wg.Done()
			return
		}
	}
}

func(self *ecount) evict() {
  self.mu.Lock()
  defer self.mu.Unlock()

  self.evictC <- self.eventCntMap
  self.eventCntMap = map[string]int{}
}

func(self *ecount) Incr(event string) {
  self.eventC <- event
}

func(self *ecount) Stop() {
  close(self.eventC)
  self.stopEvictC <- true
  self.wg.Wait()
  self.evict()
  close(self.evictC)
}

func(self *ecount) Evicted() chan map[string]int {
  return self.evictC
}
