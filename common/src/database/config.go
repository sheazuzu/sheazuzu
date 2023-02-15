package database

import (
	"flag"
	"fmt"
)

type Config struct {
	Endpoint     string
	Port         int
	DatabaseName string
	Config       string
	UserName     string
	Password     string
}

func (config *Config) GetDatabaseConn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		config.UserName,
		config.Password,
		config.Endpoint,
		config.Port,
		config.DatabaseName,
		config.Config)
	// "root:455279980@tcp(local:3306)/atb?charset=utf8&parseTime=True&loc=Local"
	// or
	// "root:455279980@/atb?charset=utf8&parseTime=True&loc=Local"
}

func BindConfig(config *Config, fs *flag.FlagSet) {
	fs.StringVar(&config.Endpoint, "database.endpoint", "localhost", "database endpoint")
	fs.IntVar(&config.Port, "database.port", 3306, "database port")
	fs.StringVar(&config.DatabaseName, "database.name", "", "database endpoint")
	fs.StringVar(&config.Config, "database.config", "parseTime=true", "database endpoint")
	fs.StringVar(&config.UserName, "database.username", "", "database username")
	fs.StringVar(&config.Password, "database.password", "", "database password")

}

// IsValid checks if the config properties URI, Database and SSLClientCertFile (in case of UseSSL=true) are set.
func (config *Config) IsValid() bool {

	if config.Endpoint == "" {
		fmt.Println("please specify a database endpoint")
		return false
	}

	if config.DatabaseName == "" {
		fmt.Println("please specify a mysql database name")
		return false
	}

	if config.UserName == "" {
		fmt.Println("please specify a username")
		return false
	}

	if config.Password == "" {
		fmt.Println("please specify a password")
		return false
	}

	return true
}
