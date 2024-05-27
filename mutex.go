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

type Holders interface {
	Holders() []Holder
}

func NewMutex() Mutex {
	if useDeadlock {
		return &deadlock.Mutex{}
	}
	if logSlowLock {
		mutex := &loggedMutex{}
		mutex.holder.Store(Holder{})
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
	currentHolder := m.holder.Load().(Holder)
	duration := timeNow().Sub(currentHolder.Time)
	if duration >= slowLockThreshold {
		l.Debugf("Mutex held for %v. Locked At %s unlocked At %s", duration, currentHolder.At, getHolder().At)
	}
	m.holder.Store(Holder{})
	m.Mutex.Unlock()
}

func (m *loggedMutex) Holders() []Holder {
	return append([]Holder{}, m.holder.Load().(Holder))
}
