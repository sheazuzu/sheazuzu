/*
 * config.go Created on 13.12.2021Copyright (C) 2021 Volkswagen AG, All rights reserved.
 */

package mongo

import (
	"flag"
	"fmt"
)

const (
	DefaultTimeout = 10 // The default timeout of the mongo DB connection in seconds
)

// Config contains properties needed for the configuration of the database connection.
type Config struct {
	URI               string
	Database          string
	Timeout           int
	UseSSL            bool
	SSLClientCertFile string
	SSLClientKeyFile  string
}

// BindConfig takes a Config and a FlagSet and stores the flags relevant for the mongo db connection in the corresponding config fields.
func BindConfig(config *Config, fs *flag.FlagSet) {

	fs.StringVar(&config.URI, "mongo.uri", "", "the mongodb uri")
	fs.StringVar(&config.Database, "mongo.database", "", "the mongodb database")
	fs.IntVar(&config.Timeout, "mongo.timeout", DefaultTimeout, "the mongodb connection timeout in sec.")
	fs.BoolVar(&config.UseSSL, "mongo.useSSL", false, "use SSL with mongo")
	fs.StringVar(&config.SSLClientCertFile, "mongo.sslClientCertFile", "", "the mongodb sslClientCertFile")
	fs.StringVar(&config.SSLClientKeyFile, "mongo.sslClientKeyFile", "", "the mongodb sslClientKeyFile")
}

// IsValid checks if the config properties URI, Database and SSLClientCertFile (in case of UseSSL=true) are set.
func (config *Config) IsValid() bool {

	if config.URI == "" {
		fmt.Println("please specify a mongodb uri")
		return false
	}

	if config.Database == "" {
		fmt.Println("please specify a mongodb database name")
		return false
	}

	if config.UseSSL && config.SSLClientCertFile == "" {
		fmt.Println("please specify a path to the SSL client certificate file (PEM)")
		return false
	}

	return true
}
