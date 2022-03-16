package config

type GoodsSrvConfig struct {
	Name string `mapstructure:"name"`
}

type JWTConfig struct {
	SigningKey string `mapstructure:"key"`
}

type ConsulConfig struct {
	Host string `mapstructuer:"host"`
	Port int    `mapstructure:"port"`
}

type AlipayConfig struct {
	AppID        string `mapstructure:"app_id"`
	PrivateKey   string `mapstructure:"private_key"`
	AliPublicKey string `mapstructure:"ali_public_key"`
	NotifyURL    string `mapstructure:"notify_url"`
	ReturnURL    string `mapstructure:"return_url"`
}

type JaegerConfig struct {
	Host string `mapstructuer:"host"`
	Port int    `mapstructure:"port"`
	Name string `mapstructure:"name"`
}

type ServerConfig struct {
	Name       string         `mapstructure:"name"`
	Address    string         `mapstructure:"address"`
	Port       int            `structure:"port"`
	Tags       []string       `mapstructure:"tags"`
	GoodsInfo  GoodsSrvConfig `mapstructure:"goods_info"`
	JWTInfo    JWTConfig      `mapstructure:"jwt"`
	ConsulInfo ConsulConfig   `mapstructure:"consul"`
	OrderInfo  GoodsSrvConfig `mapstructure:"order_info"`
	InvInfo    GoodsSrvConfig `mapstructure:"inv_info"`
	AlipayInfo AlipayConfig   `mapstructure:"alipay"`
	JaegerInfo JaegerConfig   `mapstructure:"jaeger"`
}

type NacosConfig struct {
	Host        string `mapstructuer:"host"`
	Port        uint64 `mapstructure:"port"`
	NamespaceId string `mapstructure:"namespace_id"`
	DataId      string `mapstructure:"data_id"`
	Group       string `mapstructure:"group"`
}
