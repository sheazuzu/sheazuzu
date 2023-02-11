/*
 * config.go
 * Created on 21.02.2020
 * Copyright (C) 2020 Volkswagen AG, All rights reserved
 *
 */

package server

import (
	"flag"
	"strings"
)

type Config struct {
	Port             int
	ContextPath      string
	ProfilingEnabled bool
}

// cut off any leading or trailing slashes, the context path setup can be picky otherwise
func (config *Config) GetContextPath() string {

	if config.ContextPath == "" {
		return ""
	}

	return strings.TrimPrefix(strings.TrimSuffix(config.ContextPath, "/"), "/")
}

func BindConfig(config *Config, fs *flag.FlagSet) {
	fs.IntVar(&config.Port, "server.port", 8080, "The port on which to listen")
	fs.StringVar(&config.ContextPath, "server.contextPath", "", "The context path on which to listen")
	fs.BoolVar(&config.ProfilingEnabled, "profiling.enabled", false, "Enables profiling")
}
