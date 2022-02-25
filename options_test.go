package rpool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPoolOptions(t *testing.T) {
	var options PoolOption
	WithMaxTasksCount(10)(&options)
	WithWorkerCount(20)(&options)
	WithWorkerPoolSize(30)(&options)
	WithRecoverHandler(func(i any) {})(&options)

	assert.Equal(t, 10, options.MaxTasksCount)
	assert.Equal(t, 20, options.WorkerCount)
	assert.Equal(t, 30, options.WorkerPoolSize)
	assert.NotNil(t, options.RecoverHandler)
}
