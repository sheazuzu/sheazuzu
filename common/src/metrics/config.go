/*
 * config.go
 * Created on 23.10.2019
 * Copyright (C) 2019 Volkswagen AG, All rights reserved
 *
 */

package metrics

import "flag"

type Config struct {
	Enabled     bool
	ServiceName string
}

func BindConfig(config *Config, fs *flag.FlagSet, defaultName string) {
	fs.BoolVar(&config.Enabled, "metrics.enabled", true, "enables prometheus metrics on /metrics")
	fs.StringVar(&config.ServiceName, "metrics.serviceName", defaultName, "the name used to collect metrics")
}
