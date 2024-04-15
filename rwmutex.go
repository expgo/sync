package sync

import (
	"github.com/sasha-s/go-deadlock"
	"strings"
	"sync"
	"sync/atomic"
)

type RWMutex interface {
	Mutex
	RLock()
	RUnlock()
}

func NewRWMutex() RWMutex {
	if useDeadlock {
		return &deadlock.RWMutex{}
	}
	if logSlowLock {
		mutex := &loggedRWMutex{
			readHolders: make(map[int][]holder),
			unlockers:   make(chan holder, 1024),
		}
		mutex.holder.Store(holder{})
		return mutex
	}
	return &sync.RWMutex{}
}

type loggedRWMutex struct {
	sync.RWMutex
	holder atomic.Value

	readHolders    map[int][]holder
	readHoldersMut sync.Mutex

	logUnlockers int32
	unlockers    chan holder
}

func (m *loggedRWMutex) Lock() {
	start := timeNow()

	atomic.StoreInt32(&m.logUnlockers, 1)
	m.RWMutex.Lock()
	atomic.StoreInt32(&m.logUnlockers, 0)

	h := getHolder()
	m.holder.Store(h)

	duration := h.time.Sub(start)

	if duration > slowLockThreshold {
		var unlockerStrings []string
	loop:
		for {
			select {
			case h = <-m.unlockers:
				unlockerStrings = append(unlockerStrings, h.String())
			default:
				break loop
			}
		}
		l.Debugf("RWMutex took %v to lock. Locked at %s. RUnlockers while locking:\n%s", duration, h.at, strings.Join(unlockerStrings, "\n"))
	}
}

func (m *loggedRWMutex) Unlock() {
	currentHolder := m.holder.Load().(holder)
	duration := timeNow().Sub(currentHolder.time)
	if duration >= slowLockThreshold {
		l.Debugf("RWMutex held for %v. Locked at %s unlocked at %s", duration, currentHolder.at, getHolder().at)
	}
	m.holder.Store(holder{})
	m.RWMutex.Unlock()
}

func (m *loggedRWMutex) RLock() {
	m.RWMutex.RLock()
	h := getHolder()
	m.readHoldersMut.Lock()
	m.readHolders[h.goid] = append(m.readHolders[h.goid], h)
	m.readHoldersMut.Unlock()
}

func (m *loggedRWMutex) RUnlock() {
	id := GoId()
	m.readHoldersMut.Lock()
	current := m.readHolders[id]
	if len(current) > 0 {
		m.readHolders[id] = current[:len(current)-1]
	}
	m.readHoldersMut.Unlock()
	if atomic.LoadInt32(&m.logUnlockers) == 1 {
		h := getHolder()
		select {
		case m.unlockers <- h:
		default:
			l.Debugf("Dropped holder %s as channel full", h)
		}
	}
	m.RWMutex.RUnlock()
}

func (m *loggedRWMutex) Holders() string {
	output := m.holder.Load().(holder).String() + " (writer)"
	m.readHoldersMut.Lock()
	for _, holders := range m.readHolders {
		for _, h := range holders {
			output += "\n" + h.String() + " (reader)"
		}
	}
	m.readHoldersMut.Unlock()
	return output
}
