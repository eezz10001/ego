package main

import (
	"github.com/eezz10001/ego"
	"github.com/eezz10001/ego/core/elog"
)

//  export EGO_DEBUG=false && go run main.go
func main() {
	err := ego.New().Invoker(func() error {
		elog.Info("logger info", elog.String("gopher", "ego"), elog.String("type", "command"))
		return nil
	}).Run()
	if err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}
