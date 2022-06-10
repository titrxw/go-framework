package config

import "time"

type Session struct {
	Lifetime    time.Duration `mapstructure:"life_time" json:"life_time" yaml:"life_time"`
	IdleTimeout time.Duration `mapstructure:"idle_time" json:"idle_time" yaml:"idle_time"`
}
