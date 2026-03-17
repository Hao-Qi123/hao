package bootstrap

func init() {
	InitAppConfig()
	InitNacos()
	InitMysql()
	InitRedis()
}
