package main

import (
	"github.com/eezz10001/ego"
	"github.com/eezz10001/ego/core/econf"
	"github.com/eezz10001/ego/core/elog"
	"time"
)

//  export EGO_DEBUG=true && go run main.go --config=config.toml
func main() {
	if err := ego.New(ego.WithHang(true)).Invoker(func() error {
		go func() {
			// 循环打印配置
			for {
				time.Sleep(3 * time.Second)
				peopleName := econf.GetString("people.name")
				elog.Info("people info", elog.String("name", peopleName), elog.String("type", "onelineByFileWatch"))
			}
		}()
		return nil
	}).Run(); err != nil {
		elog.Panic("startup", elog.FieldErr(err))
	}
}
