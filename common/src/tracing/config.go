/*
 * config.go
 * Created on 23.10.2019
 * Copyright (C) 2019 Volkswagen AG, All rights reserved
 *
 */

package tracing

import (
	"flag"
	"fmt"
	"net/url"
)

type Config struct {
	Enabled      bool
	ServiceName  string
	AgentURI     string
	CollectorURI string
}

func (cfg *Config) IsValid() bool {

	if !cfg.Enabled {
		return true
	}

	hasErrors := false

	if cfg.AgentURI == "" && cfg.CollectorURI == "" {
		fmt.Println("Please provide a tracing agent or collector URI!")
		hasErrors = true
	}

	if cfg.CollectorURI != "" {
		_, err := url.ParseRequestURI(cfg.CollectorURI)
		if err != nil {
			hasErrors = true
			fmt.Println("Please provide a valid tracing collector URI!")
		}
	}

	return !hasErrors
}

func BindConfig(config *Config, fs *flag.FlagSet, defaultName string) {
	fs.BoolVar(&config.Enabled, "tracing.enabled", false, "enables tracing")
	fs.StringVar(&config.ServiceName, "tracing.serviceName", defaultName, "the service name used for tracing")
	fs.StringVar(&config.AgentURI, "tracing.agent.uri", "", "the jaeger agent endpoint uri")
	fs.StringVar(&config.CollectorURI, "tracing.collector.uri", "", "the jaeger collector endpoint uri")
}
