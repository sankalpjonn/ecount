package ecount

import (
    "time"
    "sync"
)

type Ecount interface{
    Incr(string)
    Stop()
}

type ecount struct {
  eventCntMap     map[string]int
  eventC          chan string
  stopC           chan bool
  stopEvictC      chan bool
  mu              sync.RWMutex
  wg              sync.WaitGroup
  evictionTime    time.Duration
  beforeEvictHook func(map[string]int)
}

func New(evictionTime time.Duration, beforeEvictHook func(map[string]int)) Ecount {
  ec := &ecount{
    eventCntMap: map[string]int{},
    eventC: make(chan string, 1000),
    stopC: make(chan bool),
    stopEvictC: make(chan bool),
    evictionTime: evictionTime,
    beforeEvictHook: beforeEvictHook,
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
  for {
		select {
		case event := <-self.eventC:
			self.onEvent(event)
		case <-self.stopC:
			self.wg.Done()
			return
		}
	}
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

  self.beforeEvictHook(self.eventCntMap)
  self.eventCntMap = map[string]int{}
}

func(self *ecount) Incr(event string) {
  self.eventC <- event
}

func(self *ecount) Stop() {
  self.stopC <- true
  self.stopEvictC <- true
  self.wg.Wait()
  self.evict()
}
