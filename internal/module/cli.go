package module

import (
	"context"
	"flag"
	"os"
	"strings"
)

const (
	EnvPrefix = "APP_"

	FlagLevel = "level"
	FlagDSN   = "dsn"
	FlagAddr  = "addr"
)

type CLI struct {
	Level string
}

func (m *CLI) Init(ctx context.Context) (err error) {
	m.InitFlags()
	flag.Parse()
	m.InitEnvs()
	return
}

func (m *CLI) InitFlags() {
	flag.StringVar(&m.Level, FlagLevel, "info", "Log level")
}

func (m *CLI) InitEnvs() {
	StringEnv(&m.Level, FlagLevel)
}

type CLIServer struct {
	CLI
	Addr string
	DSN  string
}

func (m *CLIServer) Init(ctx context.Context) (err error) {
	m.CLI.InitFlags()
	flag.StringVar(&m.DSN, FlagDSN, "file:demo.db", "DB file path")
	flag.StringVar(&m.Addr, FlagAddr, ":8888", "Server listen address")

	flag.Parse()

	m.CLI.InitEnvs()
	StringEnv(&m.DSN, FlagDSN)
	StringEnv(&m.Addr, FlagAddr)
	return
}

func StringEnv(p *string, key string) {
	key = EnvPrefix + strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
	if value := os.Getenv(key); value != "" {
		*p = value
	}
}
