package sync

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

var timeNow = time.Now

type Holder struct {
	At       string
	Time     time.Time
	GoId     int
	ReadOnly bool
}

func (h Holder) String() string {
	if h.At == "" {
		return "not held"
	}
	return fmt.Sprintf("at %s GoId: %d for %s", h.At, h.GoId, timeNow().Sub(h.Time))
}

func getHolder() Holder {
	_, file, line, _ := runtime.Caller(2)
	file = filepath.Join(filepath.Base(filepath.Dir(file)), filepath.Base(file))
	return Holder{
		At:       fmt.Sprintf("%s:%d", file, line),
		GoId:     GoId(),
		Time:     timeNow(),
		ReadOnly: false,
	}
}
