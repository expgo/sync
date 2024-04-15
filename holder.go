package sync

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

var timeNow = time.Now

type holder struct {
	at   string
	time time.Time
	goid int
}

func (h holder) String() string {
	if h.at == "" {
		return "not held"
	}
	return fmt.Sprintf("at %s GoId: %d for %s", h.at, h.goid, timeNow().Sub(h.time))
}

func getHolder() holder {
	_, file, line, _ := runtime.Caller(2)
	file = filepath.Join(filepath.Base(filepath.Dir(file)), filepath.Base(file))
	return holder{
		at:   fmt.Sprintf("%s:%d", file, line),
		goid: GoId(),
		time: timeNow(),
	}
}
