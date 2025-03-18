package main

type ServerConfig struct {
	Host string
	Port string
}

var Config *ServerConfig

func InitServerConfig() *ServerConfig {
	Config = &ServerConfig{
		Host: "localhost",
		Port: "8080",
	}
	return Config
}
