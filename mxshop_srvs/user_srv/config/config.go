package config

type ServerConfig struct {
	MysqlInfo MysqlConfig `mapstructure:"mysql"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"db"`
}
