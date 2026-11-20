package tools

import (
	"time"
)

type Condition struct {
	ch        chan int
	interrupt chan int
}

func NewCondition() *Condition {
	c := new(Condition)
	c.ch = make(chan int)
	c.interrupt = make(chan int)
	return c
}

func (c *Condition) Await() {
	select {
	case c.ch <- 1:
		return
	case c.interrupt <- 1:
		panic("Interrupt exception")
	}
}

func (c *Condition) AwaitTime(nanos time.Duration) {
	select {
	case c.ch <- 1:
		return
	case <-time.After(nanos):
		return
	case c.interrupt <- 1:
		panic("Interrupt exception")
	}
}

func (c *Condition) Signal() {
	select {
	case _ = <-c.ch:
		return
	case <-time.After(1 * time.Millisecond):
		return
	}
}

func (c *Condition) SignalAll() {
	for {
		select {
		case _ = <-c.ch:
			break
		case <-time.After(1 * time.Millisecond):
			return
		}
	}
}

func (c *Condition) Interrupt() {
	defer c.SignalAll()
	for {
		select {
		case _ = <-c.interrupt:
			break
		case <-time.After(1 * time.Millisecond):
			return
		}
	}
}
