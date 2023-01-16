package business

import (
	"context"
	"flag"

	"server/internal/module"
)

type App struct {
	module.CLI    // 命令行
	module.SQLite // 数据库
	module.HTTP   // HTTP 服务
}

func Run(ctx context.Context) (err error) {
	return module.Run[App](ctx)
}

func (app *App) PreInit(ctx context.Context, field any) (err error) {
	switch field.(type) {
	case *module.CLI:
		flag.StringVar(&app.SQLite.DSN, "dsn", "file:demo.db", "DB file path")
		flag.StringVar(&app.HTTP.Addr, "addr", ":8888", "Server listen address")
	}
	return
}

func (app *App) PostInit(ctx context.Context, field any) (err error) {
	switch m := field.(type) {
	case *module.CLI:
		module.StringEnv(&app.SQLite.DSN, "dsn")
		module.StringEnv(&app.HTTP.Addr, "addr")
	case *module.HTTP:
		m.DB = app.SQLite.DB
	}
	return
}
