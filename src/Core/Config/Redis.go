package config

type Redis struct {
	Host     string `mapstructure:"host" json:"host" yaml:"host"`
	Port     int    `mapstructure:"port" json:"port" yaml:"port"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Username string `mapstructure:"user_name" json:"user_name" yaml:"user_name"`
	Db       int    `mapstructure:"db" json:"db" yaml:"db"`
	PoolSize int    `mapstructure:"pool_size" json:"pool_size" yaml:"pool_size"`
}
