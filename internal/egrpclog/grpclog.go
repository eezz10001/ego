package egrpclog

import (
	"sync"

	"github.com/eezz10001/ego/core/elog"
)

var (
	once   sync.Once
	logger *elog.Component
)

// Build 构建日志
func Build() *elog.Component {
	once.Do(func() {
		logger = elog.EgoLogger.With(elog.FieldComponentName("component.grpc"))
	})
	return logger
}
