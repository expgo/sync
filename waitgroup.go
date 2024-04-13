package sync

import "sync"

type WaitGroup interface {
	Add(int)
	Done()
	Wait()
}

func NewWaitGroup() WaitGroup {
	if Opts.Debug {
		return &loggedWaitGroup{}
	}
	return &sync.WaitGroup{}
}

type loggedWaitGroup struct {
	sync.WaitGroup
}

func (wg *loggedWaitGroup) Wait() {
	start := timeNow()
	wg.WaitGroup.Wait()
	duration := timeNow().Sub(start)
	if duration >= Opts.Threshold {
		Opts.log.Debugf("WaitGroup took %v at %s", duration, getHolder())
	}
}
