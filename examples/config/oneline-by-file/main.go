package main

import (
	"github.com/eezz10001/ego"
	"github.com/eezz10001/ego/core/econf"
	"github.com/eezz10001/ego/core/elog"
)

// export EGO_DEBUG=true && go run main.go  --config=config.toml --watch=false
func main() {
	if err := ego.New().Invoker(func() error {
		peopleName := econf.GetString("people.name")
		elog.Info("people info", elog.String("name", peopleName), elog.String("type", "onelineByFile"))
		return nil
	}).Run(); err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}
