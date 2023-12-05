package config

type DatabasePostgres struct {
	Host       string
	Port       string
	UserName   string
	Password   string
	SchemaName string
}

type DatabaseRedis struct {
	Host     string
	Port     string
	Password string
	Network  string
}

type DBConfig struct {
	DbPostgres DatabasePostgres `mapstructure:"db_postgres"`
	DbRedis    DatabaseRedis    `mapstructure:"db_redis"`
}

type MainConfig struct {
	Main Main `mapstructure:"main"`
}

type Main struct {
	HostGrpc string
	PortGrpc string
	Network  string
	PortHttp string
	Database string
}
