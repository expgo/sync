package sync

import (
	"time"
)

var Opts = struct {
	Debug     bool
	Threshold time.Duration
	log       Logger

	UseDeadlock     bool
	DeadlockTimeout time.Duration
}{
	Debug:           false,
	Threshold:       100 * time.Millisecond,
	UseDeadlock:     false,
	DeadlockTimeout: 10 * time.Second,
	log:             &logger{},
}

func SetLog(log Logger) {
	Opts.log = log
}
