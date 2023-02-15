/*
 *  config.go
 *  Created on 22.02.2021
 *  Copyright (C) 2021 Volkswagen AG, All rights reserved.
 */

package tracing

import (
	"flag"
	"fmt"
	"net/url"
)

// Config contains values needed for the configuration of the tracing.
type Config struct {
	Enabled      bool
	ServiceName  string
	AgentURI     string
	CollectorURI string
}

// IsValid returns true if AgentURI and CollectorURI are given and the CollectorURI is a valid URL.
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

// BindConfig takes a Config and a FlagSet and stores the tracing-relevant flags in the corresponding config fields.
func BindConfig(config *Config, fs *flag.FlagSet, defaultName string) {
	fs.BoolVar(&config.Enabled, "tracing.enabled", false, "enables tracing")
	fs.StringVar(&config.ServiceName, "tracing.serviceName", defaultName, "the service name used for tracing")
	fs.StringVar(&config.AgentURI, "tracing.agent.uri", "", "the jaeger agent endpoint uri")
	fs.StringVar(&config.CollectorURI, "tracing.collector.uri", "", "the jaeger collector endpoint uri")
}
