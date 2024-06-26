package sync

import (
	"fmt"
	"sync"
)

type Logger interface {
	Debugf(template string, args ...any)
}

type logger struct {
	mu          sync.Mutex
	enableDebug bool
	hook        func(msg string)
}

func (l *logger) Debugf(template string, args ...any) {
	if !l.enableDebug {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	msg := fmt.Sprintf(template, args...)
	fmt.Println(msg)
	if l.hook != nil {
		l.hook(msg)
	}
}
