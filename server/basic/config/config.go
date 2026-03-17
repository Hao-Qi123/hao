package config

type Mysql struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}
type Redis struct {
	Host     string
	Port     int
	Password string
	Database int
}
type Consul struct {
	Host        string
	Port        int
	ServiceName string
	ServicePort int
	TTL         int
}
type Nacos struct {
	Addr      string
	Port      int
	Namespace string
	DataId    string
	Group     string
}
type AliPay struct {
	AppID      string
	PrivateKey string
	NotifyUrl  string
	ReturnUrl  string
}
type AppConfig struct {
	Mysql
	Redis
	Consul
	Nacos
	AliPay
}
