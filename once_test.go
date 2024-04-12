package sync

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {
	once := Once{}

	err := once.DoTimeout(100*time.Millisecond, func() error {
		time.Sleep(200 * time.Millisecond)
		return nil
	})

	assert.Equal(t, context.DeadlineExceeded, err)
}
