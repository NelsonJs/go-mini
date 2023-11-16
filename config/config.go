package config

type Config struct {
	EC  EtcdConfig
	HC  HttpConfig
	RSC ReferenceServiceConfig
}

type EtcdConfig struct {
	Endpoints         []string
	Timeout           int
	ServiceNamePrefix string
	Version           string
}

type HttpConfig struct {
	Ip   string
	Port int
}

type ReferenceServiceConfig struct {
	Services []ServiceEnv
}

type ServiceEnv struct {
	ServiceName string
}
