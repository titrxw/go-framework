package config

import "time"

type Session struct {
	Lifetime     time.Duration `mapstructure:"life_time" json:"life_time" yaml:"life_time"`
	DbConnection string        `mapstructure:"db_connection" json:"db_connection" yaml:"db_connection"`
}
