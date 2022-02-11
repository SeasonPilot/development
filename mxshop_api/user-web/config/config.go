package config

type UserSrvConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}
type JWTConfig struct {
	SigningKey string `mapstructure:"key"`
}

type ServerConfig struct {
	Name     string        `mapstructure:"name"`
	Port     int           `structure:"port"`
	UserInfo UserSrvConfig `mapstructure:"user-srv"`
	JWTInfo  JWTConfig     `mapstructure:"jwt"`
}
