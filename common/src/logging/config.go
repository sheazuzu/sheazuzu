/*
 * config.go
 * Created on 23.10.2019
 * Copyright (C) 2019 Volkswagen AG, All rights reserved
 *
 */

package logging

import (
	"flag"
	"fmt"
)

var DefaultLogFormat = "json"

type Config struct {
	Level  string
	Format string
}

func (cfg *Config) IsValid() bool {

	hasErrors := false

	if cfg.Level != "info" && cfg.Level != "debug" && cfg.Level != "warn" && cfg.Level != "error" {
		fmt.Println("log level must either be 'info', 'debug', 'warn' or 'error'")
		hasErrors = true
	}

	if cfg.Format != "json" && cfg.Format != "console" && cfg.Format != "minimal" {
		fmt.Println("log format must either be 'json' 'console', or 'minimal'")
		hasErrors = true
	}

	return !hasErrors
}

func BindConfig(cfg *Config, fs *flag.FlagSet) {
	fs.StringVar(&cfg.Level, "logging.level", "info", "The configurable log level of the service, either 'info', 'debug', 'warn' or 'error'")
	fs.StringVar(&cfg.Format, "logging.format", DefaultLogFormat, "The configurable log format of the service, either 'json', 'console' or 'minimal'")
}
