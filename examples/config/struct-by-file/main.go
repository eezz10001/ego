package main

import (
	"github.com/eezz10001/ego"
	"github.com/eezz10001/ego/core/econf"
	"github.com/eezz10001/ego/core/elog"
)

//  export EGO_DEBUG=true && go run main.go --config=config.toml --watch=false
func main() {
	err := ego.New().Invoker(func() error {
		p := People{}
		err := econf.UnmarshalKey("people", &p)
		if err != nil {
			panic(err.Error())
		}
		elog.Info("people info", elog.String("name", p.Name), elog.String("type", "structByFile"))
		return nil
	}).Run()
	if err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}

// People ...
type People struct {
	Name string
}
