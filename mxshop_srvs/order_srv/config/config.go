package config

type ServerConfig struct {
	Name       string         `mapstructure:"name"`
	Host       string         `mapstructure:"host"`
	Tags       []string       `mapstructure:"tags"`
	MysqlInfo  MysqlConfig    `mapstructure:"mysql"`
	ConsulInfo ConsulConfig   `mapstructure:"consul"`
	RedisInfo  RedisConfig    `mapstructure:"redis"`
	GoodsInfo  GoodsSrvConfig `mapstructure:"goods_info"`
	InvInfo    InvSrvConfig   `mapstructure:"inv_info"`
}

type GoodsSrvConfig struct {
	Name string `mapstructure:"name"`
}

type InvSrvConfig struct {
	Name string `mapstructure:"name"`
}

type RedisConfig struct {
	Host   string `mapstructuer:"host"`
	Port   int    `mapstructure:"port"`
	Expire int    `mapstruture:"expire"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"db"`
}

type NacosConfig struct {
	Host        string `mapstructuer:"host"`
	Port        uint64 `mapstructure:"port"`
	NamespaceId string `mapstructure:"namespace_id"`
	DataId      string `mapstructure:"data_id"`
	Group       string `mapstructure:"group"`
}
