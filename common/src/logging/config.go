/*
 *  config.go
 *  Created on 22.02.2021
 *  Copyright (C) 2021 Volkswagen AG, All rights reserved.
 */

package logging

import (
	"flag"
	"fmt"
)

var DefaultLogFormat = "json"

// Config contains the attributes needed to configure the Logger.
// Level determines which level a message must have (at least) to be logged.
// Format defines which encoding will be used and influences the encoding config.
type Config struct {
	Level  string
	Format string
}

// IsValid returns true, if the config attributes Level and Format both have valid values.
// Otherwise it returns false and prints advice with the available values.
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

// BindConfig takes a Config and a FlagSet and stores the logging-relevant flags in the corresponding config fields.
func BindConfig(cfg *Config, fs *flag.FlagSet) {
	fs.StringVar(&cfg.Level, "logging.level", "info", "The configurable log level of the service, either 'info', 'debug', 'warn' or 'error'")
	fs.StringVar(&cfg.Format, "logging.format", DefaultLogFormat, "The configurable log format of the service, either 'json', 'console' or 'minimal'")
}
