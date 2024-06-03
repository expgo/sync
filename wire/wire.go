package sync

import "github.com/expgo/sync"

//go:generate ag

type Once sync.Once

type Mutex sync.Mutex

type RWMutex sync.RWMutex

type WaitGroup sync.WaitGroup

// NewOnce is a factory of Once
// @Factory
func NewOnce() Once {
	return sync.NewOnce()
}

// NewMutex is a factory of Mutex
// @Factory
func NewMutex() Mutex {
	return sync.NewMutex()
}

// NewRWMutex is a factory of RWMutex
// @Factory
func NewRWMutex() RWMutex {
	return sync.NewRWMutex()
}

// NewWaitGroup is a factory of WaitGroup
// @Factory
func NewWaitGroup() WaitGroup {
	return sync.NewWaitGroup()
}
