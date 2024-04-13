package sync

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
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
	return fmt.Sprintf("at %s goid: %d for %s", h.at, h.goid, timeNow().Sub(h.time))
}

func getHolder() holder {
	_, file, line, _ := runtime.Caller(2)
	file = filepath.Join(filepath.Base(filepath.Dir(file)), filepath.Base(file))
	return holder{
		at:   fmt.Sprintf("%s:%d", file, line),
		goid: goid(),
		time: timeNow(),
	}
}

func goid() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		return -1
	}
	return id
}
