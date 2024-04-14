package sync

import (
	"github.com/sasha-s/go-deadlock"
	"sync"
	"sync/atomic"
)

type Mutex interface {
	Lock()
	Unlock()
}

func NewMutex() Mutex {
	if useDeadlock {
		return &deadlock.Mutex{}
	}
	if debug {
		mutex := &loggedMutex{}
		mutex.holder.Store(holder{})
		return mutex
	}
	return &sync.Mutex{}
}

type loggedMutex struct {
	sync.Mutex
	holder atomic.Value
}

func (m *loggedMutex) Lock() {
	m.Mutex.Lock()
	m.holder.Store(getHolder())
}

func (m *loggedMutex) Unlock() {
	currentHolder := m.holder.Load().(holder)
	duration := timeNow().Sub(currentHolder.time)
	if duration >= threshold {
		l.Debugf("Mutex held for %v. Locked at %s unlocked at %s", duration, currentHolder.at, getHolder().at)
	}
	m.holder.Store(holder{})
	m.Mutex.Unlock()
}

func (m *loggedMutex) Holders() string {
	return m.holder.Load().(holder).String()
}
