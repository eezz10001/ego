package main

import (
	"errors"
	"fmt"

	"github.com/eezz10001/ego"
	"github.com/eezz10001/ego/core/elog"
	"github.com/eezz10001/ego/core/etrace"
	"github.com/eezz10001/ego/task/ejob"
	"go.uber.org/zap"
)

// export EGO_DEBUG=true && go run main.go --job=jobrunner  --config=config.toml
func main() {
	if err := ego.New().Job(
		ejob.Job("job1", job1),
		ejob.Job("job2", job2),
	).Run(); err != nil {
		elog.Error("start up", zap.Error(err))
	}
}

func job2(ctx ejob.Context) error {
	fmt.Println("i am error job runner, traceId: ", etrace.ExtractTraceID(ctx.Ctx))
	return errors.New("i am error")
}

func job1(ctx ejob.Context) error {
	fmt.Println("i am job runner, traceId: ", etrace.ExtractTraceID(ctx.Ctx))
	return nil
}
