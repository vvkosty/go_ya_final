package app

import (
	"log"

	"github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"
)

type ServerConfig struct {
	Host                 string `env:"HOST" envDefault:"localhost"`
	Address              string `env:"RUN_ADDRESS,notEmpty" envDefault:"localhost:8080"`
	DatabaseDsn          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func (c *ServerConfig) LoadEnvs() {
	if err := env.Parse(c); err != nil {
		log.Printf("%+v\n", err)
	}
}

func (c *ServerConfig) ParseCommandLine() {
	if flag.Lookup("a") == nil {
		flag.StringVarP(&c.Address, "a", "a", c.Address, "-a localhost:8080")
	}
	if flag.Lookup("d") == nil {
		flag.StringVarP(&c.DatabaseDsn, "d", "d", c.DatabaseDsn, "-d db_driver://user:pass@domain:port/db_name")
	}
	if flag.Lookup("r") == nil {
		flag.StringVarP(&c.AccrualSystemAddress, "r", "r", c.DatabaseDsn, "-r ?")
	}

	flag.Parse()
}
