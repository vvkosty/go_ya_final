package app

import (
	"log"

	"github.com/caarlos0/env/v6"
	flag "github.com/spf13/pflag"
)

type Config struct {
	Host                 string `env:"HOST" envDefault:"localhost"`
	Address              string `env:"RUN_ADDRESS,notEmpty" envDefault:"localhost:8080"`
	DatabaseDsn          string `env:"DATABASE_URI"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	AppKey               string `env:"APP_KEY" envDefault:"123e4567-e89b-12d3-a456-42661417"`
}

func (c *Config) LoadEnvs() {
	if err := env.Parse(c); err != nil {
		log.Printf("%+v\n", err)
	}
}

func (c *Config) ParseCommandLine() {
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
