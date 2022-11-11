package ejob

import (
	"testing"

	"github.com/eezz10001/ego/core/elog"
	"github.com/stretchr/testify/assert"
)

func TestComponent_new(t *testing.T) {
	comp := newComponent("test-cmp", defaultConfig(), elog.EgoLogger)
	assert.Equal(t, "test-cmp", comp.Name())
	assert.Equal(t, "task.ejob", comp.PackageName())
}
