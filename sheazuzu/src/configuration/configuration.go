package configuration

import (
	"flag"
	"sheazuzu/common/src/database"
	"sheazuzu/common/src/logging"
	"sheazuzu/common/src/server"
)

type Configuration struct {
	Server   server.Config
	Logging  logging.Config
	Database database.Config
}

func New() *Configuration {
	return &Configuration{}
}

func (cfg *Configuration) SetupFlags(serviceName string) *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	server.BindConfig(&cfg.Server, fs)
	logging.BindConfig(&cfg.Logging, fs)
	database.BindConfig(&cfg.Database, fs)

	return fs
}

func (cfg *Configuration) Validate() bool {

	hasErrors := false
	hasErrors = !cfg.Logging.IsValid() || hasErrors

	return !hasErrors
}
