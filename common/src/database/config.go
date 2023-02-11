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
	fs.IntVar(&config.Port, "database.port", 8080, "database port")
	fs.StringVar(&config.DatabaseName, "database.name", "sheazuzu", "database endpoint")
	fs.StringVar(&config.Config, "database.config", "parseTime=true", "database endpoint")
	fs.StringVar(&config.UserName, "database.username", "root", "database username")
	fs.StringVar(&config.Password, "database.password", "password", "database password")

}
